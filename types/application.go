package types

import (
	"database/sql"
	"errors"
	"fmt"
)

type Application struct {
	ID                   int64
	Status               string
	EnrollmentID         string `db:"enrollment_id"`
	Option               int
	LanguagesScore       float64 `db:"languages_score"`
	HumanitiesScore      float64 `db:"humanities_score"`
	NaturalSciencesScore float64 `db:"natural_sciences_score"`
	MathematicsScore     float64 `db:"mathematics_score"`
	EssayScore           float64 `db:"essay_score"`
	CompositeScore       float64 `db:"composite_score"`
	Ranking              int
	Applicant            `db:"applicant"`
	Rollcall             `db:"rollcall" json:"-"`
	Class                `db:"class"`
}

var (
	ErrApplicationNotEditable = errors.New("this application is not editable")
	ErrRollcallEnded          = errors.New("this application's rollcall is over")
)

type ApplicationRepository interface {
	CreateApplication(application *Application, selectionID int64) error
	FindApplicationByID(id int64) (*Application, error)
	FindApplications(filter ApplicationsFilter) ([]*Application, error)
	UpdateApplication(id int64, update ApplicationUpdate) (*Application, error)
	DeleteApplication(id int64) error
}

type ApplicationsFilter struct {
	SelectionID *int64
	Status      *string
	RollcallID  *int64
	ClassID     *int64
}

type ApplicationUpdate struct {
	Status     *string
	RollcallID *int64
}

func CanUpdateApplication(repo ApplicationRepository, rollcallRepo RollcallRepository, id int64, update ApplicationUpdate) error {
	application, err := repo.FindApplicationByID(id)
	if err != nil {
		return fmt.Errorf("find application by id: %w", err)
	}

	if application.Status == "WAITING" && *update.Status != "APPROVED" {
		return ErrApplicationNotEditable
	}

	rollcall, err := rollcallRepo.FindRollcallByID(application.Rollcall.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}

		return fmt.Errorf("find rollcall by id: %w", err)
	}

	if rollcall.Status == "DONE" {
		return ErrRollcallEnded
	}

	return nil
}
