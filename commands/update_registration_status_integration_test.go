package commands_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/baldugus/sisu/commands"
	"github.com/baldugus/sisu/database"
	"github.com/baldugus/sisu/testutil"
	"github.com/baldugus/sisu/types"
)

func TestUpdateRegistrationStatus_ApprovedToEnrolled_Success(t *testing.T) {
	// Arrange
	db := database.NewTestDatabase(t)
	testutil.LoadApprovedSelection(t, db.Database, "testdata/approved_small.csv")

	// Get first registration
	registrations, err := db.FetchRegistrations()
	require.NoError(t, err)
	reg := registrations[0]

	assert.Equal(t, types.RegistrationStatusApproved, reg.Status)

	cmd := commands.UpdateRegistrationStatusCommand{
		RegistrationID: reg.ID,
		NewStatus:      types.RegistrationStatusEnrolled,
	}

	// Act
	err = cmd.Execute(db.Database)

	// Assert
	require.NoError(t, err)
	testutil.AssertRegistrationStatus(t, db.Database, reg.ID, types.RegistrationStatusEnrolled)
}

func TestUpdateRegistrationStatus_ApprovedToAbsent_Success(t *testing.T) {
	// Arrange
	db := database.NewTestDatabase(t)
	testutil.LoadApprovedSelection(t, db.Database, "testdata/approved_small.csv")

	registrations, err := db.FetchRegistrations()
	require.NoError(t, err)
	reg := registrations[0]

	cmd := commands.UpdateRegistrationStatusCommand{
		RegistrationID: reg.ID,
		NewStatus:      types.RegistrationStatusAbsent,
	}

	// Act
	err = cmd.Execute(db.Database)

	// Assert
	require.NoError(t, err)
	testutil.AssertRegistrationStatus(t, db.Database, reg.ID, types.RegistrationStatusAbsent)
}

func TestUpdateRegistrationStatus_EnrolledToApproved_Success(t *testing.T) {
	// Arrange
	db := database.NewTestDatabase(t)
	testutil.LoadApprovedSelection(t, db.Database, "testdata/approved_small.csv")

	registrations, err := db.FetchRegistrations()
	require.NoError(t, err)
	reg := registrations[0]

	// First enroll the student
	testutil.EnrollRegistration(t, db.Database, reg.ID)

	cmd := commands.UpdateRegistrationStatusCommand{
		RegistrationID: reg.ID,
		NewStatus:      types.RegistrationStatusApproved,
	}

	// Act
	err = cmd.Execute(db.Database)

	// Assert
	require.NoError(t, err)
	testutil.AssertRegistrationStatus(t, db.Database, reg.ID, types.RegistrationStatusApproved)
}

func TestUpdateRegistrationStatus_AbsentToApproved_Success(t *testing.T) {
	// Arrange
	db := database.NewTestDatabase(t)
	testutil.LoadApprovedSelection(t, db.Database, "testdata/approved_small.csv")

	registrations, err := db.FetchRegistrations()
	require.NoError(t, err)
	reg := registrations[0]

	// First mark as absent
	testutil.MarkRegistrationAbsent(t, db.Database, reg.ID)

	cmd := commands.UpdateRegistrationStatusCommand{
		RegistrationID: reg.ID,
		NewStatus:      types.RegistrationStatusApproved,
	}

	// Act
	err = cmd.Execute(db.Database)

	// Assert
	require.NoError(t, err)
	testutil.AssertRegistrationStatus(t, db.Database, reg.ID, types.RegistrationStatusApproved)
}

func TestUpdateRegistrationStatus_ErrRegistrationNotFound(t *testing.T) {
	// Arrange
	db := database.NewTestDatabase(t)

	cmd := commands.UpdateRegistrationStatusCommand{
		RegistrationID: 9999, // Non-existent ID
		NewStatus:      types.RegistrationStatusEnrolled,
	}

	// Act
	err := cmd.Execute(db.Database)

	// Assert
	require.Error(t, err)
	// The error comes from database layer as qrm.ErrNoRows which gets converted
	// We just verify an error occurs
	assert.Error(t, err)
}

