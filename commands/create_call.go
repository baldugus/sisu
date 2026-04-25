package commands

import (
	"github.com/baldugus/sisu/database"
	"github.com/baldugus/sisu/types"
	"github.com/go-jet/jet/v2/qrm"
)

type CreateCallCommand struct {
	SemesterID int32
}

func (cmd *CreateCallCommand) Execute(db *database.Database) error {
	hasOpenCall, err := db.HasOpenCall()
	if err != nil {
		return err
	}
	if hasOpenCall {
		return ErrOpenCallExists{}
	}

	lastCallNumber, err := db.GetLastCallNumber()
	if err != nil {
		return err
	}

	courses, err := db.FetchCourses()
	if err != nil {
		return err
	}

	return db.RunInTx(func(tx qrm.DB) error {
		call := &types.Call{
			Number:     lastCallNumber + 1,
			Status:     types.CallStatusCalling,
			SemesterID: cmd.SemesterID,
		}

		callID, err := database.CreateCall(tx, call)
		if err != nil {
			return err
		}

		totalPromoted := 0
		totalAvailableSeats := int32(0)

		semester, err := database.FetchSemesterByID(tx, cmd.SemesterID)
		if err != nil {
			return err
		}

		var sem2ID *int32
		if semester.Number == 1 {
			sem2, err := database.FetchSemesterByYearAndNumber(tx, semester.Year, 2)
			if err == nil {
				sem2ID = &sem2.ID
			}
		}

		for _, course := range courses {
			occupiedSeats, err := database.CountCourseOccupiedSeats(tx, course.ID)
			if err != nil {
				return err
			}

			availableSeats := course.Seats - occupiedSeats
			if availableSeats <= 0 {
				continue
			}

			totalAvailableSeats += availableSeats

			var promotedIDs []int32
			if sem2ID != nil {
				promotedIDs, err = database.FetchPriorityRegistrationsByCourse(tx, course.ID, *sem2ID, availableSeats)
				if err != nil {
					return err
				}
			}

			remainingSeats := availableSeats - int32(len(promotedIDs))
			var waitlistedIDs []int32
			if remainingSeats > 0 {
				waitlistedIDs, err = database.FetchWaitlistedRegistrationsByCourse(tx, course.ID, remainingSeats)
				if err != nil {
					return err
				}
			}

			registrationIDs := append(promotedIDs, waitlistedIDs...)

			for _, regID := range registrationIDs {
				if err := database.AssignRegistrationToCall(tx, regID, callID); err != nil {
					return err
				}
				totalPromoted++
			}
		}

		if totalPromoted == 0 {
			if totalAvailableSeats == 0 {
				return ErrAllCoursesFull{}
			}
			return ErrNoWaitlistedRegistrations{}
		}

		return nil
	})
}
