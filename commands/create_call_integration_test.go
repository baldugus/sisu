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

func ptr(v int32) *int32 { return &v }

// splitCycle loads approved_promotion.csv (one course, 4 seats, ranks 1-4: ranks 1-2
// in semester 1, 3-4 in semester 2) and waitlist_promotion.csv (ranks 5-8), and
// returns call 1 and its registrations by rank (index 0 = rank 1).
func splitCycle(t *testing.T) (*database.TestDB, int32, []*types.Registration) {
	t.Helper()

	db := database.NewTestDatabase(t)
	testutil.LoadApprovedSelection(t, db.Database, "testdata/approved_promotion.csv")
	testutil.LoadWaitlistSelection(t, db.Database, "testdata/waitlist_promotion.csv")

	call1, err := db.FetchCallByNumber(1)
	require.NoError(t, err)

	regs := testutil.RegistrationsInCall(t, db.Database, call1.ID)
	require.Len(t, regs, 4)

	return db, call1.ID, regs
}

func waitlistByRank(t *testing.T, db *database.Database) []*types.Registration {
	t.Helper()

	selection := testutil.AssertSelectionExists(t, db, types.SelectionKindWaitlist)
	regs, err := db.FetchRegistrationsBySelectionID(selection.ID)
	require.NoError(t, err)
	require.Len(t, regs, 4)

	return regs
}

func entryIDs(entries []*types.CallEntryDetail, kind types.CallEntryKind, semester int32) []int32 {
	var ids []int32
	for _, e := range entries {
		if e.Entry.Kind == kind && e.Entry.Semester == semester {
			ids = append(ids, e.Entry.RegistrationID)
		}
	}

	return ids
}

func TestLoadApproved_SplitsSemestersInCallOne(t *testing.T) {
	db, call1, regs := splitCycle(t)

	entries := testutil.CallEntries(t, db.Database, call1)
	assert.Equal(t, []int32{regs[0].ID, regs[1].ID}, entryIDs(entries, types.CallEntryKindInitial, 1))
	assert.Equal(t, []int32{regs[2].ID, regs[3].ID}, entryIDs(entries, types.CallEntryKindInitial, 2))

	testutil.AssertRegistrationSemester(t, db.Database, regs[0].ID, ptr(1))
	testutil.AssertRegistrationSemester(t, db.Database, regs[3].ID, ptr(2))
	testutil.AssertSemesterOccupancy(t, db.Database, 2, 2)
}

func TestCreateCall_ErrOpenCallExists(t *testing.T) {
	db := database.NewTestDatabase(t)
	testutil.LoadApprovedSelection(t, db.Database, "testdata/approved_small.csv")

	cmd := commands.CreateCallCommand{}
	err := cmd.Execute(db.Database)

	testutil.AssertErrorType(t, err, commands.ErrOpenCallExists{})
}

func TestCreateCall_ErrNoCalls(t *testing.T) {
	db := database.NewTestDatabase(t)

	cmd := commands.CreateCallCommand{}
	err := cmd.Execute(db.Database)

	testutil.AssertErrorType(t, err, commands.ErrNoCalls{})
}

func TestCreateCall_ErrNoCandidatesToCall(t *testing.T) {
	db := database.NewTestDatabase(t)
	testutil.LoadApprovedSelection(t, db.Database, "testdata/approved_small.csv")

	call1, err := db.FetchCallByNumber(1)
	require.NoError(t, err)
	testutil.CloseCallWithEnrollment(t, db.Database, call1.ID)

	cmd := commands.CreateCallCommand{}
	err = cmd.Execute(db.Database)

	testutil.AssertErrorType(t, err, commands.ErrNoCandidatesToCall{})
}

func TestCreateCall_ErrAllCoursesFull(t *testing.T) {
	db, call1, _ := splitCycle(t)
	testutil.CloseCallWithEnrollment(t, db.Database, call1)

	cmd := commands.CreateCallCommand{}
	err := cmd.Execute(db.Database)

	testutil.AssertErrorType(t, err, commands.ErrAllCoursesFull{})
}

func TestCreateCall_WaitlistFillsSemesterOneFirst(t *testing.T) {
	db, call1, regs := splitCycle(t)
	waitlist := waitlistByRank(t, db.Database)

	// One absence in each semester opens one seat in each.
	testutil.SetOutcome(t, db.Database, call1, regs[0].ID, types.CallEntryOutcomeAbsent)
	testutil.SetOutcome(t, db.Database, call1, regs[2].ID, types.CallEntryOutcomeAbsent)
	testutil.CloseCallWithEnrollment(t, db.Database, call1)

	call2 := testutil.CreateCall(t, db.Database)

	entries := testutil.CallEntries(t, db.Database, call2)
	assert.Equal(t, []int32{waitlist[0].ID}, entryIDs(entries, types.CallEntryKindWaitlist, 1), "best ranked goes to semester 1")
	assert.Equal(t, []int32{waitlist[1].ID}, entryIDs(entries, types.CallEntryKindWaitlist, 2))
	assert.Empty(t, entryIDs(entries, types.CallEntryKindPromotion, 1))

	testutil.AssertRegistrationStatus(t, db.Database, waitlist[0].ID, types.RegistrationStatusApproved)
	testutil.AssertRegistrationSemester(t, db.Database, waitlist[0].ID, ptr(1))
	testutil.AssertRegistrationStatus(t, db.Database, waitlist[2].ID, types.RegistrationStatusWaitlisted)
	testutil.AssertSemesterOccupancy(t, db.Database, 2, 2)
}

