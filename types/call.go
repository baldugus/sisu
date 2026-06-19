package types

// ENUM(calling, done)
type CallStatus string

type Call struct {
	ID         int32      `csv:"-"`
	Status     CallStatus `csv:"-"`
	Number     int32      `csv:"CHAMADA"`
	SemesterID int32      `csv:"-"`
}
