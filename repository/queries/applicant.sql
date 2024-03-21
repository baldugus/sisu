-- name: CreateApplicant :one
INSERT INTO applicants (
    cpf, name, birthdate, sex, mother_name, address_line, address_line2,
    house_number, neighborhood, municipality, state, cep, email, phone1, social_name,
    phone2
) VALUES (
  ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
) RETURNING *;

-- name: GetApplicant :one
SELECT * FROM applicants WHERE id = ?;

-- name: UpdateApplicantEmail :exec
UPDATE applicants SET email = ? WHERE id = ?;

-- name: DeleteApplicant :exec
DELETE FROM applicants WHERE id = ?;
