package database

import (
	"github.com/baldugus/sisu/database/.gen/model"
	. "github.com/baldugus/sisu/database/.gen/table"
	"github.com/go-jet/jet/v2/qrm"
	. "github.com/go-jet/jet/v2/sqlite"
)

func CreateQuota(db qrm.DB, quota string) (int32, error) {
	quotaModel := model.Quotas{Name: quota}

	stmt := Quotas.INSERT(Quotas.MutableColumns).
		MODEL(quotaModel).
		ON_CONFLICT(Quotas.Name).DO_UPDATE(
		SET(
			Quotas.ID.SET(Quotas.ID),
		)).
		RETURNING(Quotas.ID)

	result, err := insertOne[model.Quotas](db, stmt)
	if err != nil {
		return 0, err
	}

	return result.ID, nil
}

func DeleteAllQuotas(db qrm.DB) error {
	stmt := Quotas.DELETE().
		WHERE(Bool(true))

	_, err := stmt.Exec(db)
	return err
}
