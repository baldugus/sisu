package database

import (
	"github.com/baldugus/sisu/database/.gen/model"
	. "github.com/baldugus/sisu/database/.gen/table"
	"github.com/baldugus/sisu/database/.gen/view"
	"github.com/baldugus/sisu/types"
	"github.com/go-jet/jet/v2/qrm"
	. "github.com/go-jet/jet/v2/sqlite"
)

type CreateRegistrationArgs struct {
	Registration *types.Registration
	CandidateID  int32
	CourseID     int32
	SelectionID  int32
}

func CreateRegistration(db qrm.DB, args *CreateRegistrationArgs) (int32, error) {
	registrationModel := toRegistrationModel(args.Registration)
	registrationModel.CandidateID = args.CandidateID
	registrationModel.CourseID = args.CourseID
	registrationModel.SelectionID = args.SelectionID

	stmt := Registrations.INSERT(Registrations.MutableColumns).
		MODEL(registrationModel).
		RETURNING(Registrations.ID)

	result, err := insertOne[model.Registrations](db, stmt)
	if err != nil {
		return 0, err
	}

	return result.ID, nil
}

// registrationsWithPlacement is the FROM clause shared by every registration
// query: the registration, its candidate and its derived placement.
func registrationsWithPlacement() ReadableTable {
	return Registrations.
		INNER_JOIN(Candidates, Candidates.ID.EQ(Registrations.CandidateID)).
		LEFT_JOIN(view.RegistrationPlacements, view.RegistrationPlacements.RegistrationID.EQ(Registrations.ID))
}

func fetchRegistrations(db qrm.DB, where BoolExpression) ([]*types.Registration, error) {
	stmt := SELECT(
		Registrations.AllColumns,
		Candidates.AllColumns,
		view.RegistrationPlacements.AllColumns,
	).FROM(
		registrationsWithPlacement(),
	).WHERE(
		where,
	).ORDER_BY(
		Registrations.ID.ASC(),
	)

	var result registrationsResult

	err := stmt.Query(db, &result)
	if err != nil {
		return nil, err
	}

	return result.toRegistrationsDomain(), nil
}

func (d *Database) FetchRegistrations() ([]*types.Registration, error) {
	return fetchRegistrations(d.db, Bool(true))
}

func (d *Database) FetchRegistrationsBySelectionID(selectionID int32) ([]*types.Registration, error) {
	return fetchRegistrations(d.db, Registrations.SelectionID.EQ(Int32(selectionID)))
}

func (d *Database) FetchRegistrationsByCourseID(courseID int32) ([]*types.Registration, error) {
	return fetchRegistrations(d.db, Registrations.CourseID.EQ(Int32(courseID)))
}

// FetchRegistrationsByIDs returns the given registrations (any order).
func FetchRegistrationsByIDs(db qrm.DB, ids []int32) ([]*types.Registration, error) {
	if len(ids) == 0 {
		return []*types.Registration{}, nil
	}

	return fetchRegistrations(db, Registrations.ID.IN(int32Exprs(ids)...))
}

// FetchEnrolledRegistrationsByCourseAndSemester returns the students currently
// enrolled in a course for one semester.
func (d *Database) FetchEnrolledRegistrationsByCourseAndSemester(
	courseID int32,
	semester int32,
) ([]*types.Registration, error) {
	return fetchRegistrations(d.db,
		Registrations.CourseID.EQ(Int32(courseID)).
			AND(view.RegistrationPlacements.Outcome.EQ(String(types.CallEntryOutcomeEnrolled.String()))).
			AND(view.RegistrationPlacements.Semester.EQ(Int32(semester))),
	)
}

func fullRegistrationSelect() SelectStatement {
	return SELECT(
		Registrations.AllColumns,
		Candidates.AllColumns,
		view.RegistrationPlacements.AllColumns,
		Courses.AllColumns,
		Quotas.AllColumns,
	).FROM(
		registrationsWithPlacement().
			INNER_JOIN(Courses, Courses.ID.EQ(Registrations.CourseID)).
			INNER_JOIN(Quotas, Quotas.ID.EQ(Courses.QuotaID)),
	)
}

func (d *Database) FetchEnrolledRegistrationDetails() ([]*types.RegistrationDetail, error) {
	stmt := fullRegistrationSelect().WHERE(
		view.RegistrationPlacements.Outcome.EQ(String(types.CallEntryOutcomeEnrolled.String())),
	).ORDER_BY(
		Registrations.ID.ASC(),
	)

	var result fullRegistrationsResult

	err := stmt.Query(d.db, &result)
	if err != nil {
		return nil, err
	}

	return result.toRegistrationDetails()
}

// FetchRegistrationByID returns a registration with its course and full call history.
func FetchRegistrationByID(db qrm.DB, registrationID int32) (*types.RegistrationDetail, error) {
	stmt := fullRegistrationSelect().WHERE(
		Registrations.ID.EQ(Int32(registrationID)),
	)

	var result fullRegistrationResult

	err := stmt.Query(db, &result)
	if err != nil {
		return nil, err
	}

	detail, err := result.toRegistrationDetail()
	if err != nil {
		return nil, err
	}

	detail.History, err = FetchRegistrationHistory(db, registrationID)
	if err != nil {
		return nil, err
	}

	return detail, nil
}

func (d *Database) FetchRegistrationByID(registrationID int32) (*types.RegistrationDetail, error) {
	return FetchRegistrationByID(d.db, registrationID)
}

func (d *Database) FetchSelectionKindByRegistrationID(registrationID int32) (types.SelectionKind, error) {
	stmt := SELECT(
		Selections.Kind,
	).FROM(
		Registrations.
			INNER_JOIN(Selections, Selections.ID.EQ(Registrations.SelectionID)),
	).WHERE(
		Registrations.ID.EQ(Int32(registrationID)),
	)

	var result struct {
		Kind string
	}

	err := stmt.Query(d.db, &result)
	if err != nil {
		return "", err
	}

	return types.ParseSelectionKind(result.Kind)
}

func DeleteAllRegistrations(db qrm.DB) error {
	stmt := Registrations.DELETE().
		WHERE(Bool(true))

	_, err := stmt.Exec(db)
	return err
}

func FetchCandidateIDsBySelectionID(db qrm.DB, selectionID int32) ([]int32, error) {
	stmt := SELECT(
		Registrations.CandidateID,
	).DISTINCT().FROM(
		Registrations,
	).WHERE(
		Registrations.SelectionID.EQ(Int32(selectionID)),
	)

	var candidateIDs []int32
	err := stmt.Query(db, &candidateIDs)
	if err != nil {
		return nil, err
	}

	return candidateIDs, nil
}

func int32Exprs(values []int32) []Expression {
	exprs := make([]Expression, len(values))
	for i, v := range values {
		exprs[i] = Int32(v)
	}

	return exprs
}
