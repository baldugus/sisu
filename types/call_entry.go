package types

// ENUM(initial, waitlist)
type CallEntryKind string

// ENUM(pending, enrolled, absent)
type CallEntryOutcome string

// CallEntry records one registration taking part in one call. Everything that
// happens during the admission cycle lives in call entries, owned by the call
// that produced them, so deleting a call reverts all of its effects.
//
//   - initial:  approved list, placed by the import in call 1.
//   - waitlist: called from the waitlist.
type CallEntry struct {
	CallID         int32
	CallNumber     int32
	RegistrationID int32
	Kind           CallEntryKind
	Semester       int32
	Outcome        CallEntryOutcome
}

// CallEntryDetail is a call entry with the registration and course it refers to.
type CallEntryDetail struct {
	Entry        *CallEntry
	Registration *Registration
	Course       *Course
}
