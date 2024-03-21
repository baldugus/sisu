-- name: CreatePeriod :one
INSERT INTO periods (name) VALUES (?) RETURNING *;

-- name: GetPeriod :one
SELECT * FROM periods WHERE id = ?;

-- name: GetPeriodByName :one
SELECT * FROM periods WHERE name = ?;

-- name: ListPeriods :many
SELECT * FROM periods;
