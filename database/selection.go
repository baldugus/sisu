package database

import (
	"errors"

	"github.com/baldugus/sisu/database/.gen/model"
	. "github.com/baldugus/sisu/database/.gen/table"
	"github.com/baldugus/sisu/types"
	"github.com/go-jet/jet/v2/qrm"
	. "github.com/go-jet/jet/v2/sqlite"
)

func (d *Database) FetchSelection(kind types.SelectionKind) (*types.Selection, error) {
	stmt := SELECT(
		Selections.AllColumns,
	).FROM(
		Selections,
	).WHERE(
		Selections.Kind.EQ(String(kind.String())),
	)

	var result selectionResult

	err := stmt.Query(d.db, &result)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return nil, ErrNotFound{Err: err}
		}
		return nil, err
	}

	return result.toSelectionDomain(), nil
}

func (d *Database) FetchSelectionKinds() ([]types.SelectionKind, error) {
	stmt := SELECT(
		Selections.Kind,
	).FROM(
		Selections,
	)

	var result []string

	err := stmt.Query(d.db, &result)
	if err != nil {
		return nil, err
	}

	kinds := make([]types.SelectionKind, len(result))

	for i, kindStr := range result {
		kind, _ := types.ParseSelectionKind(kindStr)
		kinds[i] = kind
	}

	return kinds, nil
}

func CreateSelection(db qrm.DB, selection *types.Selection) (int32, error) {
	selectionModel := toSelectionModel(selection)
	stmt := Selections.INSERT(Selections.MutableColumns).
		MODEL(selectionModel).
		RETURNING(Selections.ID)

	result, err := insertOne[model.Selections](db, stmt)
	if err != nil {
		return 0, err
	}

	return result.ID, nil
}

func (d *Database) SelectionHasModifiedRegistrations(kind types.SelectionKind) (bool, error) {
	stmt := SELECT(
		Registrations.ID,
	).FROM(
		Registrations.
			INNER_JOIN(Selections, Selections.ID.EQ(Registrations.SelectionID)),
	).WHERE(
		Selections.Kind.EQ(String(kind.String())).
			AND(Registrations.Status.IN(
				String(types.RegistrationStatusAbsent.String()),
				String(types.RegistrationStatusEnrolled.String()),
			)),
	).LIMIT(1)

	var result []int32

	err := stmt.Query(d.db, &result)
	if err != nil {
		return false, err
	}

	return len(result) > 0, nil
}

func DeleteSelection(db qrm.DB, kind types.SelectionKind) error {
	stmt := Selections.DELETE().
		WHERE(Selections.Kind.EQ(String(kind.String())))

	_, err := stmt.Exec(db)
	return err
}
