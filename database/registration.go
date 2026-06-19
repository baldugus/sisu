package database

import (
	"github.com/baldugus/sisu/database/.gen/model"
	. "github.com/baldugus/sisu/database/.gen/table"
	"github.com/baldugus/sisu/types"
	"github.com/go-jet/jet/v2/qrm"
	. "github.com/go-jet/jet/v2/sqlite"
)

type CreateRegistrationArgs struct {
	Registration *types.Registration
	CandidateID  int32
	CourseID     int32
	SelectionID  int32
	CallID       *int32
	SemesterID   *int32
}

func CreateRegistration(db qrm.DB, args *CreateRegistrationArgs) error {
	registrationModel := toRegistrationModel(args.Registration)
	registrationModel.CandidateID = args.CandidateID
	registrationModel.CourseID = args.CourseID
	registrationModel.SelectionID = args.SelectionID
	registrationModel.CallID = args.CallID
	registrationModel.SemesterID = args.SemesterID

	stmt := Registrations.INSERT(Registrations.MutableColumns).
		MODEL(registrationModel).
		RETURNING(Registrations.ID)

	_, err := insertOne[model.Registrations](db, stmt)
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) FetchRegistrations() ([]*types.Registration, error) {
	stmt := SELECT(
		Registrations.AllColumns,
		Candidates.AllColumns,
	).FROM(
		Registrations.
			INNER_JOIN(Candidates, Candidates.ID.EQ(Registrations.CandidateID)),
	)

	var result registrationsResult

	err := stmt.Query(d.db, &result)
	if err != nil {
		return nil, err
	}

	return result.toRegistrationsDomain(), nil
}

func (d *Database) FetchRegistrationsBySelectionID(selectionID int32) ([]*types.Registration, error) {
	stmt := SELECT(
		Registrations.AllColumns,
		Candidates.AllColumns,
	).FROM(
		Registrations.
			INNER_JOIN(Candidates, Candidates.ID.EQ(Registrations.CandidateID)),
	).WHERE(
		Registrations.SelectionID.EQ(Int32(selectionID)),
	)

	var result registrationsResult

	err := stmt.Query(d.db, &result)
	if err != nil {
		return nil, err
	}

	return result.toRegistrationsDomain(), nil
}

func (d *Database) FetchRegistrationsByCourseID(courseID int32) ([]*types.Registration, error) {
	stmt := SELECT(
		Registrations.AllColumns,
		Candidates.AllColumns,
	).FROM(
		Registrations.
			INNER_JOIN(Candidates, Candidates.ID.EQ(Registrations.CandidateID)),
	).WHERE(
		Registrations.CourseID.EQ(Int32(courseID)),
	)

	var result registrationsResult

	err := stmt.Query(d.db, &result)
	if err != nil {
		return nil, err
	}

	return result.toRegistrationsDomain(), nil
}

func (d *Database) FetchRegistrationsByCallID(callID int32) ([]*types.Registration, error) {
	stmt := SELECT(
		Registrations.AllColumns,
		Candidates.AllColumns,
	).FROM(
		Registrations.
			INNER_JOIN(Candidates, Candidates.ID.EQ(Registrations.CandidateID)),
	).WHERE(
		Registrations.CallID.EQ(Int32(callID)),
	)

	var result registrationsResult

	err := stmt.Query(d.db, &result)
	if err != nil {
		return nil, err
	}

	return result.toRegistrationsDomain(), nil
}

func (d *Database) FetchRegistrationsByCallIDAndCourseID(callID, courseID int32) ([]*types.Registration, error) {
	stmt := SELECT(
		Registrations.AllColumns,
		Candidates.AllColumns,
	).FROM(
		Registrations.
			INNER_JOIN(Candidates, Candidates.ID.EQ(Registrations.CandidateID)),
	).WHERE(
		Registrations.CallID.EQ(Int32(callID)).
			AND(Registrations.CourseID.EQ(Int32(courseID))),
	)

	var result registrationsResult

	err := stmt.Query(d.db, &result)
	if err != nil {
		return nil, err
	}

	return result.toRegistrationsDomain(), nil
}

func (d *Database) FetchEnrolledRegistrationsByCourseID(courseID int32) ([]*types.Registration, error) {
	stmt := SELECT(
		Registrations.AllColumns,
		Candidates.AllColumns,
	).FROM(
		Registrations.
			INNER_JOIN(Candidates, Candidates.ID.EQ(Registrations.CandidateID)),
	).WHERE(
		Registrations.CourseID.EQ(Int32(courseID)).
			AND(Registrations.Status.EQ(String(types.RegistrationStatusEnrolled.String()))),
	)

	var result registrationsResult

	err := stmt.Query(d.db, &result)
	if err != nil {
		return nil, err
	}

	return result.toRegistrationsDomain(), nil
}

func (d *Database) FetchEnrolledRegistrationsByCourseIDs(courseIDs []int32) ([]*types.Registration, error) {
	if len(courseIDs) == 0 {
		return []*types.Registration{}, nil
	}

	courseIDExprs := make([]Expression, len(courseIDs))
	for i, id := range courseIDs {
		courseIDExprs[i] = Int32(id)
	}

	stmt := SELECT(
		Registrations.AllColumns,
		Candidates.AllColumns,
	).FROM(
		Registrations.
			INNER_JOIN(Candidates, Candidates.ID.EQ(Registrations.CandidateID)),
	).WHERE(
		Registrations.CourseID.IN(courseIDExprs...).
			AND(Registrations.Status.EQ(String(types.RegistrationStatusEnrolled.String()))),
	)

	var result registrationsResult

	err := stmt.Query(d.db, &result)
	if err != nil {
		return nil, err
	}

	return result.toRegistrationsDomain(), nil
}

