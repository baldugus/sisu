// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package repository

import (
	"database/sql"
)

type Applicant struct {
	ID           int64
	Cpf          string
	Name         string
	SocialName   string
	Birthdate    string
	Sex          string
	MotherName   string
	AddressLine  string
	AddressLine2 string
	HouseNumber  string
	Neighborhood string
	Municipality string
	State        string
	Cep          string
	Email        string
	Phone1       string
	Phone2       string
}

type Application struct {
	ID                   int64
	EnrollmentID         string
	ClassID              int64
	Option               int64
	LanguagesScore       float64
	HumanitiesScore      float64
	NaturalSciencesScore float64
	MathematicsScore     float64
	EssayScore           float64
	CompositeScore       float64
	Ranking              int64
	Status               sql.NullString
	SelectionID          int64
	RollcallID           sql.NullInt64
	ApplicantID          int64
}

type Class struct {
	ID           int64
	PeriodID     int64
	QuotaID      int64
	Seats        int64
	MinimumScore float64
}

type Period struct {
	ID   int64
	Name string
}

type Quota struct {
	ID   int64
	Name string
}

type Rollcall struct {
	ID     int64
	Number int64
	Status sql.NullString
}

type Selection struct {
	ID          int64
	Kind        int64
	Name        string
	Date        string
	Institution string
	Course      string
}
