-- name: CreateRollcall :one
INSERT INTO rollcalls (number, status) VALUES (?, ?) RETURNING *;

-- name: GetRollcall :one
SELECT * FROM rollcalls WHERE id = ?;

-- name: GetRollcallByNumber :one
SELECT * FROM rollcalls WHERE number = ?;

-- name: ListRollcalls :many
SELECT * FROM rollcalls;

-- name: UpdateRollcall :exec
UPDATE rollcalls SET status = ? WHERE id = ?;

-- name: DeleteRollcall :exec
DELETE FROM rollcalls WHERE id = ?;
