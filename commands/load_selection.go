package commands

import (
	"fmt"
	"slices"

	"github.com/baldugus/sisu/csvparser"
	"github.com/baldugus/sisu/database"
	"github.com/baldugus/sisu/types"
	"github.com/go-jet/jet/v2/qrm"
)

type LoadSelectionCommand struct {
	Year     int32
	FilePath string
	Kind     types.SelectionKind
}

func (cmd *LoadSelectionCommand) Execute(db *database.Database) error {
	existingSelectionKinds, err := db.FetchSelectionKinds()
	if err != nil {
		return fmt.Errorf("fetch selection kind: %w", err)
	}

	if cmd.Kind == types.SelectionKindApproved && len(existingSelectionKinds) > 0 {
		return ErrApprovedSelectionAlreadyExists{}
	}

	if cmd.Kind == types.SelectionKindWaitlist && !slices.Contains(existingSelectionKinds, types.SelectionKindApproved) {
		return ErrWaitlistSelectionRequiresApproved{}
	}

	parsedCsv, err := csvparser.ParseFile(cmd.FilePath)
	if err != nil {
		return fmt.Errorf("parse csv: %w", err)
	}

	parsed, err := parsedCsv.ToSelectionDomain(cmd.Kind, cmd.Year)
	if err != nil {
		return fmt.Errorf("parsing selection: %w", err)
	}

	err = db.RunInTx(func(tx qrm.DB) error {
		selectionID, err := database.CreateSelection(tx, parsed.Selection)
		if err != nil {
			return err
		}

		// The approved import opens call 1 with every approved student, already
		// split between the two semesters by the parser.
		var callID int32
		if cmd.Kind == types.SelectionKindApproved {
			callID, err = database.CreateCall(tx, &types.Call{
				Number: 1,
				Status: types.CallStatusCalling,
			})
			if err != nil {
				return err
			}
		}

		var entries []*types.CallEntry

		for _, parsedReg := range parsed.Registrations {
			quotaID, err := database.CreateQuota(tx, parsedReg.Course.Quota)
			if err != nil {
				return err
			}

			courseID, err := database.CreateCourse(tx, parsedReg.Course, quotaID)
			if err != nil {
				return err
			}

			candidateID, err := database.CreateCandidate(tx, parsedReg.Registration.Candidate)
			if err != nil {
				return err
			}

			registrationID, err := database.CreateRegistration(tx, &database.CreateRegistrationArgs{
				Registration: parsedReg.Registration,
				CandidateID:  candidateID,
				CourseID:     courseID,
				SelectionID:  selectionID,
			})
			if err != nil {
				return err
			}

			if cmd.Kind == types.SelectionKindApproved {
				entries = append(entries, &types.CallEntry{
					CallID:         callID,
					RegistrationID: registrationID,
					Kind:           types.CallEntryKindInitial,
					Semester:       parsedReg.Course.SemesterForRanking(parsedReg.Registration.Ranking),
					Outcome:        types.CallEntryOutcomePending,
				})
			}
		}

		return database.CreateCallEntries(tx, entries)
	})
	if err != nil {
		return fmt.Errorf("create selection: %w", err)
	}

	return nil
}
