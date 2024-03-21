// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: rollcall.sql

package repository

import (
	"context"
	"database/sql"
)

const createRollcall = `-- name: CreateRollcall :one
INSERT INTO rollcalls (number, status) VALUES (?, ?) RETURNING id, number, status
`

type CreateRollcallParams struct {
	Number int64
	Status sql.NullString
}

func (q *Queries) CreateRollcall(ctx context.Context, arg CreateRollcallParams) (Rollcall, error) {
	row := q.db.QueryRowContext(ctx, createRollcall, arg.Number, arg.Status)
	var i Rollcall
	err := row.Scan(&i.ID, &i.Number, &i.Status)
	return i, err
}

const deleteRollcall = `-- name: DeleteRollcall :exec
DELETE FROM rollcalls WHERE id = ?
`

func (q *Queries) DeleteRollcall(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteRollcall, id)
	return err
}

const getRollcall = `-- name: GetRollcall :one
SELECT id, number, status FROM rollcalls WHERE id = ?
`

func (q *Queries) GetRollcall(ctx context.Context, id int64) (Rollcall, error) {
	row := q.db.QueryRowContext(ctx, getRollcall, id)
	var i Rollcall
	err := row.Scan(&i.ID, &i.Number, &i.Status)
	return i, err
}

const getRollcallByNumber = `-- name: GetRollcallByNumber :one
SELECT id, number, status FROM rollcalls WHERE number = ?
`

func (q *Queries) GetRollcallByNumber(ctx context.Context, number int64) (Rollcall, error) {
	row := q.db.QueryRowContext(ctx, getRollcallByNumber, number)
	var i Rollcall
	err := row.Scan(&i.ID, &i.Number, &i.Status)
	return i, err
}

const listRollcalls = `-- name: ListRollcalls :many
SELECT id, number, status FROM rollcalls
`

func (q *Queries) ListRollcalls(ctx context.Context) ([]Rollcall, error) {
	rows, err := q.db.QueryContext(ctx, listRollcalls)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Rollcall
	for rows.Next() {
		var i Rollcall
		if err := rows.Scan(&i.ID, &i.Number, &i.Status); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateRollcall = `-- name: UpdateRollcall :exec
UPDATE rollcalls SET status = ? WHERE id = ?
`

type UpdateRollcallParams struct {
	Status sql.NullString
	ID     int64
}

func (q *Queries) UpdateRollcall(ctx context.Context, arg UpdateRollcallParams) error {
	_, err := q.db.ExecContext(ctx, updateRollcall, arg.Status, arg.ID)
	return err
}
