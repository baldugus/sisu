package repository

import (
	"database/sql"
	"fmt"

	"changeme/types"
)

type SelectionRepository struct {
	db *DB
}

func NewSelectionRepository(db *DB) *SelectionRepository {
	return &SelectionRepository{
		db: db,
	}
}

func (s *SelectionRepository) Begin() error {
	return s.db.Begin()
}

func (s *SelectionRepository) Commit() error {
	return s.db.Commit()
}

func (s *SelectionRepository) Rollback() error {
	return s.db.Rollback()
}

func (s *SelectionRepository) FindSelectionByID(id int64) (*types.Selection, error) {
	var selection types.Selection

	query := "SELECT * FROM selections WHERE id = ?"
	if err := s.db.Get(&selection, query, id); err != nil {
		return nil, fmt.Errorf("db get: %w", err)
	}

	return &selection, nil
}

func (s *SelectionRepository) FindSelectionByKind(kind types.SelectionKind) (*types.Selection, error) {
	var selection types.Selection

	query := "SELECT * FROM selections WHERE kind = ?"
	if err := s.db.Get(&selection, query, kind); err != nil {
		return nil, fmt.Errorf("db get: %w", err)
	}

	return &selection, nil
}

func (s *SelectionRepository) FindSelections() ([]*types.Selection, error) {
	var selections []*types.Selection

	query := "SELECT * FROM selections"
	if err := s.db.Select(&selections, query); err != nil {
		return nil, fmt.Errorf("db select: %w", err)
	}

	if len(selections) == 0 {
		return nil, sql.ErrNoRows
	}

	return selections, nil
}

func (s *SelectionRepository) CreateSelection(selection *types.Selection) error {
	if err := types.CanCreateSelection(s, selection.Kind); err != nil {
		return fmt.Errorf("can't create selection: %w", err)
	}

	query := "INSERT INTO selections (name, kind, date, institution, course) VALUES (:name, :kind, :date, :institution, :course)"

	result, err := s.db.NamedExec(query, selection)
	if err != nil {
		return fmt.Errorf("db named exec: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("last insert id: %w", err)
	}

	selection.ID = id

	return nil
}

func (s *SelectionRepository) DeleteSelection(id int64) error {
	rollcallRepo := NewRollcallRepository(s.db)

	applicationRepo := NewApplicationRepository(s.db)
	if err := types.CanDeleteSelection(s, rollcallRepo, applicationRepo, id); err != nil {
		return fmt.Errorf("can delete selection: %w", err)
	}

	query := "DELETE FROM selections WHERE ID = ?"

	_, err := s.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("db exec: %w", err)
	}

	return nil
}
