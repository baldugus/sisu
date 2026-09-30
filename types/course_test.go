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
		wantErr      bool
	}{
		{name: "first ranking goes to semester 1", seats: 4, ranking: 1, wantSemester: 1},
		{name: "last ranking of first half goes to semester 1", seats: 4, ranking: 2, wantSemester: 1},
		{name: "first ranking of second half goes to semester 2", seats: 4, ranking: 3, wantSemester: 2},
		{name: "last seat goes to semester 2", seats: 4, ranking: 4, wantSemester: 2},
		{name: "ranking beyond seats goes to semester 2", seats: 4, ranking: 9, wantSemester: 2},
		{name: "odd seats is an error", seats: 5, ranking: 1, wantErr: true},
		{name: "single seat is an error", seats: 1, ranking: 1, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			course := types.Course{Seats: tt.seats}

			got, err := course.SemesterForRanking(tt.ranking)
			if tt.wantErr {
				var oddErr types.ErrOddSeatsCount
				require.ErrorAs(t, err, &oddErr)
				assert.Equal(t, tt.seats, oddErr.Count)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.wantSemester, got)
		})
	}
}
