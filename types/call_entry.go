package types

// ENUM(initial, waitlist, promotion)
type CallEntryKind string

// ENUM(pending, enrolled, absent, declined)
type CallEntryOutcome string

// CallEntry records one registration taking part in one call. Everything that
// happens during the admission cycle lives in call entries, owned by the call
// that produced them, so deleting a call reverts all of its effects.
//
//   - initial:   approved list, placed by the import in call 1.
//   - waitlist:  called from the waitlist.
//   - promotion: a semester-2 student offered a seat in semester 1. Outcome
//     enrolled means accepted; declined leaves the student in semester 2.
type CallEntry struct {
	CallID         int32
	CallNumber     int32
	RegistrationID int32
	Kind           CallEntryKind
	Semester       int32
	Outcome        CallEntryOutcome
	// WantsPromotion is set during a call's enrollment window on semester-2
	// students who asked to move to semester 1 in a later call.
	WantsPromotion bool
}

// CallEntryDetail is a call entry with the registration and course it refers to.
type CallEntryDetail struct {
	Entry        *CallEntry
	Registration *Registration
	Course       *Course
}
