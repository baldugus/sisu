package database

import (
	"github.com/baldugus/sisu/database/.gen/model"
	. "github.com/baldugus/sisu/database/.gen/table"
	"github.com/baldugus/sisu/types"
	"github.com/go-jet/jet/v2/qrm"
	. "github.com/go-jet/jet/v2/sqlite"
)

func CreateCandidate(db qrm.DB, candidate *types.Candidate) (int32, error) {
	candidateModel := toCandidateModel(candidate)
	stmt := Candidates.INSERT(Candidates.MutableColumns).
		MODEL(candidateModel).
		RETURNING(Candidates.ID)

	result, err := insertOne[model.Candidates](db, stmt)
	if err != nil {
		return 0, err
	}

	return result.ID, nil
}

func DeleteAllCandidates(db qrm.DB) error {
	stmt := Candidates.DELETE().
		WHERE(Bool(true))

	_, err := stmt.Exec(db)
	return err
}
