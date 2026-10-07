package types

// Semester is one of the two fixed intakes of the admission cycle (Number 1 or 2).
// Both rows always exist.
type Semester struct {
	Number int32
	// Seats is the semester's share of all course seats (half of each course).
	Seats int32
	// Occupied counts registrations currently placed in the semester (pending
	// or enrolled).
	Occupied int32
}
