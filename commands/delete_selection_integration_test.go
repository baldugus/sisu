package commands_test

import (
	"testing"

	"github.com/go-jet/jet/v2/qrm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/baldugus/sisu/commands"
	"github.com/baldugus/sisu/database"
	"github.com/baldugus/sisu/testutil"
	"github.com/baldugus/sisu/types"
)

func TestDeleteSelectionCommand_ApprovedOnly_DeletesEverything(t *testing.T) {
	db := database.NewTestDatabase(t)

	testutil.LoadApprovedSelection(t, db.Database, "testdata/approved_small.csv")

	testutil.AssertSelectionExists(t, db.Database, types.SelectionKindApproved)
	testutil.AssertRegistrationCount(t, db.Database, 5)

	cmd := commands.DeleteSelectionCommand{
		Kind: types.SelectionKindApproved,
	}
	err := cmd.Execute(db.Database)

	require.NoError(t, err)
	testutil.AssertDatabaseEmpty(t, db.Database)
}

func TestDeleteSelectionCommand_WaitlistOnly_PreservesApproved(t *testing.T) {
	db := database.NewTestDatabase(t)

	testutil.LoadApprovedSelection(t, db.Database, "testdata/approved_small.csv")
	testutil.LoadWaitlistSelection(t, db.Database, "testdata/waitlist_small.csv")

	approvedSelection := testutil.AssertSelectionExists(t, db.Database, types.SelectionKindApproved)
	waitlistSelection := testutil.AssertSelectionExists(t, db.Database, types.SelectionKindWaitlist)

	approvedRegs, err := db.FetchRegistrationsBySelectionID(approvedSelection.ID)
	require.NoError(t, err)
	approvedCount := len(approvedRegs)

	_, err = db.FetchRegistrationsBySelectionID(waitlistSelection.ID)
	require.NoError(t, err)

	cmd := commands.DeleteSelectionCommand{
		Kind: types.SelectionKindWaitlist,
	}
	err = cmd.Execute(db.Database)
	require.NoError(t, err)

	testutil.AssertSelectionExists(t, db.Database, types.SelectionKindApproved)

	_, err = db.FetchSelection(types.SelectionKindWaitlist)
	require.Error(t, err, "waitlist selection should be deleted")

	regsAfter, err := db.FetchRegistrations()
	require.NoError(t, err)
	assert.Len(t, regsAfter, approvedCount,
		"only approved registrations should remain (waitlist deleted)")

	for _, reg := range regsAfter {
		approvedRegs, err := db.FetchRegistrationsBySelectionID(approvedSelection.ID)
		require.NoError(t, err)

		found := false
		for _, approvedReg := range approvedRegs {
			if approvedReg.ID == reg.ID {
				found = true
				break
			}
		}
		assert.True(t, found, "registration %d should belong to approved selection", reg.ID)
	}

	courses, err := db.FetchCourses()
	require.NoError(t, err)
	assert.NotEmpty(t, courses, "courses should be preserved (shared resource)")
}

func TestDeleteSelectionCommand_CannotDeleteApprovedWithWaitlist(t *testing.T) {
	db := database.NewTestDatabase(t)

	testutil.LoadApprovedSelection(t, db.Database, "testdata/approved_small.csv")
	testutil.LoadWaitlistSelection(t, db.Database, "testdata/waitlist_small.csv")

	cmd := commands.DeleteSelectionCommand{
		Kind: types.SelectionKindApproved,
	}
	err := cmd.Execute(db.Database)

	testutil.AssertErrorType(t, err, commands.ErrCannotDeleteApprovedWithWaitlist{})
}

func TestDeleteSelectionCommand_CannotDeleteWithClosedFirstCall(t *testing.T) {
	db := database.NewTestDatabase(t)

	testutil.LoadApprovedSelection(t, db.Database, "testdata/approved_small.csv")

	call1, err := db.FetchCallByNumber(1)
	require.NoError(t, err)

	testutil.CloseCallWithEnrollment(t, db.Database, call1.ID)

	cmd := commands.DeleteSelectionCommand{
		Kind: types.SelectionKindApproved,
	}
	err = cmd.Execute(db.Database)

	testutil.AssertErrorType(t, err, commands.ErrCannotDeleteWithClosedFirstCall{})
}

func TestDeleteSelectionCommand_CannotDeleteWithMultipleCalls(t *testing.T) {
	db := database.NewTestDatabase(t)

	testutil.LoadApprovedSelection(t, db.Database, "testdata/approved_small.csv")

	call1, err := db.FetchCallByNumber(1)
	require.NoError(t, err)
	testutil.CloseCallWithEnrollment(t, db.Database, call1.ID)

	semester, err := database.FetchSemesterByYearAndNumber(db.Database.DB(), 2025, 1)
	require.NoError(t, err)

	err = db.RunInTx(func(tx qrm.DB) error {
		_, err := database.CreateCall(tx, &types.Call{
			Number:     2,
			Status:     types.CallStatusCalling,
			SemesterID: semester.ID,
		})
		return err
	})
	require.NoError(t, err)

	cmd := commands.DeleteSelectionCommand{
		Kind: types.SelectionKindApproved,
	}
	err = cmd.Execute(db.Database)

	testutil.AssertErrorType(t, err, commands.ErrCannotDeleteWithMultipleCalls{})
}
