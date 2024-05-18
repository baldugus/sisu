package repository

import (
	"fmt"

	"changeme/types"

	"changeme/internal/sql"

	"github.com/jmoiron/sqlx"
)

type db interface {
	Select(dest interface{}, query string, args ...interface{}) error
}

type Repository struct {
	*sqlx.DB
}

type RepositoryTx struct {
	*sqlx.Tx
	*Repository
}

func (r *Repository) Begin() (*RepositoryTx, error) {
	tx, err := r.DB.Beginx()
	if err != nil {
		return nil, err
	}

	return &RepositoryTx{Tx: tx, Repository: r}, nil
}
