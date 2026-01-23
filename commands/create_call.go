package commands

import (
	"github.com/baldugus/sisu/database"
	"github.com/baldugus/sisu/types"
	"github.com/go-jet/jet/v2/qrm"
)

type CreateCallCommand struct{}

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
			Number: lastCallNumber + 1,
			Status: types.CallStatusCalling,
		}

		callID, err := database.CreateCall(tx, call)
		if err != nil {
			return err
		}

		totalPromoted := 0
		totalAvailableSeats := int32(0)

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

			registrationIDs, err := database.FetchWaitlistedRegistrationsByCourse(tx, course.ID, availableSeats)
			if err != nil {
				return err
			}

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
