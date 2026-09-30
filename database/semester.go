package database

import (
	"github.com/baldugus/sisu/database/.gen/model"
	. "github.com/baldugus/sisu/database/.gen/table"
	"github.com/baldugus/sisu/types"
	"github.com/go-jet/jet/v2/qrm"
	. "github.com/go-jet/jet/v2/sqlite"
)

func toSemesterDomain(s *model.Semesters) *types.Semester {
	status := types.SemesterStatusOpen
	if s.ClosedAfterCall != nil {
		status = types.SemesterStatusClosed
	}

	return &types.Semester{
		Number:          s.Number,
		Status:          status,
		ClosedAfterCall: s.ClosedAfterCall,
	}
}

func FetchSemester(db qrm.DB, number int32) (*types.Semester, error) {
	stmt := SELECT(
		Semesters.AllColumns,
	).FROM(
		Semesters,
	).WHERE(
		Semesters.Number.EQ(Int32(number)),
	)

	var result model.Semesters

	err := stmt.Query(db, &result)
	if err != nil {
		return nil, err
	}

	return toSemesterDomain(&result), nil
}

// FetchSemesters returns both semesters with their seat counts and occupancy.
func FetchSemesters(db qrm.DB) ([]*types.Semester, error) {
	stmt := SELECT(
		Semesters.AllColumns,
	).FROM(
		Semesters,
	).ORDER_BY(
		Semesters.Number.ASC(),
	)

	var result []model.Semesters

	err := stmt.Query(db, &result)
	if err != nil {
		return nil, err
	}

	courses, err := fetchCourses(db)
	if err != nil {
		return nil, err
	}

	occupied, err := FetchOccupiedSeats(db)
	if err != nil {
		return nil, err
	}

	semesters := make([]*types.Semester, len(result))
	for i := range result {
		semester := toSemesterDomain(&result[i])

		for _, course := range courses {
			semester.Seats += types.SemesterSeats(course.Seats, semester.Number)
		}

		for _, o := range occupied {
			if o.Semester == semester.Number {
				semester.Occupied += o.Count
			}
		}

		semesters[i] = semester
	}

	return semesters, nil
}

func (d *Database) FetchSemesters() ([]*types.Semester, error) {
	return FetchSemesters(d.db)
}

// SetSemesterClosedAfterCall closes a semester after the given call number, or
// reopens it when callNumber is nil.
func SetSemesterClosedAfterCall(db qrm.DB, number int32, callNumber *int32) error {
	value := Expression(NULL)
	if callNumber != nil {
		value = Int32(*callNumber)
	}

	stmt := Semesters.UPDATE().
		SET(
			Semesters.ClosedAfterCall.SET(IntExp(value)),
		).
		WHERE(Semesters.Number.EQ(Int32(number)))

	_, err := stmt.Exec(db)
	return err
}

// ReopenAllSemesters clears every closure; used when the cycle is reset.
func ReopenAllSemesters(db qrm.DB) error {
	stmt := Semesters.UPDATE().
		SET(
			Semesters.ClosedAfterCall.SET(IntExp(NULL)),
		).
		WHERE(Bool(true))

	_, err := stmt.Exec(db)
	return err
}
