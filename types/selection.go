package types

// ENUM(approved, waitlist)
type SelectionKind string

func (s SelectionKind) ToRegistrationStatus() RegistrationStatus {
	switch s.String() {
	case "approved":
		return RegistrationStatusApproved
	case "waitlist":
		return RegistrationStatusWaitlisted
	default:
		return ""
	}
}

type Selection struct {
	ID          int32
	Name        string
	Kind        SelectionKind
	Year        int32
	Institution string
	Degree      string
}
