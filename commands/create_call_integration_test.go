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

func TestCreateCall_Success(t *testing.T) {
	// Arrange
	db := database.NewTestDatabase(t)

	// Load approved and waitlist selections
	testutil.LoadApprovedSelection(t, db.Database, "testdata/approved_small.csv")
	testutil.LoadWaitlistSelection(t, db.Database, "testdata/waitlist_small.csv")

	// Close first call so we can create a new one
	call1, err := db.FetchCallByNumber(1)
	require.NoError(t, err)
	testutil.CloseCallWithEnrollment(t, db.Database, call1.ID)

	semester, _ := database.FetchSemesterByYearAndNumber(db.Database.DB(), 2025, 1)
	cmd := commands.CreateCallCommand{SemesterID: semester.ID}

	// Act
	err = cmd.Execute(db.Database)

	// Assert
	require.NoError(t, err)

	// Verify call #2 was created
	call2 := testutil.AssertCallCreated(t, db.Database, 2)
	testutil.AssertCallStatus(t, db.Database, call2, types.CallStatusCalling)

	// Verify waitlisted students were promoted
	registrationsInCall2, err := db.FetchRegistrationsByCallID(call2)
	require.NoError(t, err)
	assert.Greater(t, len(registrationsInCall2), 0, "should have promoted some students")
}

func TestCreateCall_ErrOpenCallExists(t *testing.T) {
	// Arrange
	db := database.NewTestDatabase(t)

	// Load approved selection (creates call #1 with status "calling")
	testutil.LoadApprovedSelection(t, db.Database, "testdata/approved_small.csv")

	semester, _ := database.FetchSemesterByYearAndNumber(db.Database.DB(), 2025, 1)
	cmd := commands.CreateCallCommand{SemesterID: semester.ID}

	// Act
	err := cmd.Execute(db.Database)

	// Assert
	require.Error(t, err)
	assert.ErrorAs(t, err, &commands.ErrOpenCallExists{})
}

func TestCreateCall_ErrNoWaitlistedRegistrations(t *testing.T) {
	// Arrange
	db := database.NewTestDatabase(t)

	// Load only approved selection (no waitlist)
	testutil.LoadApprovedSelection(t, db.Database, "testdata/approved_small.csv")

	// Close first call
	call1, err := db.FetchCallByNumber(1)
	require.NoError(t, err)
	testutil.CloseCallWithEnrollment(t, db.Database, call1.ID)

	semester, _ := database.FetchSemesterByYearAndNumber(db.Database.DB(), 2025, 1)
	cmd := commands.CreateCallCommand{SemesterID: semester.ID}

	// Act
	err = cmd.Execute(db.Database)

	// Assert
	require.Error(t, err)
	assert.ErrorAs(t, err, &commands.ErrNoWaitlistedRegistrations{})
}

// TestCreateCall_ErrAllCoursesFull tests the scenario where all courses are at capacity.
// Note: This test is skipped because with our small test data (5 approved, 10+15 seats),
// we can't easily fill all seats. The error case is tested in unit tests.
func TestCreateCall_ErrAllCoursesFull(t *testing.T) {
	t.Skip("Test data doesn't support filling all course seats - error case covered elsewhere")
}

func TestCreateCall_PromotesCorrectNumberOfStudents(t *testing.T) {
	// Arrange
	db := database.NewTestDatabase(t)

	// Load approved and waitlist selections
	testutil.LoadApprovedSelection(t, db.Database, "testdata/approved_small.csv")
	testutil.LoadWaitlistSelection(t, db.Database, "testdata/waitlist_small.csv")

	// Close first call
	call1, err := db.FetchCallByNumber(1)
	require.NoError(t, err)
	testutil.CloseCallWithEnrollment(t, db.Database, call1.ID)

	// Count available seats
	courses, err := db.FetchCourses()
	require.NoError(t, err)

	totalAvailableSeats := int32(0)
	for _, course := range courses {
		occupied, err := database.CountCourseOccupiedSeats(db.DB(), course.ID)
		require.NoError(t, err)
		available := course.Seats - occupied
		if available > 0 {
			totalAvailableSeats += available
		}
	}

	// Count waitlisted students
	waitlistSelection := testutil.AssertSelectionExists(t, db.Database, types.SelectionKindWaitlist)
	waitlistRegs, err := db.FetchRegistrationsBySelectionID(waitlistSelection.ID)
	require.NoError(t, err)
	waitlistCount := int32(len(waitlistRegs))

	semester, _ := database.FetchSemesterByYearAndNumber(db.Database.DB(), 2025, 1)
	cmd := commands.CreateCallCommand{SemesterID: semester.ID}

	// Act
	err = cmd.Execute(db.Database)

	// Assert
	require.NoError(t, err)

	// Verify promoted count
	call2 := testutil.AssertCallCreated(t, db.Database, 2)
	promotedRegs, err := db.FetchRegistrationsByCallID(call2)
	require.NoError(t, err)

	// Should promote min(waitlist count, available seats)
	expectedPromoted := waitlistCount
	if totalAvailableSeats < waitlistCount {
		expectedPromoted = totalAvailableSeats
	}

	assert.Equal(t, int(expectedPromoted), len(promotedRegs),
		"should promote correct number of students")
}

func TestCreateCall_IncrementCallNumber(t *testing.T) {
	// Arrange
	db := database.NewTestDatabase(t)

	// Load approved and waitlist selections
	testutil.LoadApprovedSelection(t, db.Database, "testdata/approved_small.csv")
	testutil.LoadWaitlistSelection(t, db.Database, "testdata/waitlist_small.csv")

	// Close first call
	call1, err := db.FetchCallByNumber(1)
	require.NoError(t, err)
	testutil.CloseCallWithEnrollment(t, db.Database, call1.ID)

	semester, _ := database.FetchSemesterByYearAndNumber(db.Database.DB(), 2025, 1)
	cmd := commands.CreateCallCommand{SemesterID: semester.ID}

	// Act
	err = cmd.Execute(db.Database)

	// Assert
	require.NoError(t, err)

	// Verify call number is 2
	call2, err := db.FetchCallByNumber(2)
	require.NoError(t, err)
	assert.Equal(t, int32(2), call2.Number)
}

func TestCreateCall_PromotedStudentsHaveApprovedStatus(t *testing.T) {
	// Arrange
	db := database.NewTestDatabase(t)

	// Load approved and waitlist selections
	testutil.LoadApprovedSelection(t, db.Database, "testdata/approved_small.csv")
	testutil.LoadWaitlistSelection(t, db.Database, "testdata/waitlist_small.csv")

	// Close first call
	call1, err := db.FetchCallByNumber(1)
	require.NoError(t, err)
	testutil.CloseCallWithEnrollment(t, db.Database, call1.ID)

	semester, _ := database.FetchSemesterByYearAndNumber(db.Database.DB(), 2025, 1)
	cmd := commands.CreateCallCommand{SemesterID: semester.ID}

	// Act
	err = cmd.Execute(db.Database)

	// Assert
	require.NoError(t, err)

	// Verify promoted students have "approved" status (not waitlisted anymore)
	call2 := testutil.AssertCallCreated(t, db.Database, 2)
	promotedRegs, err := db.FetchRegistrationsByCallID(call2)
	require.NoError(t, err)

	for _, reg := range promotedRegs {
		assert.Equal(t, types.RegistrationStatusApproved, reg.Status,
			"promoted registration %d should have approved status", reg.ID)
	}
}

