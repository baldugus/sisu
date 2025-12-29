package types

import "encoding/json"

// ENUM(approved, waitlist)
type SelectionKind int

func (s SelectionKind) ToRegistrationStatus() RegistrationStatus {
	switch s.String() {
	case "approved":
		return RegistrationStatusApproved
	case "waitlist":
		return RegistrationStatusWaitlisted
	default:
		return -1
	}
}

func (s SelectionKind) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

type Selection struct {
	ID          int32
	Name        string
	Kind        SelectionKind
	Year        int32
	Semester    int32
	Institution string
	Degree      string
}
