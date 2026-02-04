package database

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/go-jet/jet/v2/qrm"
	"github.com/go-jet/jet/v2/sqlite"
	"github.com/jmoiron/sqlx"
	sqlite3 "modernc.org/sqlite"
)

type Database struct {
	db       *sqlx.DB
	filePath string
}

func NewDatabase(db *sqlx.DB, filePath string) *Database {
	return &Database{db: db, filePath: filePath}
}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) Destroy() error {
	if err := d.db.Close(); err != nil {
		return fmt.Errorf("close db: %w", err)
	}

	if err := os.Remove(d.filePath); err != nil {
		return fmt.Errorf("remove db file: %w", err)
	}

	return nil
}

func (d *Database) RunInTx(fn func(db qrm.DB) error) error {
	tx, err := d.db.Beginx()
	if err != nil {
		return err
	}

	err = fn(tx)
	if err == nil {
		return tx.Commit()
	}

	rollbackErr := tx.Rollback()
	if rollbackErr != nil {
		return errors.Join(err, rollbackErr)
	}

	return err
}

func insertOne[T any](db qrm.DB, stmt sqlite.InsertStatement) (T, error) {
	var dest []T

	err := stmt.Query(db, &dest)
	if err != nil {
		query, _ := stmt.Sql()
		return *new(T), ErrQuery{Query: query, Err: err}
	}

	if len(dest) != 1 {
		return *new(T), ErrUnexpectedNumOfRowsAffected{Expected: 1, Actual: len(dest)}
	}

	return dest[0], nil
}

type backuper interface {
	NewBackup(string) (*sqlite3.Backup, error)
	NewRestore(string) (*sqlite3.Backup, error)
}

func (d *Database) Backup(file string) error {
	conn, err := d.db.Conn(context.Background())
	if err != nil {
		return fmt.Errorf("db conn: %w", err)
	}
	defer conn.Close()

	err = conn.Raw(func(driverConn any) error {
		bck, err := driverConn.(backuper).NewBackup(file)
		if err != nil {
			return fmt.Errorf("new backup: %w", err)
		}

		for more := true; more; {
			more, err = bck.Step(-1)
			if err != nil {
				return fmt.Errorf("backup step: %w", err)
			}
		}

		return bck.Finish()
	})

	if err != nil {
		return fmt.Errorf("backup: %w", err)
	}

	return nil
}

func (d *Database) Restore(file string) error {
	conn, err := d.db.Conn(context.Background())
	if err != nil {
		return fmt.Errorf("db conn: %w", err)
	}
	defer conn.Close()

	err = conn.Raw(func(driverConn any) error {
		bck, err := driverConn.(backuper).NewRestore(file)
		if err != nil {
			return fmt.Errorf("new restore: %w", err)
		}

		for more := true; more; {
			more, err = bck.Step(-1)
			if err != nil {
				return fmt.Errorf("restore step: %w", err)
			}
		}

		return bck.Finish()
	})

	if err != nil {
		return fmt.Errorf("restore: %w", err)
	}

	return nil
}
