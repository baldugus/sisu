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

func DeleteCandidatesBySelectionID(db qrm.DB, selectionID int32) error {
	// Delete candidates that are linked to registrations in this selection
	stmt := Candidates.DELETE().
		WHERE(Candidates.ID.IN(
			SELECT(Registrations.CandidateID).
				FROM(Registrations).
				WHERE(Registrations.SelectionID.EQ(Int32(selectionID))),
		))

	_, err := stmt.Exec(db)
	return err
}

func DeleteCandidatesByIDs(db qrm.DB, candidateIDs []int32) error {
	if len(candidateIDs) == 0 {
		return nil
	}

	// Convert []int32 to []Expression for IN clause
	ids := make([]Expression, len(candidateIDs))
	for i, id := range candidateIDs {
		ids[i] = Int32(id)
	}

	stmt := Candidates.DELETE().
		WHERE(Candidates.ID.IN(ids...))

	_, err := stmt.Exec(db)
	return err
}
