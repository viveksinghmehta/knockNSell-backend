-- name: CreateFlyerOrder :one
INSERT INTO flyer_orders (user_id, service_type, total_cost)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetFlyerOrder :one
SELECT * FROM flyer_orders WHERE id = $1;

-- name: UpdateFlyerOrderStatus :exec
UPDATE flyer_orders
SET status = $2, updated_at = now()
WHERE id = $1;

-- name: ListFlyerOrdersByUser :many
SELECT * FROM flyer_orders
WHERE user_id = $1
ORDER BY order_date DESC
LIMIT $2 OFFSET $3;