func (d *Database) FetchEnrolledRegistrationDetails() ([]*types.RegistrationDetail, error) {
	stmt := SELECT(
		Registrations.AllColumns,
		Candidates.AllColumns,
		Courses.AllColumns,
		Quotas.AllColumns,
		Calls.AllColumns,
	).FROM(
		Registrations.
			INNER_JOIN(Candidates, Candidates.ID.EQ(Registrations.CandidateID)).
			INNER_JOIN(Courses, Courses.ID.EQ(Registrations.CourseID)).
			INNER_JOIN(Quotas, Quotas.ID.EQ(Courses.QuotaID)).
			LEFT_JOIN(Calls, Calls.ID.EQ(Registrations.CallID)),
	).WHERE(
		Registrations.Status.EQ(String(types.RegistrationStatusEnrolled.String())),
	)

	var result fullRegistrationsResult

	err := stmt.Query(d.db, &result)
	if err != nil {
		return nil, err
	}

	return result.toRegistrationDetails(), nil
}

func (d *Database) FetchRegistrationByID(registrationID int32) (*types.RegistrationDetail, error) {
	stmt := SELECT(
		Registrations.AllColumns,
		Candidates.AllColumns,
		Courses.AllColumns,
		Quotas.AllColumns,
		Calls.AllColumns,
	).FROM(
		Registrations.
			INNER_JOIN(Candidates, Candidates.ID.EQ(Registrations.CandidateID)).
			INNER_JOIN(Courses, Courses.ID.EQ(Registrations.CourseID)).
			INNER_JOIN(Quotas, Quotas.ID.EQ(Courses.QuotaID)).
			LEFT_JOIN(Calls, Calls.ID.EQ(Registrations.CallID)),
	).WHERE(
		Registrations.ID.EQ(Int32(registrationID)),
	)

	var result fullRegistrationResult

	err := stmt.Query(d.db, &result)
	if err != nil {
		return nil, err
	}

	return result.toRegistrationDetail(), nil
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

func CountCourseOccupiedSeats(db qrm.DB, courseID int32) (int32, error) {
	stmt := SELECT(
		COUNT(Registrations.ID).AS("count"),
	).FROM(
		Registrations,
	).WHERE(
		Registrations.CourseID.EQ(Int32(courseID)).
			AND(Registrations.Status.IN(
				String(types.RegistrationStatusApproved.String()),
				String(types.RegistrationStatusEnrolled.String()),
			)),
	)

	var result struct {
		Count int32
	}

	err := stmt.Query(db, &result)
	if err != nil {
		return 0, err
	}

	return result.Count, nil
}

func (d *Database) UpdateRegistrationStatus(registrationID int32, status types.RegistrationStatus) error {
	stmt := Registrations.UPDATE().
		SET(
			Registrations.Status.SET(String(status.String())),
		).
		WHERE(Registrations.ID.EQ(Int32(registrationID)))

	_, err := stmt.Exec(d.db)
	return err
}

func FetchWaitlistedRegistrationsByCourse(db qrm.DB, courseID int32, limit int32) ([]int32, error) {
	stmt := SELECT(
		Registrations.ID,
	).FROM(
		Registrations,
	).WHERE(
		Registrations.CourseID.EQ(Int32(courseID)).
			AND(Registrations.Status.EQ(String(types.RegistrationStatusWaitlisted.String()))),
	).ORDER_BY(
		Registrations.Ranking.DESC(),
	).LIMIT(int64(limit))

	var result []int32

	err := stmt.Query(db, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func FetchPriorityRegistrationsByCourse(db qrm.DB, courseID int32, semesterID int32, limit int32) ([]int32, error) {
	stmt := SELECT(
		Registrations.ID,
	).FROM(
		Registrations,
	).WHERE(
		Registrations.CourseID.EQ(Int32(courseID)).
			AND(Registrations.Status.EQ(String(types.RegistrationStatusApproved.String()))).
			AND(Registrations.SemesterID.EQ(Int32(semesterID))),
	).ORDER_BY(
		Registrations.Ranking.ASC(),
	).LIMIT(int64(limit))

	var result []int32

	err := stmt.Query(db, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func AssignRegistrationToCall(db qrm.DB, registrationID int32, callID int32) error {
	stmt := Registrations.UPDATE().
		SET(
			Registrations.Status.SET(String(types.RegistrationStatusApproved.String())),
			Registrations.CallID.SET(Int32(callID)),
		).
		WHERE(Registrations.ID.EQ(Int32(registrationID)))

	_, err := stmt.Exec(db)
	return err
}

func RevertRegistrationsToWaitlisted(db qrm.DB, callID int32) error {
	stmt := Registrations.UPDATE().
		SET(
			Registrations.Status.SET(String(types.RegistrationStatusWaitlisted.String())),
			Registrations.CallID.SET(CAST(Raw("NULL")).AS_INTEGER()),
		).
		WHERE(Registrations.CallID.EQ(Int32(callID)))

	_, err := stmt.Exec(db)
	return err
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
