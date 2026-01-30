package commands_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/baldugus/sisu/commands"
	"github.com/baldugus/sisu/database"
	"github.com/baldugus/sisu/testutil"
	"github.com/baldugus/sisu/types"
)

func TestDeleteCallCommand_CannotDeleteClosedCall(t *testing.T) {
	db := database.NewTestDatabase(t)

	testutil.LoadApprovedSelection(t, db.Database, "testdata/approved_small.csv")

	call1, err := db.FetchCallByNumber(1)
	require.NoError(t, err)

	testutil.CloseCallWithEnrollment(t, db.Database, call1.ID)

	cmd := commands.DeleteCallCommand{
		ID: call1.ID,
	}
	err = cmd.Execute(db.Database)

	testutil.AssertErrorType(t, err, commands.ErrCannotDeleteClosedCall{})

	call, err := db.FetchCallByID(call1.ID)
	require.NoError(t, err)
	assert.Equal(t, types.CallStatusDone, call.Status, "closed call should still exist")
}

func TestDeleteCallCommand_CanDeleteOpenCall(t *testing.T) {
	db := database.NewTestDatabase(t)

	testutil.LoadApprovedSelection(t, db.Database, "testdata/approved_small.csv")
	testutil.LoadWaitlistSelection(t, db.Database, "testdata/waitlist_small.csv")

	call1, err := db.FetchCallByNumber(1)
	require.NoError(t, err)

	registrations, err := db.FetchRegistrationsByCallID(call1.ID)
	require.NoError(t, err)

	enrollCount := 0
	for _, reg := range registrations {
		if reg.Status == types.RegistrationStatusApproved {
			if enrollCount < 2 {
				testutil.EnrollRegistration(t, db.Database, reg.ID)
				enrollCount++
			} else {
				testutil.MarkRegistrationAbsent(t, db.Database, reg.ID)
			}
		}
	}

	testutil.CloseCall(t, db.Database, call1.ID)
	call2ID := testutil.CreateCall(t, db.Database)

	cmd := commands.DeleteCallCommand{
		ID: call2ID,
	}
	err = cmd.Execute(db.Database)
	require.NoError(t, err)

	_, err = db.FetchCallByID(call2ID)
	require.Error(t, err, "call 2 should be deleted")

	_, err = db.FetchCallByID(call1.ID)
	require.NoError(t, err, "call 1 should still exist")
}

func TestDeleteCallCommand_CannotDeleteNonLastCall(t *testing.T) {
	db := database.NewTestDatabase(t)

	testutil.LoadApprovedSelection(t, db.Database, "testdata/approved_small.csv")
	testutil.LoadWaitlistSelection(t, db.Database, "testdata/waitlist_small.csv")

	call1, err := db.FetchCallByNumber(1)
	require.NoError(t, err)

	registrations, err := db.FetchRegistrationsByCallID(call1.ID)
	require.NoError(t, err)

	enrollCount := 0
	for _, reg := range registrations {
		if reg.Status == types.RegistrationStatusApproved {
			if enrollCount < 2 {
				testutil.EnrollRegistration(t, db.Database, reg.ID)
				enrollCount++
			} else {
				testutil.MarkRegistrationAbsent(t, db.Database, reg.ID)
			}
		}
	}

	testutil.CloseCall(t, db.Database, call1.ID)
	testutil.CreateCall(t, db.Database)

	cmd := commands.DeleteCallCommand{
		ID: call1.ID,
	}
	err = cmd.Execute(db.Database)

	testutil.AssertErrorType(t, err, commands.ErrCannotDeleteClosedCall{})
}
