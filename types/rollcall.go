package types

import (
	"database/sql"
	"errors"
	"fmt"
)

var (
	ErrActiveRollcalls         = errors.New("there's still open rollcalls")
	ErrPendingApplications     = errors.New("there's pending applications in this calling")
	ErrUnknownStatus           = errors.New("unknown status")
	ErrNewerRollcallsExists    = errors.New("there's rollcalls newer than this one")
	ErrCantDeleteFirstRollcall = errors.New("can't delete first rollcall")
)

type Rollcall struct {
	ID     *int64 `csv:"RollcallID"`
	Status string `csv:"RollcallStatus"`
	Number int64
}

type RollcallRepository interface {
	CreateRollcall() (*Rollcall, error)
	FindRollcallByID(id *int64) (*Rollcall, error)
	FindRollcallByNumber(number int64) (*Rollcall, error)
	FindRollcalls(filter RollcallsFilter) ([]*Rollcall, error)
	UpdateRollcall(id int64, update RollcallUpdate) (*Rollcall, error)
	DeleteRollcall(id int64) error
}

type RollcallsFilter struct {
	Status *string
}

type RollcallUpdate struct {
	Status *string
}

func CanCreateRollcall(repo RollcallRepository) error {
	status := "CALLING"
	filter := RollcallsFilter{Status: &status}

	_, err := repo.FindRollcalls(filter)
	if err == nil {
		return ErrActiveRollcalls
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("find rollcalls: %w", err)
	}

	return nil
}

func CanDeleteRollcall(repo RollcallRepository, id int64) error {
	var filter RollcallsFilter

	rollcalls, err := repo.FindRollcalls(filter)
	if err != nil {
		return fmt.Errorf("find rollcalls: %w", err)
	}

	rollcall, err := repo.FindRollcallByID(&id)
	if err != nil {
		return fmt.Errorf("find rollcall: %w", err)
	}

	if int64(len(rollcalls)) > rollcall.Number {
		return ErrNewerRollcallsExists
	}

	if rollcall.Number == 1 {
		return ErrCantDeleteFirstRollcall
	}

	return nil
}

func CanOpenRollcall(repo RollcallRepository, id int64) error {
	var filter RollcallsFilter

	rollcalls, err := repo.FindRollcalls(filter)
	if err != nil {
		return fmt.Errorf("find rollcalls: %w", err)
	}

	rollcall, err := repo.FindRollcallByID(&id)
	if err != nil {
		return fmt.Errorf("find rollcall: %w", err)
	}

	if int64(len(rollcalls)) > rollcall.Number {
		return ErrNewerRollcallsExists
	}

	return nil
}

func CanUpdateRollcall(repo RollcallRepository, applicationRepo ApplicationRepository, id int64, status string) error {
	switch status {
	case "CALLING":
		return CanOpenRollcall(repo, id)
	case "DONE":
		return canCloseRollcall(applicationRepo, id)
	default:
		return ErrUnknownStatus
	}
}

func canCloseRollcall(repo ApplicationRepository, id int64) error {
	status := []string{"WAITING", "APPROVED"}

	for i := 0; i < len(status); i++ {
		var filter ApplicationsFilter
		filter.RollcallID = &id
		filter.Status = &status[i]

		applications, err := repo.FindApplications(filter)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("find applications: %w", err)
		}

		if len(applications) > 0 {
			return ErrPendingApplications
		}
	}

	return nil
}
