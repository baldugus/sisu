package repository

import (
	"fmt"
	"go/types"

	"changeme/internal/sql"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	*sqlx.DB
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
