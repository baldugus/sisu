package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"changeme/types"
)

type ApplicationRepository struct {
	db *DB
}

func NewApplicationRepository(db *DB) *ApplicationRepository {
	return &ApplicationRepository{
		db: db,
	}
}

func (a *ApplicationRepository) FindApplicationByID(id int64) (*types.Application, error) {
	var application types.Application

	query := `
		SELECT 
		applications.id, enrollment_id, option, languages_score,
		humanities_score, natural_sciences_score, mathematics_score, essay_score,
		composite_score, ranking, applications.status,
		
		classes.id as "class.id", seats as "class.seats", minimum_score as "class.minimum_score",
		quotas.id as "class.quota.id", quotas.name as "class.quota.name",
		periods.id as "class.period.id", periods.name as "class.period.name",

		rollcalls.id as "rollcall.id", IFNULL(rollcalls.number, 0) as "rollcall.number", IFNULL(rollcalls.status, "") as "rollcall.status",
		
		applicants.id as "applicant.id", cpf as "applicant.cpf", applicants.name as "applicant.name", social_name as "applicant.social_name", birthdate as "applicant.birthdate", sex as "applicant.sex", mother_name as "applicant.mother_name", address_line as "applicant.address_line", address_line2 as "applicant.address_line2", house_number as "applicant.house_number", neighborhood as "applicant.neighborhood", municipality as "applicant.municipality", state as "applicant.state", cep as "applicant.cep", email as "applicant.email", phone1 as "applicant.phone1", phone2 as "applicant.phone2"
		FROM applications
		JOIN classes ON applications.class_id = classes.id
		JOIN periods ON classes.period_id = periods.id
		JOIN quotas ON classes.quota_id = quotas.id
		LEFT JOIN rollcalls ON applications.rollcall_id = rollcalls.id
		JOIN applicants ON applications.applicant_id = applicants.id
		WHERE applications.id = ?
		`
	if err := a.db.Get(&application, query, id); err != nil {
		return nil, fmt.Errorf("db get: %w", err)
	}

	return &application, nil
}

func (a *ApplicationRepository) FindApplications(filter types.ApplicationsFilter) ([]*types.Application, error) {
	var applications []*types.Application

	var builder strings.Builder

	builder.WriteString(`
		SELECT 
		applications.id, enrollment_id, option, languages_score,
		humanities_score, natural_sciences_score, mathematics_score, essay_score,
		composite_score, ranking, applications.status,
		
		classes.id as "class.id", seats as "class.seats", minimum_score as "class.minimum_score",
		quotas.id as "class.quota.id", quotas.name as "class.quota.name",
		periods.id as "class.period.id", periods.name as "class.period.name",

		rollcalls.id as "rollcall.id", IFNULL(rollcalls.number, 0) as "rollcall.number", IFNULL(rollcalls.status, "") as "rollcall.status",
		
		applicants.id as "applicant.id", cpf as "applicant.cpf", applicants.name as "applicant.name", social_name as "applicant.social_name", birthdate as "applicant.birthdate", sex as "applicant.sex", mother_name as "applicant.mother_name", address_line as "applicant.address_line", address_line2 as "applicant.address_line2", house_number as "applicant.house_number", neighborhood as "applicant.neighborhood", municipality as "applicant.municipality", state as "applicant.state", cep as "applicant.cep", email as "applicant.email", phone1 as "applicant.phone1", phone2 as "applicant.phone2"
		FROM applications
		JOIN classes ON applications.class_id = classes.id
		JOIN periods ON classes.period_id = periods.id
		JOIN quotas ON classes.quota_id = quotas.id
		LEFT JOIN rollcalls ON applications.rollcall_id = rollcalls.id
		JOIN applicants ON applications.applicant_id = applicants.id
		WHERE 1 = 1
    `)

	var args []any

	if filter.ClassID != nil {
		builder.WriteString(" AND class_id = ?")

		args = append(args, *filter.ClassID)
	}

	if filter.RollcallID != nil {
		builder.WriteString(" AND rollcall_id = ?")

		args = append(args, *filter.RollcallID)
	}

	if filter.SelectionID != nil {
		builder.WriteString(" AND selection_id = ?")

		args = append(args, *filter.SelectionID)
	}

	if filter.Status != nil {
		builder.WriteString(" AND applications.status = ?")

		args = append(args, *filter.Status)
	}

	if err := a.db.Select(&applications, builder.String(), args...); err != nil {
		return nil, fmt.Errorf("db select: %w", err)
	}

	if len(applications) == 0 {
		return nil, sql.ErrNoRows
	}

	return applications, nil
}

func (a *ApplicationRepository) CreateApplication(application *types.Application, selectionID int64) error {
	query := `
	INSERT INTO applications (
		enrollment_id, class_id, option, languages_score,
		humanities_score, natural_sciences_score, mathematics_score, essay_score,
		composite_score, ranking, status, selection_id, rollcall_id, applicant_id
  	) VALUES (
		:enrollment_id, :class.id, :option, :languages_score,
		:humanities_score, :natural_sciences_score, :mathematics_score, :essay_score,
		:composite_score, :ranking, :status, :selection_id, :rollcall.id, :applicant.id
  	)
	`

	data := struct {
		types.Application
		SelectionID int64 `db:"selection_id"`
	}{
		Application: *application,
		SelectionID: selectionID,
	}

	result, err := a.db.NamedExec(query, data)
	if err != nil {
		return fmt.Errorf("db named exec: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("last insert id: %w", err)
	}

	application.ID = id

	return nil
}

func (a *ApplicationRepository) UpdateApplication(id int64, update types.ApplicationUpdate) (*types.Application, error) {
	rollcallRepo := NewRollcallRepository(a.db)
	if err := types.CanUpdateApplication(a, rollcallRepo, id, update); err != nil {
		return nil, fmt.Errorf("can update application: %w", err)
	}

	application, err := a.FindApplicationByID(id)
	if err != nil {
		return nil, fmt.Errorf("find application by id: %w", err)
	}

	if update.Status != nil {
		application.Status = *update.Status
	}

	if update.RollcallID != nil {
		var rollcall types.Rollcall
		rollcall.ID = update.RollcallID
		if *update.RollcallID == 0 {
			rollcall.ID = nil
		}
		application.Rollcall = rollcall
	}

	query := "UPDATE applications SET status = :status, rollcall_id = :rollcall.id WHERE id = :id"

	if _, err := a.db.NamedExec(query, application); err != nil {
		return nil, fmt.Errorf("db select: %w", err)
	}

	return application, nil
}

func (a *ApplicationRepository) DeleteApplication(id int64) error {
	query := "DELETE FROM applications WHERE id = ?"

	_, err := a.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("db exec: %w", err)
	}

	return nil
}
