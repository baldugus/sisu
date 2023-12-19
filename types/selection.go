package types

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type SelectionKind int

const (
	ApprovedSelection SelectionKind = iota + 1
	InterestedSelection
)

var (
	ErrUnkownKind                 = errors.New("unknown selection kind")
	ErrSelectionAlreadyExists     = errors.New("another selection of this type already exists")
	ErrApprovedSelectionMissing   = errors.New("approved selection missing")
	ErrEditedApplication          = errors.New("applications were edited")
	ErrFinishedRollcall           = errors.New("rollcalls were finished")
	ErrInterestedSelectionPresent = errors.New("interested selection present")
)

type Selection struct {
	ID           int64
	Name         string
	Kind         SelectionKind
	Date         string
	Institution  string
	Course       string
	Applications []*Application
}

type SelectionRepository interface {
	FindSelectionByID(id int64) (*Selection, error)
	FindSelectionByKind(kind SelectionKind) (*Selection, error)
	FindSelections() ([]*Selection, error)
	CreateSelection(selection *Selection) error
	DeleteSelection(id int64) error
}

func CanCreateSelection(repo SelectionRepository, kind SelectionKind) error {
	switch kind {
	case ApprovedSelection:
		return canCreateApprovedKind(repo)
	case InterestedSelection:
		return canCreateInterestedKind(repo)
	default:
		return ErrUnkownKind
	}
}

func canCreateApprovedKind(repo SelectionRepository) error {
	_, err := repo.FindSelections()
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("find selections: %w", err)
	} else if err == nil {
		return ErrSelectionAlreadyExists
	}

	return nil
}

func canCreateInterestedKind(repo SelectionRepository) error {
	selections, err := repo.FindSelections()
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("find selections: %w", err)
	} else if errors.Is(err, sql.ErrNoRows) {
		return ErrApprovedSelectionMissing
	}

	if len(selections) > 1 {
		return ErrSelectionAlreadyExists
	}

	return nil
}

// TODO: eliminate single filters, replace by functions.
func CanDeleteSelection(repo SelectionRepository, rollcallRepo RollcallRepository, applicationRepo ApplicationRepository, id int64) error {
	status := []string{"ENROLLED", "ABSENT"}
	for i := 0; i < len(status); i++ {
		filter := ApplicationsFilter{Status: &status[i]}

		_, err := applicationRepo.FindApplications(filter)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("find applications: %w", err)
		} else if err == nil {
			return ErrEditedApplication
		}

	}

	rollcallStatus := "DONE"
	filter := RollcallsFilter{Status: &rollcallStatus}

	_, err := rollcallRepo.FindRollcalls(filter)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("find rollcalls: %w", err)
	} else if err == nil {
		return ErrRollcallEnded
	}

	selection, err := repo.FindSelectionByID(id)
	if err != nil {
		return fmt.Errorf("find selection: %w", err)
	}

	if selection.Kind == ApprovedSelection {
		return canDeleteApprovedSelection(repo)
	}

	return nil
}

func canDeleteApprovedSelection(repo SelectionRepository) error {
	_, err := repo.FindSelectionByKind(InterestedSelection)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	} else if err == nil {
		return ErrInterestedSelectionPresent
	}

	return nil
}

func (s *Selection) ParseDate() (time.Time, error) {
	date, err := time.Parse(time.DateTime, s.Date)
	if err != nil {
		return time.Time{}, fmt.Errorf("time parse: %w", err)
	}

	return date, nil
}
