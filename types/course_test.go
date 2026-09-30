package types_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/baldugus/sisu/types"
)

func TestCourseSemesterForRanking(t *testing.T) {
	tests := []struct {
		name         string
		seats        int32
		ranking      int32
		wantSemester int32
	}{
		{name: "first ranking goes to semester 1", seats: 4, ranking: 1, wantSemester: 1},
		{name: "last ranking of first half goes to semester 1", seats: 4, ranking: 2, wantSemester: 1},
		{name: "first ranking of second half goes to semester 2", seats: 4, ranking: 3, wantSemester: 2},
		{name: "last seat goes to semester 2", seats: 4, ranking: 4, wantSemester: 2},
		{name: "ranking beyond seats goes to semester 2", seats: 4, ranking: 9, wantSemester: 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			seats, err := types.NewSeats(tt.seats)
			require.NoError(t, err)

			course := types.Course{Seats: seats}

			assert.Equal(t, tt.wantSemester, course.SemesterForRanking(tt.ranking))
		})
	}
}
