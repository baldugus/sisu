package types

// ENUM(morning, evening)
type CoursePeriod string

type Course struct {
	ID           int32        `csv:"-"`
	Seats        Seats        `csv:"-" ts_type:"number"`
	MinimumScore *Score       `csv:"-" ts_type:"string"`
	Period       CoursePeriod `csv:"TURNO"`
	Quota        string       `csv:"MODALIDADE"`
}

// SemesterForRanking returns the semester intake (1 or 2) of the candidate with
// the given ranking in this course: rankings within the first half of the seats
// go to Semester 1, the rest to Semester 2.
func (c Course) SemesterForRanking(ranking int32) int32 {
	if ranking <= c.Seats.PerSemester() {
		return 1
	}

	return 2
}
