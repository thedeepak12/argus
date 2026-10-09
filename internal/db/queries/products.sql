-- name: CreateProduct :one
INSERT INTO products (org_id, name, slug, description)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetProductByID :one
SELECT * FROM products WHERE id = $1;

-- name: ListProductsByOrg :many
SELECT * FROM products WHERE org_id = $1 ORDER BY created_at DESC;

-- name: UpdateProduct :one
UPDATE products
SET name = $2, slug = $3, description = $4, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeleteProduct :exec
DELETE FROM products WHERE id = $1;
