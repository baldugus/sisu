package types

type RegistrationDetail struct {
	Registration *Registration
	Course       *Course
	// History lists every call entry of the registration, oldest call first.
	History []*CallEntry
}
