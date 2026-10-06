package types

// ENUM(calling, done)
type CallStatus string

type Call struct {
	ID     int32      `csv:"-"`
	Status CallStatus `csv:"-"`
	Number int32      `csv:"CHAMADA"`
}

// CallSummary is a call plus per-semester counts of the entries it holds.
type CallSummary struct {
	ID        int32
	Status    CallStatus
	Number    int32
	Pending   int32
	Semesters []*CallSemesterSummary
}

type CallSemesterSummary struct {
	Semester int32
	Initial  int32
	Waitlist int32
}
