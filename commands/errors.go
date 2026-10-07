package commands

type ErrApprovedSelectionAlreadyExists struct{}

func (e ErrApprovedSelectionAlreadyExists) Error() string {
	return "cannot load approved list when another selection exists"
}

type ErrWaitlistSelectionRequiresApproved struct{}

func (e ErrWaitlistSelectionRequiresApproved) Error() string {
	return "cannot load waitlist without approved list"
}

type ErrCallHasPendingRegistrations struct{}

func (e ErrCallHasPendingRegistrations) Error() string {
	return "cannot close rollcall with pending registrations"
}

type ErrSelectionNotFound struct{}

func (e ErrSelectionNotFound) Error() string {
	return "selection not found"
}

type ErrCannotDeleteApprovedWithWaitlist struct{}

func (e ErrCannotDeleteApprovedWithWaitlist) Error() string {
	return "cannot delete approved selection while waitlist selection exists"
}

type ErrSelectionHasModifiedRegistrations struct{}

func (e ErrSelectionHasModifiedRegistrations) Error() string {
	return "cannot delete selection with enrolled or absent registrations"
}

type ErrRegistrationNotFound struct{}

func (e ErrRegistrationNotFound) Error() string {
	return "registration not found"
}

type ErrCallNotOpen struct{}

func (e ErrCallNotOpen) Error() string {
	return "call is not open"
}

type ErrInvalidStatusTransition struct{}

func (e ErrInvalidStatusTransition) Error() string {
	return "invalid status transition"
}

type ErrOpenCallExists struct{}

func (e ErrOpenCallExists) Error() string {
	return "cannot create new call while another is open"
}

type ErrAllCoursesFull struct{}

func (e ErrAllCoursesFull) Error() string {
	return "all seats in all courses are occupied"
}

type ErrCannotReopenCallWithLaterCalls struct{}

func (e ErrCannotReopenCallWithLaterCalls) Error() string {
	return "cannot reopen call while there are calls after it"
}

type ErrCannotDeleteWithMultipleCalls struct{}

func (e ErrCannotDeleteWithMultipleCalls) Error() string {
	return "cannot delete selection when there are multiple calls"
}

type ErrCannotDeleteWithClosedFirstCall struct{}

func (e ErrCannotDeleteWithClosedFirstCall) Error() string {
	return "cannot delete selection when first call is closed"
}

type ErrCannotDeleteCallWithLaterCalls struct{}

func (e ErrCannotDeleteCallWithLaterCalls) Error() string {
	return "cannot delete call while there are calls after it"
}

type ErrCannotDeleteClosedCall struct{}

func (e ErrCannotDeleteClosedCall) Error() string {
	return "cannot delete a closed call"
}

type ErrCallNotFound struct{}

func (e ErrCallNotFound) Error() string {
	return "call not found"
}

type ErrCannotDeleteFirstCall struct{}

func (e ErrCannotDeleteFirstCall) Error() string {
	return "the first call can only be removed by deleting the approved selection"
}

type ErrNoCandidatesToCall struct{}

func (e ErrNoCandidatesToCall) Error() string {
	return "no one to call: no promotion requests nor waitlisted registrations for the open seats"
}

type ErrNoCalls struct{}

func (e ErrNoCalls) Error() string {
	return "no calls exist; import the approved list first"
}

type ErrInvalidSemester struct{}

func (e ErrInvalidSemester) Error() string {
	return "semester must be 1 or 2"
}

type ErrPromotionNotAllowed struct{}

func (e ErrPromotionNotAllowed) Error() string {
	return "only semester-2 students who were not absent can ask for promotion"
}
