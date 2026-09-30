package types_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/baldugus/sisu/types"
)

func TestNewSeats(t *testing.T) {
	tests := []struct {
		name            string
		total           int32
		wantPerSemester int32
		wantErr         bool
	}{
		{name: "even total", total: 20, wantPerSemester: 10},
		{name: "zero seats", total: 0, wantPerSemester: 0},
		{name: "odd total is an error", total: 15, wantErr: true},
		{name: "single seat is an error", total: 1, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			seats, err := types.NewSeats(tt.total)
			if tt.wantErr {
				var oddErr types.ErrOddSeatsCount
				require.ErrorAs(t, err, &oddErr)
				assert.Equal(t, tt.total, oddErr.Count)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.total, seats.Total())
			assert.Equal(t, tt.wantPerSemester, seats.PerSemester())
		})
	}
}

func TestSeatsMarshalJSON(t *testing.T) {
	seats, err := types.NewSeats(16)
	require.NoError(t, err)

	got, err := json.Marshal(types.Course{Seats: seats})
	require.NoError(t, err)

	assert.Contains(t, string(got), `"Seats":16`)
}
