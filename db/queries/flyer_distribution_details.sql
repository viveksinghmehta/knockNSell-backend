-- name: CreateDistributionDetail :one
INSERT INTO flyer_distribution_details (order_id, target_area, distribution_date, pavement_required)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: ListDistributionByOrder :many
SELECT * FROM flyer_distribution_details
WHERE order_id = $1;