func TestCreateCall_PreviewMatchesCreatedCall(t *testing.T) {
	db, call1, regs := splitCycle(t)

	testutil.SetOutcome(t, db.Database, call1, regs[1].ID, types.CallEntryOutcomeAbsent)
	testutil.SetWantsPromotion(t, db.Database, call1, regs[2].ID, true)
	testutil.CloseCallWithEnrollment(t, db.Database, call1)

	preview := commands.PreviewCallCommand{}
	plan, err := preview.Execute(db.Database)
	require.NoError(t, err)

	assert.Equal(t, int32(2), plan.Number)
	require.Len(t, plan.Courses, 1)
	assert.Equal(t, int32(1), plan.Courses[0].Vacancies1)
	assert.Equal(t, int32(0), plan.Courses[0].Vacancies2)
	require.Len(t, plan.Courses[0].Promoted, 1)
	assert.Equal(t, regs[2].ID, plan.Courses[0].Promoted[0].ID)
	assert.Empty(t, plan.Courses[0].Waitlist1)
	assert.Equal(t, int32(1), plan.Total())

	call2 := testutil.CreateCall(t, db.Database)
	entries := testutil.CallEntries(t, db.Database, call2)
	require.Len(t, entries, 1)
	assert.Equal(t, types.CallEntryKindPromotion, entries[0].Entry.Kind)
	assert.Equal(t, regs[2].ID, entries[0].Entry.RegistrationID)
}

// TestPromotionFlow walks the promotion rules: a semester-2 student who asked to
// move up is offered semester 1 before the waitlist; while the offer is pending
// or declined they stay in semester 2; accepting moves them and frees their
// semester-2 seat for the next call.
func TestPromotionFlow(t *testing.T) {
	db, call1, regs := splitCycle(t)
	waitlist := waitlistByRank(t, db.Database)

	// Call 1: rank 1 (sem 1) absent; ranks 3 and 4 (sem 2) both want semester 1.
	testutil.SetOutcome(t, db.Database, call1, regs[0].ID, types.CallEntryOutcomeAbsent)
	testutil.SetOutcome(t, db.Database, call1, regs[2].ID, types.CallEntryOutcomeEnrolled)
	testutil.SetOutcome(t, db.Database, call1, regs[3].ID, types.CallEntryOutcomeEnrolled)
	testutil.SetWantsPromotion(t, db.Database, call1, regs[2].ID, true)
	testutil.SetWantsPromotion(t, db.Database, call1, regs[3].ID, true)
	testutil.CloseCallWithEnrollment(t, db.Database, call1)

	// Call 2: the single semester-1 seat goes to the best ranked requester.
	call2 := testutil.CreateCall(t, db.Database)
	entries := testutil.CallEntries(t, db.Database, call2)
	assert.Equal(t, []int32{regs[2].ID}, entryIDs(entries, types.CallEntryKindPromotion, 1))
	assert.Empty(t, entryIDs(entries, types.CallEntryKindWaitlist, 1), "promotion took the only seat")
	assert.Empty(t, entryIDs(entries, types.CallEntryKindWaitlist, 2), "semester 2 is full until the promotion is accepted")

	// A pending offer leaves the student enrolled in semester 2.
	testutil.AssertRegistrationStatus(t, db.Database, regs[2].ID, types.RegistrationStatusEnrolled)
	testutil.AssertRegistrationSemester(t, db.Database, regs[2].ID, ptr(2))

	// Declined: still semester 2, and never offered again.
	testutil.SetOutcome(t, db.Database, call2, regs[2].ID, types.CallEntryOutcomeDeclined)
	testutil.AssertRegistrationSemester(t, db.Database, regs[2].ID, ptr(2))
	testutil.AssertSemesterOccupancy(t, db.Database, 1, 2)
	testutil.CloseCall(t, db.Database, call2)

	// Call 3: the next requester gets the seat.
	call3 := testutil.CreateCall(t, db.Database)
	entries = testutil.CallEntries(t, db.Database, call3)
	assert.Equal(t, []int32{regs[3].ID}, entryIDs(entries, types.CallEntryKindPromotion, 1))

	// Accepted: moves to semester 1, freeing a semester-2 seat.
	testutil.SetOutcome(t, db.Database, call3, regs[3].ID, types.CallEntryOutcomeEnrolled)
	testutil.AssertRegistrationStatus(t, db.Database, regs[3].ID, types.RegistrationStatusEnrolled)
	testutil.AssertRegistrationSemester(t, db.Database, regs[3].ID, ptr(1))
	testutil.AssertSemesterOccupancy(t, db.Database, 2, 1)
	testutil.CloseCall(t, db.Database, call3)

	// Call 4: the freed semester-2 seat goes to the waitlist.
	call4 := testutil.CreateCall(t, db.Database)
	entries = testutil.CallEntries(t, db.Database, call4)
	assert.Empty(t, entryIDs(entries, types.CallEntryKindPromotion, 1))
	assert.Equal(t, []int32{waitlist[0].ID}, entryIDs(entries, types.CallEntryKindWaitlist, 2))

	// History shows every step for the promoted student.
	detail, err := db.FetchRegistrationByID(regs[3].ID)
	require.NoError(t, err)
	require.Len(t, detail.History, 2)
	assert.Equal(t, int32(1), detail.History[0].CallNumber)
	assert.Equal(t, types.CallEntryKindInitial, detail.History[0].Kind)
	assert.True(t, detail.History[0].WantsPromotion)
	assert.Equal(t, int32(3), detail.History[1].CallNumber)
	assert.Equal(t, types.CallEntryKindPromotion, detail.History[1].Kind)
}

