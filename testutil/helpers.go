package testutil

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/baldugus/sisu/commands"
	"github.com/baldugus/sisu/database"
	"github.com/baldugus/sisu/types"
)

// LoadApprovedSelection loads an approved selection from a CSV file.
func LoadApprovedSelection(t *testing.T, db *database.Database, filePath string) {
	t.Helper()

	cmd := commands.LoadSelectionCommand{
		Year:     2025,
		FilePath: filePath,
		Kind:     types.SelectionKindApproved,
	}

	err := cmd.Execute(db)
	require.NoError(t, err, "failed to load approved selection")
}

// LoadWaitlistSelection loads a waitlist selection from a CSV file.
func LoadWaitlistSelection(t *testing.T, db *database.Database, filePath string) {
	t.Helper()

	cmd := commands.LoadSelectionCommand{
		Year:     2025,
		FilePath: filePath,
		Kind:     types.SelectionKindWaitlist,
	}

	err := cmd.Execute(db)
	require.NoError(t, err, "failed to load waitlist selection")
}

// CreateCall creates the next call and returns its ID.
func CreateCall(t *testing.T, db *database.Database) int32 {
	t.Helper()

	cmd := commands.CreateCallCommand{}
	err := cmd.Execute(db)
	require.NoError(t, err, "failed to create call")

	lastCallNumber, err := db.GetLastCallNumber()
	require.NoError(t, err)

	call, err := db.FetchCallByNumber(lastCallNumber)
	require.NoError(t, err)

	return call.ID
}

// CloseCall closes a call by ID.
// Note: The call must have no pending entries to be closed.
func CloseCall(t *testing.T, db *database.Database, callID int32) {
	t.Helper()

	cmd := commands.CloseCallCommand{ID: callID}
	err := cmd.Execute(db)
	require.NoError(t, err, "failed to close call")
}

// CallEntries returns the entries of a call, best ranked first.
func CallEntries(t *testing.T, db *database.Database, callID int32) []*types.CallEntryDetail {
	t.Helper()

	entries, err := database.FetchCallEntryDetails(db.DB(), callID, nil, nil)
	require.NoError(t, err, "failed to fetch call entries")

	return entries
}

// RegistrationsInCall returns the registrations of a call, best ranked first.
func RegistrationsInCall(t *testing.T, db *database.Database, callID int32) []*types.Registration {
	t.Helper()

	entries := CallEntries(t, db, callID)

	regs := make([]*types.Registration, len(entries))
	for i, e := range entries {
		regs[i] = e.Registration
	}

	return regs
}

// EnrollAllInCall marks every pending entry of a call as enrolled.
func EnrollAllInCall(t *testing.T, db *database.Database, callID int32) {
	t.Helper()

	for _, e := range CallEntries(t, db, callID) {
		if e.Entry.Outcome == types.CallEntryOutcomePending {
			SetOutcome(t, db, callID, e.Entry.RegistrationID, types.CallEntryOutcomeEnrolled)
		}
	}
}

// CloseCallWithEnrollment enrolls all students and then closes the call.
func CloseCallWithEnrollment(t *testing.T, db *database.Database, callID int32) {
	t.Helper()

	EnrollAllInCall(t, db, callID)
	CloseCall(t, db, callID)
}

// SetOutcome records a registration's outcome in a call.
func SetOutcome(
	t *testing.T,
	db *database.Database,
	callID int32,
	regID int32,
	outcome types.CallEntryOutcome,
) {
	t.Helper()

	cmd := commands.SetCallEntryOutcomeCommand{
		CallID:         callID,
		RegistrationID: regID,
		Outcome:        outcome,
	}

	err := cmd.Execute(db)
	require.NoError(t, err, "failed to set outcome %s", outcome)
}

// latestCallID returns the call of the registration's most recent entry.
func latestCallID(t *testing.T, db *database.Database, regID int32) int32 {
	t.Helper()

	history, err := database.FetchRegistrationHistory(db.DB(), regID)
	require.NoError(t, err)
	require.NotEmpty(t, history, "registration %d has no call entries", regID)

	return history[len(history)-1].CallID
}

// EnrollRegistration marks a registration as enrolled in its latest call.
func EnrollRegistration(t *testing.T, db *database.Database, regID int32) {
	t.Helper()

	SetOutcome(t, db, latestCallID(t, db, regID), regID, types.CallEntryOutcomeEnrolled)
}

// MarkRegistrationAbsent marks a registration as absent in its latest call.
func MarkRegistrationAbsent(t *testing.T, db *database.Database, regID int32) {
	t.Helper()

	SetOutcome(t, db, latestCallID(t, db, regID), regID, types.CallEntryOutcomeAbsent)
}

// ClearRegistrationStatus resets a registration's latest entry to pending.
func ClearRegistrationStatus(t *testing.T, db *database.Database, regID int32) {
	t.Helper()

	SetOutcome(t, db, latestCallID(t, db, regID), regID, types.CallEntryOutcomePending)
}

// DeleteApprovedSelection deletes the approved selection.
func DeleteApprovedSelection(t *testing.T, db *database.Database) {
	t.Helper()

	cmd := commands.DeleteSelectionCommand{
		Kind: types.SelectionKindApproved,
	}

	err := cmd.Execute(db)
	require.NoError(t, err, "failed to delete approved selection")
}

// DeleteWaitlistSelection deletes the waitlist selection.
func DeleteWaitlistSelection(t *testing.T, db *database.Database) {
	t.Helper()

	cmd := commands.DeleteSelectionCommand{
		Kind: types.SelectionKindWaitlist,
	}

	err := cmd.Execute(db)
	require.NoError(t, err, "failed to delete waitlist selection")
}
