-- name: CreateQuota :one
INSERT INTO quotas (name) VALUES (?) RETURNING *;

-- name: GetQuota :one
SELECT * FROM quotas WHERE id = ?;

-- name: GetQuotaByName :one
SELECT * FROM quotas WHERE name = ?;
