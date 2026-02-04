package repository

import (
	"fmt"

	"changeme/types"
)

type ApplicantRepository struct {
	db *DB
}

func NewApplicantRepository(db *DB) *ApplicantRepository {
	return &ApplicantRepository{
		db: db,
	}
}

func (a *ApplicantRepository) CreateApplicant(applicant *types.Applicant) error {
	query := `
		INSERT INTO applicants (
			cpf, name, birthdate, sex, mother_name, address_line, address_line2,
    		house_number, neighborhood, municipality, state, cep, email, phone1, social_name,
   			phone2
		) VALUES (
			:cpf, :name, :birthdate, :sex, :mother_name, :address_line, :address_line2,
			:house_number, :neighborhood, :municipality, :state, :cep, :email, :phone1, :social_name,
			:phone2
		)
	`

	result, err := a.db.NamedExec(query, applicant)
	if err != nil {
		return fmt.Errorf("db named exec: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("last insert id: %w", err)
	}

	applicant.ID = id

	return nil
}

func (a *ApplicantRepository) FindApplicantByID(id int64) (*types.Applicant, error) {
	var applicant types.Applicant

	query := "SELECT * FROM applicants WHERE id = ?"
	if err := a.db.Get(&applicant, query, id); err != nil {
		return nil, fmt.Errorf("db get: %w", err)
	}

	return &applicant, nil
}

func (a *ApplicantRepository) FindApplicantByApplicationID(id int64) (*types.Applicant, error) {
	var applicant types.Applicant

	query := "SELECT * FROM applicants WHERE application_id = ?"
	if err := a.db.Get(&applicant, query, id); err != nil {
		return nil, fmt.Errorf("db get: %w", err)
	}

	return &applicant, nil
}

func (a *ApplicantRepository) UpdateApplicant(id int64, update types.ApplicantUpdate) (*types.Applicant, error) {
	applicant, err := a.FindApplicantByID(id)
	if err != nil {
		return nil, fmt.Errorf("find rollcall by id: %w", err)
	}

	if update.Email != nil {
		applicant.Email = *update.Email
	}

	query := "UPDATE applicant SET email = :email WHERE id = :id"

	if _, err := a.db.NamedExec(query, applicant); err != nil {
		return nil, fmt.Errorf("db select: %w", err)
	}

	return applicant, nil
}

func (a *ApplicantRepository) DeleteApplicant(id int64) error {
	query := "DELETE FROM applicants WHERE id = ?"

	_, err := a.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("db exec: %w", err)
	}

	return nil
}
