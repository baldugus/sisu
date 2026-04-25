package types

// ENUM(open, closed)
type SemesterStatus int

type Semester struct {
	ID     int32
	Year   int32
	Number int32
	Status SemesterStatus
}
