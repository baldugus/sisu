package repository

/*
// nolint: godox,nolintlint
// TODO: make some helper functions to avoid the repetition.
package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"strings"

	"changeme/types"

	"github.com/jmoiron/sqlx"
)

type (
	SelectionArgument []types.SelectionKind
	StatusArgument    []string
	RollcallArgument  []int64
	ClassArgument     []int64
)

// TODO: move restriction logic to SISU, not here.
// TODO: also move insertion dependency logic (tx) to SISU.
// TODO: anything that doesn't relate to database logic (such as checks and restrictions before actions are made) should be in SISU.
// FIXME: these errors make no sense, make them reasonable and move them closer to where they're used.
var (
	ErrImportExists              = errors.New("import already exists")
	ErrInterestedSelectionExists = errors.New("interested selection exists")
	ErrRollCallExists            = errors.New("rollcall already exists")
	ErrImportMissing             = errors.New("first import does not exist")
	ErrDoneRollCallExists        = errors.New("finished rollcall exists")
	ErrRollCallClose             = errors.New("rollcall has pending applications")
	ErrRollCallMissing           = errors.New("no rollcall with the provided id")
	ErrUpdatedApplicationExists  = errors.New("updated application exists")
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) SaveSelection(selection *types.Selection) error {
	tx, err := r.db.Beginx()
	defer func() { _ = tx.Rollback() }()

	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	selectionID, err := saveSelection(tx, selection)
	if err != nil {
		return fmt.Errorf("save selection: %w", err)
	}

	var rollCallID *int64

	if selection.Kind == types.ApprovedSelection {
		id, err := saveRollCall(tx)
		if err != nil {
			return fmt.Errorf("save roll call: %w", err)
		}

		rollCallID = &id
	}

	// TODO: test null rollcall insertion
	for _, application := range selection.Applications {
		_, err := saveApplication(tx, application, selectionID, rollCallID)
		if err != nil {
			return fmt.Errorf("save application: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

func (r *Repository) DeleteRollCall(number int64) error {
	numberArgument := NumberArgument{number}
}

// FIXME: why did this code die and why is it this long and cyclomatic?
func (r *Repository) DeleteSelection(kind types.SelectionKind) error { //nolint: funlen, cyclop
	if kind == types.ApprovedSelection {
		if _, err := r.FindSelection(types.InterestedSelection); err == nil {
			return ErrInterestedSelectionExists
		}
	}

	selection := SelectionArgument{kind}
	status := StatusArgument{"ENROLLED", "ABSENT"}

	_, err := r.SearchApplications(nil, selection, status, nil, nil)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("search application: %w", err)
	} else if err == nil {
		return ErrUpdatedApplicationExists
	}

	status = StatusArgument{"DONE"}

	if _, err := r.SearchRollCalls(nil, selection, status); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("search rollcalls: %w", err)
	} else if err == nil {
		return ErrDoneRollCallExists
	}

	tx, err := r.db.Beginx()
	defer func() { _ = tx.Rollback() }()

	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	if err := deleteRollcallBySelection(tx, kind); err != nil {
		return fmt.Errorf("delete rollcall by selection id: %w", err)
	}

	if err := deleteApplicantBySelection(tx, kind); err != nil {
		return fmt.Errorf("delete applicant by selection id: %w", err)
	}

	if err := deleteApplicationBySelection(tx, kind); err != nil {
		return fmt.Errorf("delete application by selection id: %w", err)
	}

	query := "DELETE FROM selections WHERE id = ?"

	if _, err := tx.Exec(query, kind); err != nil {
		return fmt.Errorf("execute query: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

func (r *Repository) FindSelection(kind types.SelectionKind) (*types.Selection, error) {
	var selection types.Selection

	query := "SELECT id, name, date FROM selections WHERE id = $1"

	if err := r.db.Get(&selection, query, kind); err != nil {
		return nil, fmt.Errorf("get selections query: %w", err)
	}

	query = `
  SELECT
    applications.id, enrollment_id, address_line, address_line2, status,
    birthdate, sex, cep, composite_score, cpf, email, essay_score, house_number,
    humanities_score, languages_score, mathematics_score, mother_name,
    municipality, state, applicants.name, natural_sciences_score, neighborhood,
    option, phone1, phone2, ranking, quotas.name as quota, applicants.social_name,
    periods.name as period, seats, minimum_score
  FROM applications
  JOIN applicants ON applications.id = applicants.application_id
  JOIN classes ON applications.class_id = classes.id
  JOIN periods ON periods.id = classes.period_id
  JOIN quotas ON quotas.id = classes.quota_id
  WHERE selection_id = $1
  `

	err := r.db.Select(&selection.Applications, query, kind)
	if err != nil {
		return nil, fmt.Errorf("select from applications: %w", err)
	}

	return &selection, nil
}

// FIXME: for some reason its too complex.
func (r *Repository) SearchApplications( //nolint: funlen, cyclop
	ids ApplicationArgument,
	selections SelectionArgument,
	status StatusArgument,
	rollcalls RollcallArgument,
	classes ClassArgument,
) (
	[]*types.Application,
	error,
) {
	var builder strings.Builder

	builder.WriteString(`
	SELECT
      applications.id, enrollment_id, address_line, address_line2, status,
      birthdate, sex, cep, composite_score, cpf, email, essay_score, house_number,
      humanities_score, languages_score, mathematics_score, mother_name,
      municipality, state, applicants.name, natural_sciences_score, neighborhood,
      option, phone1, phone2, ranking, quotas.name as quota, applicants.social_name,
      periods.name as period, seats, minimum_score
    FROM applications
    JOIN applicants ON applications.id = applicants.application_id
    JOIN classes ON applications.class_id = classes.id
    JOIN periods ON periods.id = classes.period_id
    JOIN quotas ON quotas.id = classes.quota_id
    WHERE 1 = 1
	`)

	var args []any

	if len(ids) > 0 {
		builder.WriteString(" AND applications.id IN (?)")

		args = append(args, ids)
	}

	if len(classes) > 0 {
		builder.WriteString(" AND class_id IN (?)")

		args = append(args, classes)
	}

	if len(selections) > 0 {
		builder.WriteString(" AND selection_id IN (?)")

		args = append(args, selections)
	}

	if len(status) > 0 {
		builder.WriteString(" AND applications.status IN (?)")

		args = append(args, status)
	}

	if len(rollcalls) > 0 {
		builder.WriteString(" AND rollcall_id IN (?)")

		args = append(args, rollcalls)
	}

	builder.WriteString(" GROUP BY applications.id")

	query, args, err := sqlx.In(builder.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("sqlx in: %w", err)
	}

	query = r.db.Rebind(query)

	applications := []*types.Application{}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("select applications with status=%v: %w", status, err)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("select applications with status=%v: %w", status, err)
	}

	if err := sqlx.StructScan(rows, &applications); err != nil {
		return nil, fmt.Errorf("struct scan: %w", err)
	}

	if len(applications) == 0 {
		return nil, sql.ErrNoRows
	}

	return applications, nil
}

func (r *Repository) SearchClasses() ([]*types.Class, error) {
	query := `
	  SELECT classes.id,
	         periods.name as period,
	         quotas.name as quota,
	         seats,
	         minimum_score
	  FROM classes
	  JOIN periods ON classes.period_id = periods.id
	  JOIN quotas ON classes.quota_id = quotas.id
	`

	var classes []*types.Class

	if err := r.db.Select(&classes, query); err != nil {
		return nil, fmt.Errorf("db select: %w", err)
	}

	return classes, nil
}

func (r *Repository) SearchRollCalls( //nolint: funlen
	numbers NumberArgument,
	selections SelectionArgument,
	status StatusArgument,
) (
	[]*types.RollCall,
	error,
) {
	var builder strings.Builder

	builder.WriteString(`
    SELECT rollcalls.id as number, rollcalls.status
    FROM rollcalls
	LEFT JOIN applications ON rollcalls.id = applications.rollcall_id
    WHERE 1 = 1
    `)

	var args []any

	if len(numbers) > 0 {
		builder.WriteString(" AND number IN (?)")

		args = append(args, numbers)
	}

	if len(selections) > 0 {
		builder.WriteString(" AND applications.selection_id IN (?)")

		args = append(args, selections)
	}

	if len(status) > 0 {
		builder.WriteString(" AND rollcalls.status IN (?)")

		args = append(args, status)
	}

	builder.WriteString(" GROUP BY rollcalls.id")

	query, args, err := sqlx.In(builder.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("sqlx in: %w", err)
	}

	query = r.db.Rebind(query)

	var rollcalls []*types.RollCall

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("select rollcalls: %w", err)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows: %w", err)
	}

	if err := sqlx.StructScan(rows, &rollcalls); err != nil {
		return nil, fmt.Errorf("struct scan: %w", err)
	}

	if len(rollcalls) == 0 {
		return nil, sql.ErrNoRows
	}

	return rollcalls, nil
}

func (r *Repository) UpdateRollCallStatus(id int64, status string) error { //nolint: varnamelen
	query := "UPDATE rollcalls SET status = $1 WHERE id = $2"

	result, err := r.db.Exec(query, status, id)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if rows < 1 {
		return sql.ErrNoRows
	}

	return nil
}

// TODO: write test for this.
func (r *Repository) UpdateApplicationStatus(id int64, status string) error {
	query := "UPDATE applications SET status = $1 WHERE id = $2"

	result, err := r.db.Exec(query, status, id)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if rows < 1 {
		return ErrRollCallMissing
	}

	return nil
}

// TODO: this is a clone.
func updateApplication(tx *sqlx.Tx, id int64, status string, rollCallNumber int64) error {
	query := "UPDATE applications SET status = $1, rollcall_id = $2 WHERE id = $3"

	result, err := tx.Exec(query, status, rollCallNumber, id)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if rows < 1 {
		return ErrRollCallMissing
	}

	return nil
}

func saveSelection(tx *sqlx.Tx, selection *types.Selection) (int64, error) {
	var ids []int64

	query := "SELECT id FROM selections"

	err := tx.Select(&ids, query, selection.Kind)
	if err != nil {
		return 0, fmt.Errorf("select selections: %w", err)
	}

	if len(ids) == 0 && selection.Kind != types.ApprovedSelection {
		return 0, ErrImportMissing
	}

	if len(ids) >= int(selection.Kind) {
		return 0, ErrImportExists
	}

	query = "INSERT INTO selections (id, name, date) VALUES (:id, :name, :date)"

	result, err := tx.NamedExec(query, selection)
	if err != nil {
		return 0, fmt.Errorf("insert selection: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("selection last insert id: %w", err)
	}

	return id, nil
}

func saveApplication(
	tx *sqlx.Tx,
	application *types.Application,
	selectionID int64,
	rollcallID *int64,
) (int64, error) {
	classID, err := saveClass(tx, &application.Class)
	if err != nil {
		return 0, fmt.Errorf("save class: %w", err)
	}

	query := `
  INSERT INTO applications (
    enrollment_id, class_id, option, languages_score,
    humanities_score, natural_sciences_score, mathematics_score, essay_score,
    composite_score, ranking, status, selection_id, rollcall_id
  ) VALUES (
    :enrollment_id, :class_id, :option, :languages_score,
    :humanities_score, :natural_sciences_score, :mathematics_score, :essay_score,
    :composite_score, :ranking, :status, :selection_id, :rollcall_id
  )
  `

	data := struct {
		*types.Application
		ClassID     int64  `db:"class_id"`
		SelectionID int64  `db:"selection_id"`
		RollcallID  *int64 `db:"rollcall_id"`
	}{
		Application: application,
		ClassID:     classID,
		SelectionID: selectionID,
		RollcallID:  rollcallID,
	}

	result, err := tx.NamedExec(query, data)
	if err != nil {
		return 0, fmt.Errorf("insert into applications: %w", err)
	}

	id, err := result.LastInsertId() //nolint: varnamelen
	if err != nil {
		return 0, fmt.Errorf("applications last insert id: %w", err)
	}

	if _, err := saveApplicant(tx, &application.Applicant, id); err != nil {
		return 0, fmt.Errorf("save applicant %w", err)
	}

	return id, nil
}

func saveApplicant(tx *sqlx.Tx, applicant *types.Applicant, applicationID int64) (int64, error) {
	query := `
  INSERT INTO applicants (
    cpf, application_id, name, birthdate, sex, mother_name, address_line, address_line2,
    house_number, neighborhood, municipality, state, cep, email, phone1, social_name,
    phone2
  ) VALUES (
    :cpf, :application_id, :name, :birthdate, :sex, :mother_name, :address_line, :address_line2,
    :house_number, :neighborhood, :municipality, :state, :cep, :email, :phone1, :social_name,
    :phone2
  )
  `

	data := struct {
		*types.Applicant
		ApplicationID int64 `db:"application_id"`
	}{
		Applicant:     applicant,
		ApplicationID: applicationID,
	}

	result, err := tx.NamedExec(query, data)
	if err != nil {
		return 0, fmt.Errorf("insert into applicants: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("applicants last insert id: %w", err)
	}

	return id, nil
}

func saveClass(tx *sqlx.Tx, class *types.Class) (int64, error) {
	periodID, err := savePeriod(tx, class.Period)
	if err != nil {
		return 0, fmt.Errorf("save period: %w", err)
	}

	quotaID, err := saveQuota(tx, class.Quota)
	if err != nil {
		return 0, fmt.Errorf("save quota: %w", err)
	}

	var classID int64

	query := `SELECT classes.id FROM classes
						JOIN periods ON classes.period_id = periods.id
						JOIN quotas ON classes.quota_id = quotas.id
						WHERE periods.id = $1 AND quotas.id = $2
	`

	if err := tx.Get(&classID, query, periodID, quotaID); err != nil &&
		!errors.Is(err, sql.ErrNoRows) {
		return 0, fmt.Errorf("select from classes: %w", err)
	} else if err == nil {
		return classID, nil
	}

	query = "INSERT INTO classes (period_id, quota_id, seats, minimum_score) VALUES ($1, $2, $3, $4)"

	result, err := tx.Exec(query, periodID, quotaID, class.Seats, class.MinimumScore)
	if err != nil {
		return 0, fmt.Errorf("insert into classes: %w", err)
	}

	classID, err = result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("classes last insert id: %w", err)
	}

	return classID, nil
}

func savePeriod(tx *sqlx.Tx, period string) (int64, error) {
	var periodID int64

	query := "SELECT id FROM periods WHERE name = $1"
	if err := tx.Get(&periodID, query, period); err != nil &&
		!errors.Is(err, sql.ErrNoRows) {
		return 0, fmt.Errorf("select from periods: %w", err)
	} else if err == nil {
		return periodID, nil
	}

	query = "INSERT INTO periods (name) VALUES ($1)"

	result, err := tx.Exec(query, period)
	if err != nil {
		return 0, fmt.Errorf("insert into periods: %w", err)
	}

	periodID, err = result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("periods last insert id: %w", err)
	}

	return periodID, nil
}

func saveQuota(tx *sqlx.Tx, quota string) (int64, error) {
	var quotaID int64

	query := "SELECT id FROM quotas WHERE name = $1"
	if err := tx.Get(&quotaID, query, quota); err != nil &&
		!errors.Is(err, sql.ErrNoRows) {
		return 0, fmt.Errorf("select from quotas: %w", err)
	} else if err == nil {
		return quotaID, nil
	}

	query = "INSERT INTO quotas (name) VALUES ($1)"

	result, err := tx.Exec(query, quota)
	if err != nil {
		return 0, fmt.Errorf("insert into quotas: %w", err)
	}

	quotaID, err = result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("quotas last insert id: %w", err)
	}

	return quotaID, nil
}

func (r *Repository) CreateRollCall() error {
	tx, err := r.db.Beginx()
	defer func() { _ = tx.Rollback() }()

	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	id, err := saveRollCall(tx)
	if err != nil {
		return fmt.Errorf("save rollcall: %w", err)
	}

	if err := r.allocApplications(tx, id); err != nil {
		return fmt.Errorf("alloc applications: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

func (r *Repository) allocApplications(tx *sqlx.Tx, rollCallNumber int64) error {
	classes, err := r.SearchClasses()
	if err != nil {
		return fmt.Errorf("search classes: %w", err)
	}

	for _, class := range classes {
		seats := class.Seats
		statusArgument := StatusArgument{"ENROLLED", "ABSENT"}
		classArgument := ClassArgument{class.ID}

		applications, err := r.SearchApplications(nil, nil, statusArgument, nil, classArgument)
		if err != nil {
			return fmt.Errorf("search applications: %w", err)
		}

		for _, application := range applications {
			if application.Status == "ENROLLED" {
				seats--

				continue
			}
		}

		if err := r.fillClass(tx, class, seats, rollCallNumber); err != nil {
			return fmt.Errorf("fill class: %w", err)
		}
	}

	return nil
}

type applicationsRanked []*types.Application

func (v applicationsRanked) Len() int           { return len(v) }
func (v applicationsRanked) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }
func (v applicationsRanked) Less(i, j int) bool { return v[i].Scores.Ranking < v[j].Scores.Ranking }

func (r *Repository) fillClass(tx *sqlx.Tx, class *types.Class, seats int, rollCallNumber int64) error {
	classArgument := ClassArgument{class.ID}
	statusArgument := StatusArgument{"WAITING"}

	applications, err := r.SearchApplications(nil, nil, statusArgument, nil, classArgument)
	if err != nil {
		return fmt.Errorf("search applications: %w", err)
	}

	sort.Sort(applicationsRanked(applications))

	applications = applications[:seats]

	for _, application := range applications {
		if err := updateApplication(tx, application.ID, "APPROVED", rollCallNumber); err != nil {
			return fmt.Errorf("update application: %w", err)
		}
	}

	return nil
}

func saveRollCall(tx *sqlx.Tx) (int64, error) {
	query := "INSERT INTO rollcalls (status) VALUES ('CALLING')"

	result, err := tx.Exec(query)
	if err != nil {
		return 0, fmt.Errorf("insert rollcall: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("rollcalls last insert id: %w", err)
	}

	return id, nil
}

type (
	NumberArgument      []int64
	ApplicationArgument []int64
)

func deleteRollcallBySelection(tx *sqlx.Tx, kind types.SelectionKind) error {
	query := `
      DELETE FROM rollcalls
      WHERE id IN (
        SELECT rollcall_id
        FROM applications
        WHERE selection_id = ?
      )
	`

	if _, err := tx.Exec(query, kind); err != nil {
		return fmt.Errorf("execute query: %w", err)
	}

	return nil
}

func deleteApplicationBySelection(tx *sqlx.Tx, kind types.SelectionKind) error {
	query := "DELETE FROM applications WHERE selection_id = ?"

	if _, err := tx.Exec(query, kind); err != nil {
		return fmt.Errorf("execute query: %w", err)
	}

	return nil
}

func deleteApplicantBySelection(tx *sqlx.Tx, kind types.SelectionKind) error {
	query := `
	  DELETE FROM applicants
	  WHERE application_id IN (
	    SELECT id
	    FROM applications
	    WHERE selection_id = ?
	  )
    `

	if _, err := tx.Exec(query, kind); err != nil {
		return fmt.Errorf("execute query: %w", err)
	}

	return nil
}
*/
