package commands_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/baldugus/sisu/database"
	"github.com/baldugus/sisu/testutil"
	"github.com/baldugus/sisu/types"
)

// TestCompleteAdmissionCycle tests the full workflow of the admission system
// from loading data to multiple enrollment calls.
func TestCompleteAdmissionCycle(t *testing.T) {
	// Arrange
	db := database.NewTestDatabase(t)

	// Step 1: Load approved selection
	t.Log("Step 1: Loading approved selection")
	testutil.LoadApprovedSelection(t, db.Database, "testdata/approved_small.csv")

	// Verify approved selection created
	approvedSelection := testutil.AssertSelectionExists(t, db.Database, types.SelectionKindApproved)
	assert.Equal(t, int32(2025), approvedSelection.Year)
	assert.Equal(t, int32(1), approvedSelection.Semester)

	// Verify call #1 created automatically
	call1 := testutil.AssertCallCreated(t, db.Database, 1)
	testutil.AssertCallStatus(t, db.Database, call1, types.CallStatusCalling)

	// Verify registrations created
	testutil.AssertRegistrationCount(t, db.Database, 5) // 5 students in approved_small.csv
	testutil.AssertRegistrationsInCall(t, db.Database, call1, 5)

	// Step 2: Load waitlist selection
	t.Log("Step 2: Loading waitlist selection")
	testutil.LoadWaitlistSelection(t, db.Database, "testdata/waitlist_small.csv")

	// Verify waitlist selection created
	testutil.AssertSelectionExists(t, db.Database, types.SelectionKindWaitlist)

	// Verify total registrations (5 approved + 3 waitlisted)
	testutil.AssertRegistrationCount(t, db.Database, 8)

	// Step 3: Process call #1 - enroll some, mark some absent
	t.Log("Step 3: Processing call #1")
	call1Regs, err := db.FetchRegistrationsByCallID(call1)
	require.NoError(t, err)
	require.Len(t, call1Regs, 5)

	// Enroll 3 students
	for i := 0; i < 3; i++ {
		testutil.EnrollRegistration(t, db.Database, call1Regs[i].ID)
	}

	// Mark 2 students absent
	for i := 3; i < 5; i++ {
		testutil.MarkRegistrationAbsent(t, db.Database, call1Regs[i].ID)
	}

	// Verify statuses
	testutil.AssertRegistrationStatus(t, db.Database, call1Regs[0].ID, types.RegistrationStatusEnrolled)
	testutil.AssertRegistrationStatus(t, db.Database, call1Regs[3].ID, types.RegistrationStatusAbsent)

	// Step 4: Close call #1
	t.Log("Step 4: Closing call #1")
	testutil.CloseCall(t, db.Database, call1)
	testutil.AssertCallStatus(t, db.Database, call1, types.CallStatusDone)

	// Step 5: Create call #2 (promotes waitlisted students)
	t.Log("Step 5: Creating call #2")
	call2 := testutil.CreateCall(t, db.Database)
	testutil.AssertCallStatus(t, db.Database, call2, types.CallStatusCalling)

	// Verify students were promoted to call #2
	call2Regs, err := db.FetchRegistrationsByCallID(call2)
	require.NoError(t, err)
	assert.Greater(t, len(call2Regs), 0, "should have promoted waitlisted students")

	// Verify promoted students have approved status
	for _, reg := range call2Regs {
		assert.Equal(t, types.RegistrationStatusApproved, reg.Status,
			"promoted student should have approved status")
	}

	// Step 6: Process call #2
	t.Log("Step 6: Processing call #2")
	if len(call2Regs) > 0 {
		// Enroll first promoted student
		testutil.EnrollRegistration(t, db.Database, call2Regs[0].ID)
		testutil.AssertRegistrationStatus(t, db.Database, call2Regs[0].ID, types.RegistrationStatusEnrolled)

		// Mark remaining as absent
		for i := 1; i < len(call2Regs); i++ {
			testutil.MarkRegistrationAbsent(t, db.Database, call2Regs[i].ID)
		}
	}

	// Step 7: Close call #2
	t.Log("Step 7: Closing call #2")
	testutil.CloseCall(t, db.Database, call2)
	testutil.AssertCallStatus(t, db.Database, call2, types.CallStatusDone)

	// Step 8: Verify final state
	t.Log("Step 8: Verifying final state")

	// Get all calls
	calls, err := db.FetchCalls()
	require.NoError(t, err)
	assert.Len(t, calls, 2, "should have 2 calls")

	// All calls should be closed
	for _, call := range calls {
		assert.Equal(t, types.CallStatusDone, call.Status,
			"all calls should be closed")
	}

	// Count enrolled students
	allRegs, err := db.FetchRegistrations()
	require.NoError(t, err)
	enrolledCount := 0
	absentCount := 0
	for _, reg := range allRegs {
		if reg.Status == types.RegistrationStatusEnrolled {
			enrolledCount++
		} else if reg.Status == types.RegistrationStatusAbsent {
			absentCount++
		}
	}

	assert.Equal(t, 4, enrolledCount, "should have 4 enrolled students")
	assert.Greater(t, absentCount, 0, "should have some absent students")

	t.Log("Admission cycle completed successfully!")
}

