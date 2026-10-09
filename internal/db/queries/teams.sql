-- name: CreateTeam :one
INSERT INTO teams (product_id, name, slug, description)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetTeamByID :one
SELECT * FROM teams WHERE id = $1;

-- name: ListTeamsByProduct :many
SELECT * FROM teams WHERE product_id = $1 ORDER BY created_at DESC;

-- name: UpdateTeam :one
UPDATE teams
SET name = $2, slug = $3, description = $4, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeleteTeam :exec
DELETE FROM teams WHERE id = $1;
