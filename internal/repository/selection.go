package repository

import (
	"fmt"

	"changeme/internal/sql"
	"changeme/types"
)

func (r *Repository) FetchSelections() ([]*types.Selection, error) {
	query := sql.Query("select_selections")

	var selections []*types.Selection
	if err := r.Select(&selections, query); err != nil {
		return nil, fmt.Errorf("select selections: %w", err)
	}

	return selections, nil
}

func (r *Repository) SaveSelection(selection *types.Selection) (int64, error) {
	query := sql.Query("insert_selection")

	res, err := r.NamedExec(query, selection)
	if err != nil {
		return 0, fmt.Errorf("insert selection: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("insert selection id: %w", err)
	}

	return id, nil
}
