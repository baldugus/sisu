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

func TestDeleteCallCommand_CannotDeleteFirstCall(t *testing.T) {
	db := database.NewTestDatabase(t)

	testutil.LoadApprovedSelection(t, db.Database, "testdata/approved_small.csv")

	call1, err := db.FetchCallByNumber(1)
	require.NoError(t, err)

	cmd := commands.DeleteCallCommand{ID: call1.ID}
	err = cmd.Execute(db.Database)

	testutil.AssertErrorType(t, err, commands.ErrCannotDeleteFirstCall{})

	_, err = db.FetchCallByID(call1.ID)
	require.NoError(t, err, "call 1 should still exist")
}

func TestDeleteCallCommand_CannotDeleteClosedCall(t *testing.T) {
	db := database.NewTestDatabase(t)

	testutil.LoadApprovedSelection(t, db.Database, "testdata/approved_small.csv")
	testutil.LoadWaitlistSelection(t, db.Database, "testdata/waitlist_small.csv")

	call1, err := db.FetchCallByNumber(1)
	require.NoError(t, err)
	testutil.CloseCallWithEnrollment(t, db.Database, call1.ID)

	call2 := testutil.CreateCall(t, db.Database)
	testutil.CloseCallWithEnrollment(t, db.Database, call2)

	cmd := commands.DeleteCallCommand{ID: call2}
	err = cmd.Execute(db.Database)

	testutil.AssertErrorType(t, err, commands.ErrCannotDeleteClosedCall{})
	testutil.AssertCallStatus(t, db.Database, call2, types.CallStatusDone)
}

func TestDeleteCallCommand_RevertsCalledStudentsToWaitlist(t *testing.T) {
	db := database.NewTestDatabase(t)

	testutil.LoadApprovedSelection(t, db.Database, "testdata/approved_small.csv")
	testutil.LoadWaitlistSelection(t, db.Database, "testdata/waitlist_small.csv")

	call1, err := db.FetchCallByNumber(1)
	require.NoError(t, err)

	registrations := testutil.RegistrationsInCall(t, db.Database, call1.ID)
	for i, reg := range registrations {
		if i < 2 {
			testutil.EnrollRegistration(t, db.Database, reg.ID)
		} else {
			testutil.MarkRegistrationAbsent(t, db.Database, reg.ID)
		}
	}

	testutil.CloseCall(t, db.Database, call1.ID)
	call2ID := testutil.CreateCall(t, db.Database)

	called := testutil.RegistrationsInCall(t, db.Database, call2ID)
	require.NotEmpty(t, called)

	testutil.EnrollRegistration(t, db.Database, called[0].ID)

	cmd := commands.DeleteCallCommand{ID: call2ID}
	err = cmd.Execute(db.Database)
	require.NoError(t, err)

	_, err = db.FetchCallByID(call2ID)
	require.Error(t, err, "call 2 should be deleted")

	for _, reg := range called {
		testutil.AssertRegistrationStatus(t, db.Database, reg.ID, types.RegistrationStatusWaitlisted)
		testutil.AssertRegistrationSemester(t, db.Database, reg.ID, nil)
	}

	// Call 1 outcomes are untouched.
	testutil.AssertRegistrationStatus(t, db.Database, registrations[0].ID, types.RegistrationStatusEnrolled)
	testutil.AssertRegistrationStatus(t, db.Database, registrations[2].ID, types.RegistrationStatusAbsent)
}

func TestDeleteCallCommand_CannotDeleteNonLastCall(t *testing.T) {
	db := database.NewTestDatabase(t)

	testutil.LoadApprovedSelection(t, db.Database, "testdata/approved_split.csv")
	testutil.LoadWaitlistSelection(t, db.Database, "testdata/waitlist_split.csv")

	call1, err := db.FetchCallByNumber(1)
	require.NoError(t, err)

	// Leave one seat free in each semester so two more calls can be created.
	regs := testutil.RegistrationsInCall(t, db.Database, call1.ID)
	testutil.MarkRegistrationAbsent(t, db.Database, regs[0].ID)
	testutil.MarkRegistrationAbsent(t, db.Database, regs[3].ID)
	testutil.CloseCallWithEnrollment(t, db.Database, call1.ID)

	call2 := testutil.CreateCall(t, db.Database)
	call2Regs := testutil.RegistrationsInCall(t, db.Database, call2)
	testutil.MarkRegistrationAbsent(t, db.Database, call2Regs[0].ID)
	testutil.CloseCallWithEnrollment(t, db.Database, call2)

	call3 := testutil.CreateCall(t, db.Database)
	require.NotZero(t, call3)

	openCmd := commands.OpenCallCommand{ID: call2}
	testutil.AssertErrorType(t, openCmd.Execute(db.Database), commands.ErrCannotReopenCallWithLaterCalls{})

	cmd := commands.DeleteCallCommand{ID: call2}
	err = cmd.Execute(db.Database)

	testutil.AssertErrorType(t, err, commands.ErrCannotDeleteClosedCall{})

	calls, err := db.FetchCallSummaries()
	require.NoError(t, err)
	assert.Len(t, calls, 3)
}
