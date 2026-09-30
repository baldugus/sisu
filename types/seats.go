package types

import "encoding/json"

// Seats is a course's seat count. It is always even: the seats are split
// evenly between the year's two semester intakes, so the type stores the
// per-semester count and an odd total cannot be represented. The zero value
// is a valid count of zero seats.
type Seats struct {
	perSemester int32
}

// NewSeats returns the seat count for a course with the given total number of
// seats, or ErrOddSeatsCount if the total is odd.
func NewSeats(total int32) (Seats, error) {
	if total%2 != 0 {
		return Seats{}, ErrOddSeatsCount{Count: total}
	}

	return Seats{perSemester: total / 2}, nil
}

// Total returns the course's total number of seats.
func (s Seats) Total() int32 {
	return s.perSemester * 2
}

// PerSemester returns the number of seats of each semester intake.
func (s Seats) PerSemester() int32 {
	return s.perSemester
}

// MarshalJSON serializes the seats as the total count.
func (s Seats) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Total())
}
