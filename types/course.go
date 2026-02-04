package types

import "encoding/json"

// ENUM(morning, evening)
type CoursePeriod int

func (c CoursePeriod) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

type Course struct {
	ID           int32        `csv:"-"`
	Seats        int32        `csv:"-"`
	MinimumScore *Score       `csv:"-"`
	Period       CoursePeriod `csv:"TURNO"`
	Quota        string       `csv:"MODALIDADE"`
}
