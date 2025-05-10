-- name: CreatePrintDetail :one
INSERT INTO flyer_print_details
    (order_id, upload_type, design_file, flyer_size, paper_quality, flyer_quantity)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: ListPrintDetailsByOrder :many
SELECT * FROM flyer_print_details
WHERE order_id = $1;
