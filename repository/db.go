// nolint: wrapcheck,nolintlint
package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

var (
	ErrActiveTransaction = errors.New("already inside a transaction")
	ErrNoTransaction     = errors.New("no transaction active")
)

type DB struct {
	db       *sqlx.DB
	tx       *sqlx.Tx
	activeTx bool
}

func NewDB(db *sqlx.DB) *DB {
	return &DB{
		db:       db,
		tx:       nil,
		activeTx: false,
	}
}

func (d *DB) Begin() error {
	if d.activeTx {
		return ErrActiveTransaction
	}

	tx, err := d.db.Beginx()
	if err != nil {
		return fmt.Errorf("db begin: %w", err)
	}

	d.activeTx = true
	d.tx = tx

	return nil
}

func (d *DB) Commit() error {
	if !d.activeTx {
		return ErrNoTransaction
	}

	err := d.tx.Commit()
	d.tx = nil
	d.activeTx = false

	return err
}

func (d *DB) Rollback() error {
	if !d.activeTx {
		return ErrNoTransaction
	}

	err := d.tx.Rollback()
	d.tx = nil
	d.activeTx = false

	return err
}

func (d *DB) Get(dest any, query string, args ...any) error {
	if d.activeTx {
		return d.tx.Get(dest, query, args...)
	}

	return d.db.Get(dest, query, args...)
}

func (d *DB) NamedExec(query string, arg any) (sql.Result, error) {
	if d.activeTx {
		return d.tx.NamedExec(query, arg)
	}

	return d.db.NamedExec(query, arg)
}

func (d *DB) Exec(query string, args ...any) (sql.Result, error) {
	if d.activeTx {
		return d.tx.Exec(query, args...)
	}

	return d.db.Exec(query, args...)
}

func (d *DB) Select(dest any, query string, args ...any) error {
	if d.activeTx {
		return d.tx.Select(dest, query, args...)
	}

	return d.db.Select(dest, query, args...)
}

func (d *DB) Close() error {
	return d.db.Close()
}
