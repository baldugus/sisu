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

func TestLoadApprovedSelection_Success(t *testing.T) {
	// Arrange
	db := database.NewTestDatabase(t)

	cmd := commands.LoadSelectionCommand{
		Year:     2025,
		FilePath: "testdata/approved_small.csv",
		Kind:     types.SelectionKindApproved,
	}

	// Act
	err := cmd.Execute(db.Database)

	// Assert
	require.NoError(t, err)

	// Verify selection created
	selections, err := db.FetchSelectionKinds()
	require.NoError(t, err)
	assert.Contains(t, selections, types.SelectionKindApproved)

	selection := testutil.AssertSelectionExists(t, db.Database, types.SelectionKindApproved)
	assert.Equal(t, int32(2025), selection.Year)

	// Verify first call created with "calling" status
	call, err := db.FetchCallByNumber(1)
	require.NoError(t, err)
	assert.Equal(t, types.CallStatusCalling, call.Status)

	// Verify registrations created (5 students in approved_small.csv)
	testutil.AssertRegistrationCount(t, db.Database, 5)

	// Verify all registrations are "approved" and in call #1
	registrations, err := db.FetchRegistrations()
	require.NoError(t, err)
	for _, reg := range registrations {
		assert.Equal(t, types.RegistrationStatusApproved, reg.Status)
	}

	// Verify registrations are in call #1
	testutil.AssertRegistrationsInCall(t, db.Database, call.ID, 5)

	// Verify courses were created (2 courses: Matutino and Noturno)
	courses, err := db.FetchCourses()
	require.NoError(t, err)
	assert.Len(t, courses, 2, "should have created 2 courses (morning and evening)")
}

func TestLoadWaitlistSelection_Success(t *testing.T) {
	// Arrange
	db := database.NewTestDatabase(t)

	// First load approved selection
	testutil.LoadApprovedSelection(t, db.Database, "testdata/approved_small.csv")

	cmd := commands.LoadSelectionCommand{
		Year:     2025,
		FilePath: "testdata/waitlist_small.csv",
		Kind:     types.SelectionKindWaitlist,
	}

	// Act
	err := cmd.Execute(db.Database)

	// Assert
	require.NoError(t, err)

	// Verify both selections exist
	selections, err := db.FetchSelectionKinds()
	require.NoError(t, err)
	assert.Contains(t, selections, types.SelectionKindApproved)
	assert.Contains(t, selections, types.SelectionKindWaitlist)

	// Verify total registrations (5 approved + 3 waitlisted = 8)
	testutil.AssertRegistrationCount(t, db.Database, 8)

	// Verify waitlisted registrations have correct status
	waitlistSelection := testutil.AssertSelectionExists(t, db.Database, types.SelectionKindWaitlist)
	waitlistRegs, err := db.FetchRegistrationsBySelectionID(waitlistSelection.ID)
	require.NoError(t, err)
	assert.Len(t, waitlistRegs, 3, "should have 3 waitlisted registrations")

	for _, reg := range waitlistRegs {
		assert.Equal(t, types.RegistrationStatusWaitlisted, reg.Status)
	}
}

func TestLoadSelection_BusinessRules(t *testing.T) {
	tests := []struct {
		name          string
		setup         func(t *testing.T, db *database.Database)
		cmd           commands.LoadSelectionCommand
		expectedError any
	}{
		{
			name: "cannot load approved when one exists",
			setup: func(t *testing.T, db *database.Database) {
				testutil.LoadApprovedSelection(t, db, "testdata/approved_small.csv")
			},
			cmd: commands.LoadSelectionCommand{
				Year:     2025,
				FilePath: "testdata/approved_small.csv",
				Kind:     types.SelectionKindApproved,
			},
			expectedError: &commands.ErrApprovedSelectionAlreadyExists{},
		},
		{
			name: "cannot load waitlist before approved",
			setup: func(t *testing.T, db *database.Database) {
				// No setup - empty database
			},
			cmd: commands.LoadSelectionCommand{
				Year:     2025,
				FilePath: "testdata/waitlist_small.csv",
				Kind:     types.SelectionKindWaitlist,
			},
			expectedError: &commands.ErrWaitlistSelectionRequiresApproved{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := database.NewTestDatabase(t)
			tt.setup(t, db.Database)

			err := tt.cmd.Execute(db.Database)

			require.Error(t, err)
			assert.ErrorAs(t, err, tt.expectedError)
		})
	}
}

func TestLoadSelection_TransactionRollback(t *testing.T) {
	// Arrange
	db := database.NewTestDatabase(t)

	cmd := commands.LoadSelectionCommand{
		Year:     2025,
		FilePath: "testdata/invalid_missing_fields.csv",
		Kind:     types.SelectionKindApproved,
	}

	// Act
	err := cmd.Execute(db.Database)

	// Assert
	require.Error(t, err, "should fail with invalid CSV")

	// Verify transaction rolled back - no data persisted
	testutil.AssertDatabaseEmpty(t, db.Database)
}

func TestLoadSelection_CreatesCourses(t *testing.T) {
	// Arrange
	db := database.NewTestDatabase(t)

	// Act
	testutil.LoadApprovedSelection(t, db.Database, "testdata/approved_small.csv")

	// Assert
	courses, err := db.FetchCourses()
	require.NoError(t, err)

	// Should have 2 courses (Matutino and Noturno)
	assert.Len(t, courses, 2)

	// Verify course properties
	morningCourse := findCourseByPeriod(t, courses, types.CoursePeriodMorning)
	assert.Equal(t, int32(10), morningCourse.Seats)

	eveningCourse := findCourseByPeriod(t, courses, types.CoursePeriodEvening)
	assert.Equal(t, int32(16), eveningCourse.Seats)
}

func TestLoadSelection_CreatesCandidates(t *testing.T) {
	// Arrange
	db := database.NewTestDatabase(t)

	// Act
	testutil.LoadApprovedSelection(t, db.Database, "testdata/approved_small.csv")

	// Assert
	registrations, err := db.FetchRegistrations()
	require.NoError(t, err)

	// Verify each registration has candidate data
	for _, reg := range registrations {
		assert.NotNil(t, reg.Candidate)
		assert.NotEmpty(t, reg.Candidate.CPF)
		assert.NotEmpty(t, reg.Candidate.Name)
		assert.NotEmpty(t, reg.Candidate.Email)
	}
}

func TestLoadSelection_MaintainsRanking(t *testing.T) {
	// Arrange
	db := database.NewTestDatabase(t)

	// Act
	testutil.LoadApprovedSelection(t, db.Database, "testdata/approved_small.csv")

	// Assert
	registrations, err := db.FetchRegistrations()
	require.NoError(t, err)

	// Find registrations by enrollment ID and verify ranking
	rankings := make(map[string]int32)
	for _, reg := range registrations {
		rankings[reg.EnrollmentID] = reg.Ranking
	}

	// Verify specific rankings from CSV
	assert.Equal(t, int32(1), rankings["241077018570"], "João Silva should be ranked 1")
	assert.Equal(t, int32(2), rankings["241071358853"], "Maria Santos should be ranked 2")
	assert.Equal(t, int32(3), rankings["241080353259"], "Carlos Ferreira should be ranked 3")
}

// Helper function to find course by period
func findCourseByPeriod(t *testing.T, courses []*types.Course, period types.CoursePeriod) *types.Course {
	t.Helper()
	for _, course := range courses {
		if course.Period == period {
			return course
		}
	}
	t.Fatalf("course with period %s not found", period)
	return nil
}
