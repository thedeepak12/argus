-- name: CreateOrganization :one
INSERT INTO organizations (name, slug, domain)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetOrganizationByID :one
SELECT * FROM organizations WHERE id = $1;

-- name: GetOrganizationBySlug :one
SELECT * FROM organizations WHERE slug = $1;

-- name: ListOrganizations :many
SELECT * FROM organizations ORDER BY created_at DESC;

-- name: UpdateOrganization :one
UPDATE organizations
SET name = $2, slug = $3, domain = $4, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeleteOrganization :exec
DELETE FROM organizations WHERE id = $1;
