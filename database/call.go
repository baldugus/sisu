package database

import (
	"github.com/baldugus/sisu/database/.gen/model"
	. "github.com/baldugus/sisu/database/.gen/table"
	"github.com/baldugus/sisu/types"
	"github.com/go-jet/jet/v2/qrm"
	. "github.com/go-jet/jet/v2/sqlite"
)

func (d *Database) FetchCalls() ([]*types.Call, error) {
	stmt := SELECT(
		Calls.AllColumns,
	).FROM(
		Calls,
	)

	var result callsResult

	err := stmt.Query(d.db, &result)
	if err != nil {
		return nil, err
	}

	return result.toCallsDomain(), nil
}

func (d *Database) FetchCallByID(callID int32) (*types.Call, error) {
	stmt := SELECT(
		Calls.AllColumns,
	).FROM(
		Calls,
	).WHERE(
		Calls.ID.EQ(Int32(callID)),
	)

	var result model.Calls

	err := stmt.Query(d.db, &result)
	if err != nil {
		return nil, err
	}

	return toCallDomain(&result), nil
}

func CreateCall(db qrm.DB, call *types.Call) (int32, error) {
	callModel := toCallModel(call)
	stmt := Calls.INSERT(Calls.MutableColumns).
		MODEL(callModel).
		RETURNING(Calls.ID)

	result, err := insertOne[model.Calls](db, stmt)
	if err != nil {
		return 0, err
	}

	return result.ID, nil
}

func (d *Database) CallHasPendingRegistrations(ID int32) (bool, error) {
	stmt := SELECT(
		Registrations.ID,
	).FROM(
		Registrations.
			INNER_JOIN(Calls, Calls.ID.EQ(Registrations.CallID)),
	).WHERE(
		Calls.ID.EQ(Int32(ID)).
			AND(Registrations.Status.IN(
				String(types.RegistrationStatusApproved.String()),
				String(types.RegistrationStatusWaitlisted.String()),
			)),
	)

	var result []int32

	err := stmt.Query(d.db, &result)
	if err != nil {
		return false, err
	}

	if len(result) > 0 {
		return true, nil
	}

	return false, nil
}

func (d *Database) CloseCall(ID int32) error {
	stmt := Calls.UPDATE().
		SET(
			Calls.Status.SET(String(types.CallStatusDone.String())),
		).
		WHERE(Calls.ID.EQ(Int32(ID)))

	_, err := stmt.Exec(d.db)
	return err
}

func (d *Database) HasOpenCall() (bool, error) {
	stmt := SELECT(
		Calls.ID,
	).FROM(
		Calls,
	).WHERE(
		Calls.Status.EQ(String(types.CallStatusCalling.String())),
	).LIMIT(1)

	var result []int32

	err := stmt.Query(d.db, &result)
	if err != nil {
		return false, err
	}

	return len(result) > 0, nil
}

func (d *Database) GetLastCallNumber() (int32, error) {
	stmt := SELECT(
		MAX(Calls.Number).AS("max_number"),
	).FROM(
		Calls,
	)

	var result struct {
		MaxNumber *int32
	}

	err := stmt.Query(d.db, &result)
	if err != nil {
		return 0, err
	}

	if result.MaxNumber == nil {
		return 0, nil
	}

	return *result.MaxNumber, nil
}

func (d *Database) GetCallNumber(ID int32) (int32, error) {
	stmt := SELECT(
		Calls.AllColumns,
	).FROM(
		Calls,
	).WHERE(
		Calls.ID.EQ(Int32(ID)),
	)

	var result model.Calls

	err := stmt.Query(d.db, &result)
	if err != nil {
		return 0, err
	}

	return result.Number, nil
}

func (d *Database) HasCallAfterNumber(number int32) (bool, error) {
	stmt := SELECT(
		Calls.ID,
	).FROM(
		Calls,
	).WHERE(
		Calls.Number.GT(Int32(number)),
	).LIMIT(1)

	var result []int32

	err := stmt.Query(d.db, &result)
	if err != nil {
		return false, err
	}

	return len(result) > 0, nil
}

func (d *Database) OpenCall(ID int32) error {
	stmt := Calls.UPDATE().
		SET(
			Calls.Status.SET(String(types.CallStatusCalling.String())),
		).
		WHERE(Calls.ID.EQ(Int32(ID)))

	_, err := stmt.Exec(d.db)
	return err
}

func DeleteCall(db qrm.DB, ID int32) error {
	stmt := Calls.DELETE().
		WHERE(Calls.ID.EQ(Int32(ID)))

	_, err := stmt.Exec(db)
	return err
}

func (d *Database) FetchCallByNumber(callNumber int32) (*types.Call, error) {
	stmt := SELECT(
		Calls.AllColumns,
	).FROM(
		Calls,
	).WHERE(
		Calls.Number.EQ(Int32(callNumber)),
	)

	var result model.Calls

	err := stmt.Query(d.db, &result)
	if err != nil {
		return nil, err
	}

	return toCallDomain(&result), nil
}

func DeleteAllCalls(db qrm.DB) error {
	stmt := Calls.DELETE().
		WHERE(Bool(true))

	_, err := stmt.Exec(db)
	return err
}
