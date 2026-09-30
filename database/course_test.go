package database_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/baldugus/sisu/database"
)

// The courses_seats_even triggers keep odd seat counts out of the database even
// when a row is written without going through types.NewSeats.
func TestCoursesSeatsEvenTriggers(t *testing.T) {
	tests := []struct {
		name    string
		stmt    string
		wantErr bool
	}{
		{
			name: "insert with even seats",
			stmt: `INSERT INTO courses (time_slot, quota_id, seats, minimum_score) VALUES ('evening', 1, 16, 0)`,
		},
		{
			name:    "insert with odd seats",
			stmt:    `INSERT INTO courses (time_slot, quota_id, seats, minimum_score) VALUES ('evening', 1, 15, 0)`,
			wantErr: true,
		},
		{
			name: "update to even seats",
			stmt: `UPDATE courses SET seats = 12 WHERE id = 1`,
		},
		{
			name:    "update to odd seats",
			stmt:    `UPDATE courses SET seats = 11 WHERE id = 1`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := database.NewTestDatabase(t)
			ctx := context.Background()

			_, err := db.DB().ExecContext(ctx, `INSERT INTO quotas (id, name) VALUES (1, 'Ampla concorrência')`)
			require.NoError(t, err)
			_, err = db.DB().ExecContext(ctx, `INSERT INTO courses (id, time_slot, quota_id, seats, minimum_score) VALUES (1, 'morning', 1, 10, 0)`)
			require.NoError(t, err)

			_, err = db.DB().ExecContext(ctx, tt.stmt)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "courses.seats must be even")
				return
			}

			require.NoError(t, err)
		})
	}
}
