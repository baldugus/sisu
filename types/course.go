package types

// ENUM(morning, evening)
type CoursePeriod string

type Course struct {
	ID           int32        `csv:"-"`
	Seats        int32        `csv:"-"`
	MinimumScore *Score       `csv:"-" ts_type:"string"`
	Period       CoursePeriod `csv:"TURNO"`
	Quota        string       `csv:"MODALIDADE"`
}

// SemesterForRanking returns the semester intake (1 or 2) of the candidate with
// the given ranking in this course: rankings within the first half of the seats
// go to Semester 1, the rest to Semester 2. The seat count must be even,
// otherwise ErrOddSeatsCount is returned.
func (c Course) SemesterForRanking(ranking int32) (int32, error) {
	if c.Seats%2 != 0 {
		return 0, ErrOddSeatsCount{Count: c.Seats}
	}

	if ranking <= c.Seats/2 {
		return 1, nil
	}

	return 2, nil
}
