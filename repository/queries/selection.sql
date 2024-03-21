-- name: CreateSelection :one
INSERT INTO selections (
    name, kind, date, institution, course
) VALUES (
  ?, ?, ?, ?, ?
) RETURNING *;

-- name: GetSelection :one
 SELECT *
 FROM selections
 WHERE id = ?;

-- name: GetSelectionByKind :one
SELECT *
FROM selections
WHERE kind = ?;

-- name: ListSelections :many
SELECT *
FROM selections;

-- name: DeleteSelection :exec
DELETE FROM selections 
WHERE ID = ?;
