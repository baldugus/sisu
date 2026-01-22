package testutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/baldugus/sisu/database"
	"github.com/baldugus/sisu/types"
)

// AssertErrorType checks that the error is of the expected type using errors.As.
func AssertErrorType[T error](t *testing.T, err error, expectedErr T) {
	t.Helper()
	require.Error(t, err)
	assert.ErrorAs(t, err, &expectedErr)
}

// AssertRegistrationStatus verifies that a registration has the expected status.
func AssertRegistrationStatus(
	t *testing.T,
	db *database.Database,
	regID int32,
	expectedStatus types.RegistrationStatus,
) {
	t.Helper()

	detail, err := db.FetchRegistrationByID(regID)
	require.NoError(t, err)
	assert.Equal(t, expectedStatus, detail.Registration.Status,
		"registration %d should have status %s", regID, expectedStatus)
}

// AssertCallStatus verifies that a call has the expected status.
func AssertCallStatus(
	t *testing.T,
	db *database.Database,
	callID int32,
	expectedStatus types.CallStatus,
) {
	t.Helper()

	call, err := db.FetchCallByID(callID)
	require.NoError(t, err)
	assert.Equal(t, expectedStatus, call.Status,
		"call %d should have status %s", callID, expectedStatus)
}

// AssertCourseOccupiedSeats verifies the number of occupied seats in a course.
func AssertCourseOccupiedSeats(
	t *testing.T,
	db *database.Database,
	courseID int32,
	expectedOccupied int32,
) {
	t.Helper()

	occupied, err := db.CountCourseOccupiedSeats(courseID)
	require.NoError(t, err)
	assert.Equal(t, expectedOccupied, occupied,
		"course %d should have %d occupied seats", courseID, expectedOccupied)
}

// AssertDatabaseEmpty verifies that the database has no data in main tables.
func AssertDatabaseEmpty(t *testing.T, db *database.Database) {
	t.Helper()

	selections, err := db.FetchSelectionKinds()
	require.NoError(t, err)
	assert.Empty(t, selections, "selections table should be empty")

	registrations, err := db.FetchRegistrations()
	require.NoError(t, err)
	assert.Empty(t, registrations, "registrations table should be empty")

	calls, err := db.FetchCalls()
	require.NoError(t, err)
	assert.Empty(t, calls, "calls table should be empty")

	courses, err := db.FetchCourses()
	require.NoError(t, err)
	assert.Empty(t, courses, "courses table should be empty")
}

// AssertCallCreated verifies that a call with the given number exists.
func AssertCallCreated(t *testing.T, db *database.Database, callNumber int32) int32 {
	t.Helper()

	call, err := db.FetchCallByNumber(callNumber)
	require.NoError(t, err)
	assert.Equal(t, callNumber, call.Number)
	return call.ID
}

// AssertSelectionExists verifies that a selection of the given kind exists.
func AssertSelectionExists(t *testing.T, db *database.Database, kind types.SelectionKind) *types.Selection {
	t.Helper()

	selection, err := db.FetchSelection(kind)
	require.NoError(t, err)
	assert.Equal(t, kind, selection.Kind)
	return selection
}

// AssertRegistrationCount verifies the total number of registrations.
func AssertRegistrationCount(t *testing.T, db *database.Database, expected int) {
	t.Helper()

	registrations, err := db.FetchRegistrations()
	require.NoError(t, err)
	assert.Len(t, registrations, expected, "should have %d registrations", expected)
}

// AssertRegistrationsInCall verifies the number of registrations in a specific call.
func AssertRegistrationsInCall(t *testing.T, db *database.Database, callID int32, expected int) {
	t.Helper()

	registrations, err := db.FetchRegistrationsByCallID(callID)
	require.NoError(t, err)
	assert.Len(t, registrations, expected,
		"call %d should have %d registrations", callID, expected)
}
