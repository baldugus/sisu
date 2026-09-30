package types

// ENUM(open, closed)
type SemesterStatus string

// Semester is one of the two fixed intakes of the admission cycle (Number 1 or 2).
// Both rows always exist; only their closure changes.
type Semester struct {
	Number int32
	Status SemesterStatus
	// ClosedAfterCall is the number of the last call when the semester was
	// closed. A closed semester receives no one in later calls; it can be
	// reopened only while no call was created after the closure.
	ClosedAfterCall *int32
	// Seats is the semester's share of all course seats (half of each course).
	Seats int32
	// Occupied counts registrations currently placed in the semester (pending
	// or enrolled).
	Occupied int32
}

// SemesterSeats is a course's share of seats in one semester: half each, with
// semester 1 taking the extra seat when the total is odd.
func SemesterSeats(courseSeats int32, semester int32) int32 {
	if semester == 1 {
		return (courseSeats + 1) / 2 //nolint: mnd
	}

	return courseSeats / 2 //nolint: mnd
}
