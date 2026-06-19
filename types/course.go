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
