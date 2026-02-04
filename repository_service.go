package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	sqlite3 "modernc.org/sqlite"
)

type RepositoryService struct {
	shouldDestroy bool
	file          string
	db            *sqlx.DB
}

func NewRepositoryService(file string, db *sqlx.DB) *RepositoryService {
	repositoryService := RepositoryService{
		file:          file,
		db:            db,
		shouldDestroy: false,
	}

	return &repositoryService
}

func (rs *RepositoryService) Destroy() {
	rs.shouldDestroy = true
}

func (rs *RepositoryService) Close() error {
	if rs.shouldDestroy {
		defer rs.destroy()
	}

	return rs.db.Close()
}

type backuper interface {
	NewBackup(string) (*sqlite3.Backup, error)
	NewRestore(string) (*sqlite3.Backup, error)
}

func (rs *RepositoryService) Backup(file string) error {
	conn, err := rs.db.Conn(context.Background())
	if err != nil {
		return fmt.Errorf("db conn: %w", err)
	}
	defer conn.Close()
	err = conn.Raw(func(driverConn any) error {
		bck, err := driverConn.(backuper).NewBackup(file)
		if err != nil {
			return fmt.Errorf("NewBackup: %w", err)
		}

		for more := true; more; {
			more, err = bck.Step(-1)
			if err != nil {
				return fmt.Errorf("bkp step: %w", err)
			}
		}

		return bck.Finish()
	})

	if err != nil {
		return fmt.Errorf("bkp: %w", err)
	}

	return nil
}

func (rs *RepositoryService) Restore(file string) error {
	conn, err := rs.db.Conn(context.Background())
	if err != nil {
		return fmt.Errorf("db conn: %w", err)
	}
	defer conn.Close()
	err = conn.Raw(func(driverConn any) error {
		bck, err := driverConn.(backuper).NewRestore(file)
		if err != nil {
			return fmt.Errorf("NewRestore: %w", err)
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
		return fmt.Errorf("bkp: %w", err)
	}

	return nil
}

func (rs *RepositoryService) destroy() {
	os.Remove(rs.file)
}
