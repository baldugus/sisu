package database

import (
	"github.com/baldugus/sisu/database/.gen/model"
	. "github.com/baldugus/sisu/database/.gen/table"
	"github.com/baldugus/sisu/types"
	"github.com/go-jet/jet/v2/qrm"
	. "github.com/go-jet/jet/v2/sqlite"
)

func CreateCourse(db qrm.DB, course *types.Course, quotaID int32) (int32, error) {
	courseModel := toCourseModel(course)
	courseModel.QuotaID = quotaID
	stmt := Courses.INSERT(Courses.MutableColumns).
		MODEL(courseModel).
		ON_CONFLICT(Courses.TimeSlot, Courses.QuotaID).DO_UPDATE(
		SET(
			Courses.ID.SET(Courses.ID),
		)).
		RETURNING(Courses.ID)

	result, err := insertOne[model.Courses](db, stmt)
	if err != nil {
		return 0, err
	}

	return result.ID, nil
}

type courseResult struct {
	model.Courses
	Quota model.Quotas
}

type coursesResult []courseResult

func (c coursesResult) toCoursesDomain() []*types.Course {
	courses := make([]*types.Course, len(c))
	for i := range c {
		courses[i] = toCourseDomain(&c[i].Courses, c[i].Quota.Name)
	}
	return courses
}

func (d *Database) FetchCourses() ([]*types.Course, error) {
	stmt := SELECT(
		Courses.AllColumns,
		Quotas.AllColumns,
	).FROM(
		Courses.
			INNER_JOIN(Quotas, Quotas.ID.EQ(Courses.QuotaID)),
	)

	var result coursesResult

	err := stmt.Query(d.db, &result)
	if err != nil {
		return nil, err
	}

	return result.toCoursesDomain(), nil
}

func (d *Database) FetchCourseByID(courseID int32) (*types.Course, error) {
	stmt := SELECT(
		Courses.AllColumns,
		Quotas.AllColumns,
	).FROM(
		Courses.
			INNER_JOIN(Quotas, Quotas.ID.EQ(Courses.QuotaID)),
	).WHERE(
		Courses.ID.EQ(Int32(courseID)),
	)

	var result courseResult

	err := stmt.Query(d.db, &result)
	if err != nil {
		return nil, err
	}

	return toCourseDomain(&result.Courses, result.Quota.Name), nil
}

func (d *Database) FetchCoursesByPeriod(period types.CoursePeriod) ([]*types.Course, error) {
	stmt := SELECT(
		Courses.AllColumns,
		Quotas.AllColumns,
	).FROM(
		Courses.
			INNER_JOIN(Quotas, Quotas.ID.EQ(Courses.QuotaID)),
	).WHERE(
		Courses.TimeSlot.EQ(String(period.String())),
	)

	var result coursesResult

	err := stmt.Query(d.db, &result)
	if err != nil {
		return nil, err
	}

	return result.toCoursesDomain(), nil
}

func DeleteAllCourses(db qrm.DB) error {
	stmt := Courses.DELETE().
		WHERE(Bool(true))

	_, err := stmt.Exec(db)
	return err
}
