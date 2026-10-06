package commands_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/baldugus/sisu/commands"
	"github.com/baldugus/sisu/testutil"
	"github.com/baldugus/sisu/types"
)

func TestSetCallEntryOutcome(t *testing.T) {
	tests := []struct {
		name    string
		outcome types.CallEntryOutcome
		want    types.RegistrationStatus
		wantErr error
	}{
		{name: "enrolled", outcome: types.CallEntryOutcomeEnrolled, want: types.RegistrationStatusEnrolled},
		{name: "absent", outcome: types.CallEntryOutcomeAbsent, want: types.RegistrationStatusAbsent},
		{name: "back to pending", outcome: types.CallEntryOutcomePending, want: types.RegistrationStatusApproved},
		{name: "unknown outcome", outcome: types.CallEntryOutcome("bogus"), wantErr: commands.ErrInvalidStatusTransition{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, call1, regs := splitCycle(t)

			cmd := commands.SetCallEntryOutcomeCommand{CallID: call1, RegistrationID: regs[0].ID, Outcome: tt.outcome}
			err := cmd.Execute(db.Database)

			if tt.wantErr != nil {
				testutil.AssertErrorType(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)
			testutil.AssertRegistrationStatus(t, db.Database, regs[0].ID, tt.want)
		})
	}
}

func TestSetCallEntryOutcome_ErrCallNotOpen(t *testing.T) {
	db, call1, regs := splitCycle(t)
	testutil.CloseCallWithEnrollment(t, db.Database, call1)

	cmd := commands.SetCallEntryOutcomeCommand{
		CallID:         call1,
		RegistrationID: regs[0].ID,
		Outcome:        types.CallEntryOutcomePending,
	}

	testutil.AssertErrorType(t, cmd.Execute(db.Database), commands.ErrCallNotOpen{})
}

func TestSetCallEntryOutcome_ErrRegistrationNotInCall(t *testing.T) {
	db, call1, _ := splitCycle(t)
	waitlist := waitlistByRank(t, db.Database)

	cmd := commands.SetCallEntryOutcomeCommand{
		CallID:         call1,
		RegistrationID: waitlist[0].ID,
		Outcome:        types.CallEntryOutcomeEnrolled,
	}

	testutil.AssertErrorType(t, cmd.Execute(db.Database), commands.ErrRegistrationNotFound{})
}

func TestCloseCall_ErrPendingEntries(t *testing.T) {
	db, call1, _ := splitCycle(t)

	cmd := commands.CloseCallCommand{ID: call1}

	testutil.AssertErrorType(t, cmd.Execute(db.Database), commands.ErrCallHasPendingRegistrations{})
}
