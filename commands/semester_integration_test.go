package commands_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/baldugus/sisu/database"
)

func TestFetchSemesters_AlwaysTwo(t *testing.T) {
	db := database.NewTestDatabase(t)

	semesters, err := db.FetchSemesters()
	require.NoError(t, err)
	require.Len(t, semesters, 2)
	assert.Equal(t, int32(1), semesters[0].Number)
	assert.Equal(t, int32(2), semesters[1].Number)
	assert.Equal(t, int32(0), semesters[0].Seats)
}
