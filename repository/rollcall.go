package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"changeme/types"
)

type RollcallRepository struct {
	db *DB
}

func NewRollcallRepository(db *DB) *RollcallRepository {
	return &RollcallRepository{
		db: db,
	}
}

func (r *RollcallRepository) Begin() error {
	return r.db.Begin()
}

func (r *RollcallRepository) Commit() error {
	return r.db.Commit()
}

func (r *RollcallRepository) Rollback() error {
	return r.db.Rollback()
}

func (r *RollcallRepository) Close() error {
	return r.db.Close()
}

func (r *RollcallRepository) CreateRollcall() (*types.Rollcall, error) {
	if err := types.CanCreateRollcall(r); err != nil {
		return nil, fmt.Errorf("can't create rollcall: %w", err)
	}

	rollcalls, err := r.FindRollcalls(types.RollcallsFilter{})
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	rollcall := types.Rollcall{
		ID:     nil,
		Number: int64(len(rollcalls) + 1),
		Status: "CALLING",
	}

	query := "INSERT INTO rollcalls (number, status) VALUES (:number, :status)"

	result, err := r.db.NamedExec(query, rollcall)
	if err != nil {
		return nil, fmt.Errorf("db named exec: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("last insert id: %w", err)
	}

	rollcall.ID = &id

	return &rollcall, nil
}

func (r *RollcallRepository) FindRollcallByID(id *int64) (*types.Rollcall, error) {
	var rollcall types.Rollcall

	query := "SELECT * FROM rollcalls WHERE id = ?"
	if err := r.db.Get(&rollcall, query, id); err != nil {
		return nil, fmt.Errorf("db get: %w", err)
	}

	return &rollcall, nil
}

func (r *RollcallRepository) FindRollcallByNumber(number int64) (*types.Rollcall, error) {
	var rollcall types.Rollcall

	query := "SELECT * FROM rollcalls WHERE number = ?"
	if err := r.db.Get(&rollcall, query, number); err != nil {
		return nil, fmt.Errorf("db get: %w", err)
	}

	return &rollcall, nil
}

func (r *RollcallRepository) FindRollcalls(filter types.RollcallsFilter) ([]*types.Rollcall, error) {
	var rollcalls []*types.Rollcall

	var builder strings.Builder

	builder.WriteString("SELECT * FROM rollcalls WHERE 1 = 1")

	var args []any

	if filter.Status != nil {
		builder.WriteString(" AND status = ?")

		args = append(args, *filter.Status)
	}

	if err := r.db.Select(&rollcalls, builder.String(), args...); err != nil {
		return nil, fmt.Errorf("db select: %w", err)
	}

	if len(rollcalls) == 0 {
		return nil, sql.ErrNoRows
	}

	return rollcalls, nil
}

func (r *RollcallRepository) UpdateRollcall(id int64, update types.RollcallUpdate) (*types.Rollcall, error) {
	rollcall, err := r.FindRollcallByID(&id)
	if err != nil {
		return nil, fmt.Errorf("find rollcall by id: %w", err)
	}

	if update.Status != nil {
		applicationRepo := NewApplicationRepository(r.db)
		if err := types.CanUpdateRollcall(r, applicationRepo, id, *update.Status); err != nil {
			return nil, fmt.Errorf("can update rollcall: %w", err)
		}

		rollcall.Status = *update.Status
	}

	query := "UPDATE rollcalls SET status = :status WHERE id = :id"

	if _, err := r.db.NamedExec(query, rollcall); err != nil {
		return nil, fmt.Errorf("db select: %w", err)
	}

	return rollcall, nil
}

func (r *RollcallRepository) DeleteRollcall(id int64) error {
	if err := types.CanDeleteRollcall(r, id); err != nil {
		return fmt.Errorf("can delete rollcall: %w", err)
	}

	query := "DELETE FROM rollcalls WHERE ID = ?"

	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("db exec: %w", err)
	}

	return nil
}