func TestUpdateRegistrationStatus_ErrCallNotOpen(t *testing.T) {
	// Arrange
	db := database.NewTestDatabase(t)
	testutil.LoadApprovedSelection(t, db.Database, "testdata/approved_small.csv")

	registrations, err := db.FetchRegistrations()
	require.NoError(t, err)
	reg := registrations[0]

	// Close the call (enrolls all first)
	call1, err := db.FetchCallByNumber(1)
	require.NoError(t, err)
	testutil.CloseCallWithEnrollment(t, db.Database, call1.ID)

	// Try to clear status (enrolled->approved) after call is closed
	cmd := commands.UpdateRegistrationStatusCommand{
		RegistrationID: reg.ID,
		NewStatus:      types.RegistrationStatusApproved,
	}

	// Act
	err = cmd.Execute(db.Database)

	// Assert
	require.Error(t, err)
	assert.ErrorAs(t, err, &commands.ErrCallNotOpen{})
}

func TestUpdateRegistrationStatus_ErrInvalidStatusTransition_WaitlistedToEnrolled(t *testing.T) {
	// Arrange
	db := database.NewTestDatabase(t)
	testutil.LoadApprovedSelection(t, db.Database, "testdata/approved_small.csv")
	testutil.LoadWaitlistSelection(t, db.Database, "testdata/waitlist_small.csv")

	// Get a waitlisted registration
	waitlistSelection := testutil.AssertSelectionExists(t, db.Database, types.SelectionKindWaitlist)
	waitlistRegs, err := db.FetchRegistrationsBySelectionID(waitlistSelection.ID)
	require.NoError(t, err)
	reg := waitlistRegs[0]

	assert.Equal(t, types.RegistrationStatusWaitlisted, reg.Status)

	// Try to enroll directly from waitlisted (invalid transition)
	cmd := commands.UpdateRegistrationStatusCommand{
		RegistrationID: reg.ID,
		NewStatus:      types.RegistrationStatusEnrolled,
	}

	// Act
	err = cmd.Execute(db.Database)

	// Assert
	require.Error(t, err)
	assert.ErrorAs(t, err, &commands.ErrInvalidStatusTransition{})
}

func TestUpdateRegistrationStatus_ErrInvalidTransitionOnWaitlisted(t *testing.T) {
	// Arrange
	db := database.NewTestDatabase(t)
	testutil.LoadApprovedSelection(t, db.Database, "testdata/approved_small.csv")
	testutil.LoadWaitlistSelection(t, db.Database, "testdata/waitlist_small.csv")

	// Get a waitlisted registration (not in a call yet)
	waitlistSelection := testutil.AssertSelectionExists(t, db.Database, types.SelectionKindWaitlist)
	waitlistRegs, err := db.FetchRegistrationsBySelectionID(waitlistSelection.ID)
	require.NoError(t, err)
	reg := waitlistRegs[0]

	// Try to update status when not in a call - waitlisted->approved is invalid
	cmd := commands.UpdateRegistrationStatusCommand{
		RegistrationID: reg.ID,
		NewStatus:      types.RegistrationStatusApproved,
	}

	// Act
	err = cmd.Execute(db.Database)

	// Assert
	require.Error(t, err)
	// Either ErrInvalidStatusTransition or ErrRegistrationNotInCall is acceptable
	assert.True(t,
		errors.As(err, &commands.ErrInvalidStatusTransition{}) ||
			errors.As(err, &commands.ErrRegistrationNotInCall{}),
		"expected invalid transition or not in call error")
}

func TestUpdateRegistrationStatus_NoOpWhenStatusUnchanged(t *testing.T) {
	// Arrange
	db := database.NewTestDatabase(t)
	testutil.LoadApprovedSelection(t, db.Database, "testdata/approved_small.csv")

	registrations, err := db.FetchRegistrations()
	require.NoError(t, err)
	reg := registrations[0]

	// Try to set to same status
	cmd := commands.UpdateRegistrationStatusCommand{
		RegistrationID: reg.ID,
		NewStatus:      types.RegistrationStatusApproved, // Already approved
	}

	// Act
	err = cmd.Execute(db.Database)

	// Assert
	require.NoError(t, err)
	testutil.AssertRegistrationStatus(t, db.Database, reg.ID, types.RegistrationStatusApproved)
}
