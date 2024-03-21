-- name: CreateApplication :one
INSERT INTO applications (
    enrollment_id, class_id, option, languages_score,
    humanities_score, natural_sciences_score, mathematics_score, essay_score,
    composite_score, ranking, status, selection_id, rollcall_id, applicant_id
) VALUES (
    ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
) RETURNING *;

-- name: GetApplication :one
SELECT sqlc.embed(applications), sqlc.embed(classes), sqlc.embed(quotas), sqlc.embed(periods), sqlc.embed(rollcalls), sqlc.embed(applicants)
FROM applications
JOIN classes ON applications.class_id = classes.id
JOIN periods ON classes.period_id = periods.id
JOIN quotas ON classes.quota_id = quotas.id
LEFT JOIN rollcalls ON applications.rollcall_id = rollcalls.id
JOIN applicants ON applications.applicant_id = applicants.id
WHERE applications.id = ?;

-- name: UpdateApplication :exec
UPDATE applications SET status = ?, rollcall_id = ? WHERE id = ?;

-- name: DeleteApplication :exec
DELETE FROM applications WHERE id = ?;
