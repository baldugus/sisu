-- name: CreateClass :one
INSERT INTO classes (
    period_id, quota_id, seats, minimum_score
) VALUES (
    ?, ?, ?, ?
) RETURNING *;

-- name: GetClass :one
SELECT sqlc.embed(classes), sqlc.embed(periods), sqlc.embed(quotas)
FROM classes
JOIN periods ON classes.period_id = periods.id
JOIN quotas ON classes.quota_id = quotas.id
WHERE classes.id = ?;

-- name: ListClassesByPeriodAndQuota :many
SELECT sqlc.embed(classes), sqlc.embed(periods), sqlc.embed(quotas)
FROM classes
JOIN periods ON classes.period_id = periods.id
JOIN quotas ON classes.quota_id = quotas.id
WHERE period_id = ? AND quota_id = ?;

-- name: ListClassesByPeriod :many
SELECT sqlc.embed(classes), sqlc.embed(periods), sqlc.embed(quotas)
FROM classes
JOIN periods ON classes.period_id = periods.id
JOIN quotas ON classes.quota_id = quotas.id
WHERE classes.period_id = ?;

-- name: ListClasses :many
SELECT sqlc.embed(classes), sqlc.embed(periods), sqlc.embed(quotas)
FROM classes
JOIN periods ON classes.period_id = periods.id
JOIN quotas ON classes.quota_id = quotas.id;

-- name: DeleteClass :exec
DELETE FROM classes WHERE id = ?;
