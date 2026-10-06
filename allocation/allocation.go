// Package allocation holds the rule that decides who is called when a new call
// is opened. It is a pure function over plain data so it can be tested without
// a database.
//
// For each course (time slot × quota), the semester 1 vacancies, then the
// semester 2 vacancies, are filled from the waitlist, best ranked first.
package allocation

import (
	"cmp"
	"slices"
)

// Candidate is a registration that may be called.
type Candidate struct {
	RegistrationID int32
	Ranking        int32
	CompositeScore int32
}

// Course is the allocation input for one course.
type Course struct {
	ID int32
	// SeatsPerSemester is each semester's share of the course's seats.
	SeatsPerSemester int32
	// Occupied is indexed by semester - 1.
	Occupied [2]int32
	// Waitlist lists waitlist registrations that were never called.
	Waitlist []Candidate
}

// Input is everything the rule needs.
type Input struct {
	Courses []Course
}

// CourseResult is who gets called for one course.
type CourseResult struct {
	CourseID int32
	// Vacancies is indexed by semester - 1.
	Vacancies [2]int32
	// Waitlist is indexed by semester - 1.
	Waitlist [2][]int32
}

// Total is the number of registrations called for the course.
func (r *CourseResult) Total() int {
	return len(r.Waitlist[0]) + len(r.Waitlist[1])
}

// Plan applies the rule to every course, in input order.
func Plan(in Input) []CourseResult {
	results := make([]CourseResult, len(in.Courses))

	for i, course := range in.Courses {
		results[i] = planCourse(course)
	}

	return results
}

func planCourse(course Course) CourseResult {
	result := CourseResult{CourseID: course.ID}

	for s := range 2 {
		result.Vacancies[s] = max(0, course.SeatsPerSemester-course.Occupied[s])
	}

	waitlist := SortByRanking(course.Waitlist)

	next := 0
	for s, free := range result.Vacancies {
		n := min(int(free), len(waitlist)-next)
		for _, c := range waitlist[next : next+n] {
			result.Waitlist[s] = append(result.Waitlist[s], c.RegistrationID)
		}

		next += n
	}

	return result
}

// SortByRanking returns the candidates best first: lowest ranking, with a
// missing ranking (0, common in waitlist files) last; ties broken by higher
// composite score, then registration ID for determinism.
func SortByRanking(candidates []Candidate) []Candidate {
	sorted := slices.Clone(candidates)

	slices.SortStableFunc(sorted, func(a, b Candidate) int {
		if (a.Ranking == 0) != (b.Ranking == 0) {
			if a.Ranking == 0 {
				return 1
			}

			return -1
		}

		if c := cmp.Compare(a.Ranking, b.Ranking); c != 0 {
			return c
		}

		if c := cmp.Compare(b.CompositeScore, a.CompositeScore); c != 0 {
			return c
		}

		return cmp.Compare(a.RegistrationID, b.RegistrationID)
	})

	return sorted
}