// TestRegistrationHistory checks that a registration's history lists its
// entries oldest call first.
func TestRegistrationHistory(t *testing.T) {
	db, call1, regs := splitCycle(t)
	waitlist := waitlistByRank(t, db.Database)

	testutil.SetOutcome(t, db.Database, call1, regs[0].ID, types.CallEntryOutcomeAbsent)
	testutil.CloseCallWithEnrollment(t, db.Database, call1)

	call2 := testutil.CreateCall(t, db.Database)
	testutil.SetOutcome(t, db.Database, call2, waitlist[0].ID, types.CallEntryOutcomeEnrolled)

	detail, err := db.FetchRegistrationByID(regs[0].ID)
	require.NoError(t, err)
	require.Len(t, detail.History, 1)
	assert.Equal(t, int32(1), detail.History[0].CallNumber)
	assert.Equal(t, types.CallEntryKindInitial, detail.History[0].Kind)
	assert.Equal(t, types.CallEntryOutcomeAbsent, detail.History[0].Outcome)

	detail, err = db.FetchRegistrationByID(waitlist[0].ID)
	require.NoError(t, err)
	require.Len(t, detail.History, 1)
	assert.Equal(t, int32(2), detail.History[0].CallNumber)
	assert.Equal(t, types.CallEntryKindWaitlist, detail.History[0].Kind)
	assert.Equal(t, int32(1), detail.History[0].Semester)
	assert.Equal(t, types.CallEntryOutcomeEnrolled, detail.History[0].Outcome)
}

// TestUndoToImport deletes and reopens calls in reverse order and checks that
// the cycle returns exactly to the state right after the import.
func TestUndoToImport(t *testing.T) {
	db, call1, regs := splitCycle(t)

	testutil.SetOutcome(t, db.Database, call1, regs[0].ID, types.CallEntryOutcomeAbsent)
	testutil.SetOutcome(t, db.Database, call1, regs[2].ID, types.CallEntryOutcomeEnrolled)
	testutil.SetWantsPromotion(t, db.Database, call1, regs[2].ID, true)
	testutil.CloseCallWithEnrollment(t, db.Database, call1)

	call2 := testutil.CreateCall(t, db.Database)
	testutil.CloseCallWithEnrollment(t, db.Database, call2) // accepts the promotion
	testutil.AssertRegistrationSemester(t, db.Database, regs[2].ID, ptr(1))

	closeSem := commands.CloseSemesterCommand{Number: 1}
	require.NoError(t, closeSem.Execute(db.Database))

	// Undo, most recent action first.
	reopenCall := commands.OpenCallCommand{ID: call2}
	testutil.AssertErrorType(t, reopenCall.Execute(db.Database), commands.ErrSemesterClosedAfterCall{})

	reopenSem := commands.ReopenSemesterCommand{Number: 1}
	require.NoError(t, reopenSem.Execute(db.Database))
	require.NoError(t, reopenCall.Execute(db.Database))

	deleteCall := commands.DeleteCallCommand{ID: call2}
	require.NoError(t, deleteCall.Execute(db.Database))
	testutil.AssertRegistrationSemester(t, db.Database, regs[2].ID, ptr(2))

	reopenCall1 := commands.OpenCallCommand{ID: call1}
	require.NoError(t, reopenCall1.Execute(db.Database))

	for _, reg := range regs {
		testutil.ClearRegistrationStatus(t, db.Database, reg.ID)
	}

	testutil.SetWantsPromotion(t, db.Database, call1, regs[2].ID, false)

	// Back to "just imported": the approved list can be deleted again.
	testutil.DeleteWaitlistSelection(t, db.Database)
	testutil.DeleteApprovedSelection(t, db.Database)
	testutil.AssertDatabaseEmpty(t, db.Database)
}
