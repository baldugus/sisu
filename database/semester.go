package database

import (
	"github.com/baldugus/sisu/database/.gen/model"
	. "github.com/baldugus/sisu/database/.gen/table"
	"github.com/baldugus/sisu/types"
	"github.com/go-jet/jet/v2/qrm"
	. "github.com/go-jet/jet/v2/sqlite"
)

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
		semester := &types.Semester{Number: result[i].Number}

		for _, course := range courses {
			semester.Seats += course.Seats.PerSemester()
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
