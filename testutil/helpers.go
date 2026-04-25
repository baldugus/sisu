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

// CreateCall creates a new call and returns its ID.
func CreateCall(t *testing.T, db *database.Database) int32 {
	t.Helper()

	semester, err := database.FetchSemesterByYearAndNumber(db.DB(), 2025, 1)
	require.NoError(t, err, "failed to fetch semester")

	cmd := commands.CreateCallCommand{
		SemesterID: semester.ID,
	}
	err = cmd.Execute(db)
	require.NoError(t, err, "failed to create call")

	// Get the last call number
	lastCallNumber, err := db.GetLastCallNumber()
	require.NoError(t, err)

	call, err := db.FetchCallByNumber(lastCallNumber)
	require.NoError(t, err)

	return call.ID
}

// CloseCall closes a call by ID.
// Note: The call must have no pending registrations to be closed.
func CloseCall(t *testing.T, db *database.Database, callID int32) {
	t.Helper()

	cmd := commands.CloseCallCommand{ID: callID}
	err := cmd.Execute(db)
	require.NoError(t, err, "failed to close call")
}

// EnrollAllInCall enrolls all registrations in a call.
func EnrollAllInCall(t *testing.T, db *database.Database, callID int32) {
	t.Helper()

	registrations, err := db.FetchRegistrationsByCallID(callID)
	require.NoError(t, err, "failed to fetch registrations")

	for _, reg := range registrations {
		if reg.Status == types.RegistrationStatusApproved {
			EnrollRegistration(t, db, reg.ID)
		}
	}
}

// CloseCallWithEnrollment enrolls all students and then closes the call.
func CloseCallWithEnrollment(t *testing.T, db *database.Database, callID int32) {
	t.Helper()

	EnrollAllInCall(t, db, callID)
	CloseCall(t, db, callID)
}

// EnrollRegistration marks a registration as enrolled.
func EnrollRegistration(t *testing.T, db *database.Database, regID int32) {
	t.Helper()

	cmd := commands.UpdateRegistrationStatusCommand{
		RegistrationID: regID,
		NewStatus:      types.RegistrationStatusEnrolled,
	}

	err := cmd.Execute(db)
	require.NoError(t, err, "failed to enroll registration")
}

// MarkRegistrationAbsent marks a registration as absent.
func MarkRegistrationAbsent(t *testing.T, db *database.Database, regID int32) {
	t.Helper()

	cmd := commands.UpdateRegistrationStatusCommand{
		RegistrationID: regID,
		NewStatus:      types.RegistrationStatusAbsent,
	}

	err := cmd.Execute(db)
	require.NoError(t, err, "failed to mark registration absent")
}

// ClearRegistrationStatus resets a registration status to approved.
func ClearRegistrationStatus(t *testing.T, db *database.Database, regID int32) {
	t.Helper()

	cmd := commands.UpdateRegistrationStatusCommand{
		RegistrationID: regID,
		NewStatus:      types.RegistrationStatusApproved,
	}

	err := cmd.Execute(db)
	require.NoError(t, err, "failed to clear registration status")
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
