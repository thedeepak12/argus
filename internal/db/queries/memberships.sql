-- name: CreateMembership :one
INSERT INTO memberships (user_id, org_id, product_id, team_id, role, scope)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetMembershipByID :one
SELECT * FROM memberships WHERE id = $1;

-- name: ListMembershipsByUser :many
SELECT * FROM memberships WHERE user_id = $1 AND is_active = true;

-- name: ListMembershipsByOrg :many
SELECT * FROM memberships WHERE org_id = $1 AND is_active = true;

-- name: UpdateMembership :one
UPDATE memberships
SET role = $2, scope = $3, is_active = $4, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeleteMembership :exec
DELETE FROM memberships WHERE id = $1;
