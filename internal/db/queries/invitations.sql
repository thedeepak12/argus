-- name: CreateInvitation :one
INSERT INTO invitations (org_id, email, role, product_id, team_id, token, expires_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetInvitationByToken :one
SELECT * FROM invitations WHERE token = $1;

-- name: ListInvitationsByOrg :many
SELECT * FROM invitations WHERE org_id = $1 AND accepted = false;

-- name: AcceptInvitation :one
UPDATE invitations
SET accepted = true
WHERE id = $1
RETURNING *;

-- name: DeleteInvitation :exec
DELETE FROM invitations WHERE id = $1;
