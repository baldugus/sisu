package commands

import (
	"database/sql"
	"errors"
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
		var sem1ID, sem2ID int32
		if cmd.Kind == types.SelectionKindApproved {
			sem1, err := database.FetchSemesterByYearAndNumber(tx, cmd.Year, 1)
			if err != nil {
				if errors.Is(err, qrm.ErrNoRows) || errors.Is(err, sql.ErrNoRows) {
					sem1ID, err = database.CreateSemester(tx, cmd.Year, 1)
					if err != nil {
						return err
					}
				} else {
					return err
				}
			} else {
				sem1ID = sem1.ID
			}

			sem2, err := database.FetchSemesterByYearAndNumber(tx, cmd.Year, 2)
			if err != nil {
				if errors.Is(err, qrm.ErrNoRows) || errors.Is(err, sql.ErrNoRows) {
					sem2ID, err = database.CreateSemester(tx, cmd.Year, 2)
					if err != nil {
						return err
					}
				} else {
					return err
				}
			} else {
				sem2ID = sem2.ID
			}
		}

		selectionID, err := database.CreateSelection(tx, parsed.Selection)
		if err != nil {
			return err
		}

		var callID *int32
		if cmd.Kind == types.SelectionKindApproved {
			id, err := database.CreateCall(tx, &types.Call{
				Number:     1,
				Status:     types.CallStatusCalling,
				SemesterID: sem1ID,
			})
			if err != nil {
				return err
			}
			callID = &id
		}

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

			var mappedSemesterID *int32
			if parsedReg.Registration.SemesterID != nil {
				if *parsedReg.Registration.SemesterID == 1 {
					mappedSemesterID = &sem1ID
				} else if *parsedReg.Registration.SemesterID == 2 {
					mappedSemesterID = &sem2ID
				}
			}

			err = database.CreateRegistration(tx, &database.CreateRegistrationArgs{
				Registration: parsedReg.Registration,
				CandidateID:  candidateID,
				CourseID:     courseID,
				SelectionID:  selectionID,
				CallID:       callID,
				SemesterID:   mappedSemesterID,
			})
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("create selection: %w", err)
	}

	return nil
}
