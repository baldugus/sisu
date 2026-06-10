package database

import (

	"github.com/baldugus/sisu/database/.gen/model"
	. "github.com/baldugus/sisu/database/.gen/table"
	"github.com/baldugus/sisu/types"
	"github.com/go-jet/jet/v2/qrm"
	. "github.com/go-jet/jet/v2/sqlite"
)

func CreateSemester(db qrm.DB, year int32, number int32) (int32, error) {
	insertStmt := Semesters.INSERT(Semesters.Year, Semesters.Number, Semesters.Status).
		MODEL(model.Semesters{
			Year:   year,
			Number: number,
			Status: "open",
		}).
		RETURNING(Semesters.ID)

	var insertResult model.Semesters
	err := insertStmt.Query(db, &insertResult)
	if err != nil {
		return 0, err
	}

	return insertResult.ID, nil
}

func FetchSemesterByID(db qrm.DB, id int32) (*types.Semester, error) {
	stmt := SELECT(Semesters.AllColumns).FROM(Semesters).WHERE(Semesters.ID.EQ(Int32(id)))
	var result model.Semesters
	err := stmt.Query(db, &result)
	if err != nil {
		return nil, err
	}
	return &types.Semester{
		ID:     result.ID,
		Year:   result.Year,
		Number: result.Number,
		Status: types.SemesterStatusOpen, // TODO parse properly
	}, nil
}

func FetchSemesterByYearAndNumber(db qrm.DB, year int32, number int32) (*types.Semester, error) {
	stmt := SELECT(Semesters.AllColumns).FROM(Semesters).WHERE(
		Semesters.Year.EQ(Int32(year)).AND(Semesters.Number.EQ(Int32(number))),
	)
	var result model.Semesters
	err := stmt.Query(db, &result)
	if err != nil {
		return nil, err
	}
	return &types.Semester{
		ID:     result.ID,
		Year:   result.Year,
		Number: result.Number,
	}, nil
}

func (d *Database) FetchSemesters() ([]*types.Semester, error) {
	stmt := SELECT(Semesters.AllColumns).FROM(Semesters).ORDER_BY(Semesters.Number.ASC())
	var results []model.Semesters
	err := stmt.Query(d.db, &results)
	if err != nil {
		return nil, err
	}
	semesters := make([]*types.Semester, len(results))
	for i, r := range results {
		semesters[i] = &types.Semester{
			ID:     r.ID,
			Year:   r.Year,
			Number: r.Number,
		}
	}
	return semesters, nil
}
