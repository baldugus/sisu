package commands

import (
	"database/sql"
	"errors"

	"github.com/baldugus/sisu/database"
	"github.com/baldugus/sisu/types"
)

type UpdateRegistrationStatusCommand struct {
	RegistrationID int32
	NewStatus      types.RegistrationStatus
}

func (cmd *UpdateRegistrationStatusCommand) Execute(db *database.Database) error {
	detail, err := db.FetchRegistrationByID(cmd.RegistrationID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrRegistrationNotFound{}
		}
		return err
	}

	currentStatus := detail.Registration.Status

	if currentStatus == cmd.NewStatus {
		return nil
	}

	if err := cmd.validateTransition(detail); err != nil {
		return err
	}

	return db.UpdateRegistrationStatus(cmd.RegistrationID, cmd.NewStatus)
}

func (cmd *UpdateRegistrationStatusCommand) validateTransition(
	detail *types.RegistrationDetail,
) error {
	currentStatus := detail.Registration.Status
	call := detail.Call

	switch cmd.NewStatus {
	case types.RegistrationStatusApproved:
		return cmd.validateToApproved(currentStatus, call)

	case types.RegistrationStatusEnrolled, types.RegistrationStatusAbsent, types.RegistrationStatusDeclinedPromotion:
		return cmd.validateToEnrolledOrAbsentOrDeclined(currentStatus, call)

	default:
		return ErrInvalidStatusTransition{}
	}
}

func (cmd *UpdateRegistrationStatusCommand) validateToApproved(
	currentStatus types.RegistrationStatus,
	call *types.Call,
) error {
	if currentStatus != types.RegistrationStatusEnrolled && currentStatus != types.RegistrationStatusAbsent && currentStatus != types.RegistrationStatusDeclinedPromotion {
		return ErrInvalidStatusTransition{}
	}

	return requireOpenCall(call)
}

func (cmd *UpdateRegistrationStatusCommand) validateToEnrolledOrAbsentOrDeclined(
	currentStatus types.RegistrationStatus,
	call *types.Call,
) error {
	if currentStatus != types.RegistrationStatusApproved {
		return ErrInvalidStatusTransition{}
	}

	return requireOpenCall(call)
}

func requireOpenCall(call *types.Call) error {
	if call == nil {
		return ErrRegistrationNotInCall{}
	}

	if call.Status != types.CallStatusCalling {
		return ErrCallNotOpen{}
	}

	return nil
}
