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

func TestFetchSemesters_AlwaysTwo(t *testing.T) {
	db := database.NewTestDatabase(t)

	semesters, err := db.FetchSemesters()
	require.NoError(t, err)
	require.Len(t, semesters, 2)
	assert.Equal(t, int32(1), semesters[0].Number)
	assert.Equal(t, int32(2), semesters[1].Number)
	assert.Equal(t, types.SemesterStatusOpen, semesters[0].Status)
	assert.Equal(t, int32(0), semesters[0].Seats)
}

func TestCloseSemester_ClosedSemesterReceivesNoOne(t *testing.T) {
	db, call1, regs := splitCycle(t)
	waitlist := waitlistByRank(t, db.Database)

	testutil.SetOutcome(t, db.Database, call1, regs[0].ID, types.CallEntryOutcomeAbsent)
	testutil.SetOutcome(t, db.Database, call1, regs[2].ID, types.CallEntryOutcomeAbsent)
	testutil.SetWantsPromotion(t, db.Database, call1, regs[3].ID, true)
	testutil.CloseCallWithEnrollment(t, db.Database, call1)

	closeCmd := commands.CloseSemesterCommand{Number: 1}
	require.NoError(t, closeCmd.Execute(db.Database))

	semesters, err := db.FetchSemesters()
	require.NoError(t, err)
	assert.Equal(t, types.SemesterStatusClosed, semesters[0].Status)
	assert.Equal(t, ptr(1), semesters[0].ClosedAfterCall)

	call2 := testutil.CreateCall(t, db.Database)
	entries := testutil.CallEntries(t, db.Database, call2)
	assert.Empty(t, entryIDs(entries, types.CallEntryKindPromotion, 1), "no promotion into a closed semester")
	assert.Empty(t, entryIDs(entries, types.CallEntryKindWaitlist, 1))
	assert.Equal(t, []int32{waitlist[0].ID}, entryIDs(entries, types.CallEntryKindWaitlist, 2))
}

func TestCloseSemester_Rules(t *testing.T) {
	db, call1, _ := splitCycle(t)

	closeCmd := commands.CloseSemesterCommand{Number: 2}
	testutil.AssertErrorType(t, closeCmd.Execute(db.Database), commands.ErrCannotCloseSemesterWithOpenCall{})

	invalid := commands.CloseSemesterCommand{Number: 3}
	testutil.AssertErrorType(t, invalid.Execute(db.Database), commands.ErrInvalidSemester{})

	testutil.CloseCallWithEnrollment(t, db.Database, call1)
	require.NoError(t, closeCmd.Execute(db.Database))
	testutil.AssertErrorType(t, closeCmd.Execute(db.Database), commands.ErrSemesterAlreadyClosed{})

	both := commands.CloseSemesterCommand{Number: 1}
	require.NoError(t, both.Execute(db.Database))

	create := commands.CreateCallCommand{}
	testutil.AssertErrorType(t, create.Execute(db.Database), commands.ErrAllSemestersClosed{})
}

func TestReopenSemester_OnlyWithoutLaterCalls(t *testing.T) {
	db, call1, regs := splitCycle(t)

	testutil.SetOutcome(t, db.Database, call1, regs[2].ID, types.CallEntryOutcomeAbsent)
	testutil.CloseCallWithEnrollment(t, db.Database, call1)

	closeCmd := commands.CloseSemesterCommand{Number: 1}
	require.NoError(t, closeCmd.Execute(db.Database))

	testutil.CreateCall(t, db.Database)

	reopen := commands.ReopenSemesterCommand{Number: 1}
	testutil.AssertErrorType(t, reopen.Execute(db.Database), commands.ErrCannotReopenSemesterWithLaterCalls{})

	notClosed := commands.ReopenSemesterCommand{Number: 2}
	testutil.AssertErrorType(t, notClosed.Execute(db.Database), commands.ErrSemesterNotClosed{})
}

func TestDeleteApprovedSelection_ReopensSemesters(t *testing.T) {
	db := database.NewTestDatabase(t)
	testutil.LoadApprovedSelection(t, db.Database, "testdata/approved_promotion.csv")

	require.NoError(t, database.SetSemesterClosedAfterCall(db.DB(), 2, ptr(1)))

	testutil.DeleteApprovedSelection(t, db.Database)

	semesters, err := db.FetchSemesters()
	require.NoError(t, err)

	for _, s := range semesters {
		assert.Equal(t, types.SemesterStatusOpen, s.Status)
	}
}
