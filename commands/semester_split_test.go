package commands_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/baldugus/sisu/commands"
	"github.com/baldugus/sisu/types"
)

func TestSplitApprovedBySemester(t *testing.T) {
	ranked := func(rankings ...int32) []*types.Registration {
		regs := make([]*types.Registration, len(rankings))
		for i, r := range rankings {
			regs[i] = &types.Registration{Ranking: r}
		}
		return regs
	}
	rankingsOf := func(regs []*types.Registration) []int32 {
		out := []int32{}
		for _, r := range regs {
			out = append(out, r.Ranking)
		}
		return out
	}

	tests := []struct {
		name     string
		seats    int32
		input    []*types.Registration
		wantSem1 []int32
		wantSem2 []int32
		wantErr  bool
	}{
		{
			name:     "even seats split top half to semester 1",
			seats:    4,
			input:    ranked(1, 2, 3, 4),
			wantSem1: []int32{1, 2},
			wantSem2: []int32{3, 4},
		},
		{
			name:     "fewer candidates than seats",
			seats:    10,
			input:    ranked(1, 2, 6),
			wantSem1: []int32{1, 2},
			wantSem2: []int32{6},
		},
		{
			name:     "ranking beyond seats goes to semester 2",
			seats:    2,
			input:    ranked(1, 5),
			wantSem1: []int32{1},
			wantSem2: []int32{5},
		},
		{
			name:     "input order preserved",
			seats:    4,
			input:    ranked(4, 1, 3, 2),
			wantSem1: []int32{1, 2},
			wantSem2: []int32{4, 3},
		},
		{
			name:     "empty input",
			seats:    4,
			input:    nil,
			wantSem1: []int32{},
			wantSem2: []int32{},
		},
		{
			name:    "odd seats is an error",
			seats:   5,
			input:   ranked(1, 2, 3, 4, 5),
			wantErr: true,
		},
		{
			name:    "single seat is an error",
			seats:   1,
			input:   ranked(1),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sem1, sem2, err := commands.SplitApprovedBySemester(tt.seats, tt.input)
			if tt.wantErr {
				var oddErr commands.ErrOddSeatsCount
				require.ErrorAs(t, err, &oddErr)
				assert.Equal(t, tt.seats, oddErr.Count)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.wantSem1, rankingsOf(sem1))
			assert.Equal(t, tt.wantSem2, rankingsOf(sem2))
		})
	}
}
