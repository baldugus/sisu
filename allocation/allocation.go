// Package allocation holds the rule that decides who is called when a new call
// is opened. It is a pure function over plain data so it can be tested without
// a database.
//
// For each course (time slot × quota) the rule is:
//
//  1. Semester 1 vacancies are offered first to semester-2 students who asked
//     to move up (promotion), best ranked first.
//  2. The remaining semester 1 vacancies, then the semester 2 vacancies, are
//     filled from the waitlist, best ranked first.
//
// A closed semester receives no one. Seats freed by a promotion only become
// available in the next call, once the promotion is accepted, so a student who
// declines never leaves semester 2 overbooked.
package allocation

import (
	"cmp"
	"slices"

	"github.com/baldugus/sisu/types"
)

// Candidate is a registration that may be called.
type Candidate struct {
	RegistrationID int32
	Ranking        int32
	CompositeScore int32
}

// Course is the allocation input for one course.
type Course struct {
	ID    int32
	Seats int32
	// Occupied is indexed by semester - 1.
	Occupied [2]int32
	// Promotion lists semester-2 students eligible to move to semester 1.
	Promotion []Candidate
	// Waitlist lists waitlist registrations that were never called.
	Waitlist []Candidate
}

// Input is everything the rule needs.
type Input struct {
	Courses []Course
	// Open is indexed by semester - 1.
	Open [2]bool
}

// CourseResult is who gets called for one course.
type CourseResult struct {
	CourseID int32
	// Vacancies is indexed by semester - 1; closed semesters report 0.
	Vacancies [2]int32
	// Promoted are offered a seat in semester 1.
	Promoted []int32
	// Waitlist is indexed by semester - 1.
	Waitlist [2][]int32
}

// Total is the number of registrations called for the course.
func (r *CourseResult) Total() int {
	return len(r.Promoted) + len(r.Waitlist[0]) + len(r.Waitlist[1])
}

// Plan applies the rule to every course, in input order.
func Plan(in Input) []CourseResult {
	results := make([]CourseResult, len(in.Courses))

	for i, course := range in.Courses {
		results[i] = planCourse(course, in.Open)
	}

	return results
}

func planCourse(course Course, open [2]bool) CourseResult {
	result := CourseResult{CourseID: course.ID}

	for s := range 2 {
		if !open[s] {
			continue
		}

		seats := types.SemesterSeats(course.Seats, int32(s+1)) //nolint: gosec
		result.Vacancies[s] = max(0, seats-course.Occupied[s])
	}

	promotion := SortByRanking(course.Promotion)
	waitlist := SortByRanking(course.Waitlist)

	free1 := int(result.Vacancies[0])
	for _, c := range promotion[:min(free1, len(promotion))] {
		result.Promoted = append(result.Promoted, c.RegistrationID)
	}

	free1 -= len(result.Promoted)

	next := 0
	for s, free := range []int{free1, int(result.Vacancies[1])} {
		n := min(free, len(waitlist)-next)
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