// TestAdmissionCycle_WithClearStatus tests the workflow including clearing statuses.
func TestAdmissionCycle_WithClearStatus(t *testing.T) {
	// Arrange
	db := database.NewTestDatabase(t)

	testutil.LoadApprovedSelection(t, db.Database, "testdata/approved_small.csv")

	call1, err := db.FetchCallByNumber(1)
	require.NoError(t, err)

	regs, err := db.FetchRegistrationsByCallID(call1.ID)
	require.NoError(t, err)
	reg := regs[0]

	// Enroll student
	testutil.EnrollRegistration(t, db.Database, reg.ID)
	testutil.AssertRegistrationStatus(t, db.Database, reg.ID, types.RegistrationStatusEnrolled)

	// Clear status back to approved
	testutil.ClearRegistrationStatus(t, db.Database, reg.ID)
	testutil.AssertRegistrationStatus(t, db.Database, reg.ID, types.RegistrationStatusApproved)

	// Enroll again
	testutil.EnrollRegistration(t, db.Database, reg.ID)
	testutil.AssertRegistrationStatus(t, db.Database, reg.ID, types.RegistrationStatusEnrolled)
}

// TestAdmissionCycle_MultipleCallsProgression tests creating multiple calls sequentially.
func TestAdmissionCycle_MultipleCallsProgression(t *testing.T) {
	// Arrange
	db := database.NewTestDatabase(t)

	testutil.LoadApprovedSelection(t, db.Database, "testdata/approved_small.csv")
	testutil.LoadWaitlistSelection(t, db.Database, "testdata/waitlist_small.csv")

	// Close call #1
	call1, err := db.FetchCallByNumber(1)
	require.NoError(t, err)
	testutil.CloseCallWithEnrollment(t, db.Database, call1.ID)

	// Create and close call #2
	call2 := testutil.CreateCall(t, db.Database)
	testutil.CloseCallWithEnrollment(t, db.Database, call2)

	// Verify call numbers are sequential
	calls, err := db.FetchCalls()
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(calls), 2, "should have at least 2 calls")

	// Verify call #1 and #2 exist with correct numbers
	foundCall1 := false
	foundCall2 := false
	for _, call := range calls {
		if call.Number == 1 {
			foundCall1 = true
		}
		if call.Number == 2 {
			foundCall2 = true
		}
	}
	assert.True(t, foundCall1, "call #1 should exist")
	assert.True(t, foundCall2, "call #2 should exist")

	// Both calls should be closed
	for _, call := range calls {
		assert.Equal(t, types.CallStatusDone, call.Status,
			"all calls should be closed")
	}
}
