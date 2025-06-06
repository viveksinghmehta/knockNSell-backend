// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: flyer_orders.sql

package db

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const createFlyerOrder = `-- name: CreateFlyerOrder :one
INSERT INTO flyer_orders (user_id, service_type, total_cost)
VALUES ($1, $2, $3)
RETURNING id, user_id, service_type, order_date, status, total_cost, created_at, updated_at
`

type CreateFlyerOrderParams struct {
	UserID      uuid.NullUUID  `json:"user_id"`
	ServiceType string         `json:"service_type"`
	TotalCost   sql.NullString `json:"total_cost"`
}

func (q *Queries) CreateFlyerOrder(ctx context.Context, arg CreateFlyerOrderParams) (FlyerOrder, error) {
	row := q.queryRow(ctx, q.createFlyerOrderStmt, createFlyerOrder, arg.UserID, arg.ServiceType, arg.TotalCost)
	var i FlyerOrder
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.ServiceType,
		&i.OrderDate,
		&i.Status,
		&i.TotalCost,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getFlyerOrder = `-- name: GetFlyerOrder :one
SELECT id, user_id, service_type, order_date, status, total_cost, created_at, updated_at FROM flyer_orders WHERE id = $1
`

func (q *Queries) GetFlyerOrder(ctx context.Context, id uuid.UUID) (FlyerOrder, error) {
	row := q.queryRow(ctx, q.getFlyerOrderStmt, getFlyerOrder, id)
	var i FlyerOrder
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.ServiceType,
		&i.OrderDate,
		&i.Status,
		&i.TotalCost,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listFlyerOrdersByUser = `-- name: ListFlyerOrdersByUser :many
SELECT id, user_id, service_type, order_date, status, total_cost, created_at, updated_at FROM flyer_orders
WHERE user_id = $1
ORDER BY order_date DESC
LIMIT $2 OFFSET $3
`

type ListFlyerOrdersByUserParams struct {
	UserID uuid.NullUUID `json:"user_id"`
	Limit  int32         `json:"limit"`
	Offset int32         `json:"offset"`
}

func (q *Queries) ListFlyerOrdersByUser(ctx context.Context, arg ListFlyerOrdersByUserParams) ([]FlyerOrder, error) {
	rows, err := q.query(ctx, q.listFlyerOrdersByUserStmt, listFlyerOrdersByUser, arg.UserID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FlyerOrder
	for rows.Next() {
		var i FlyerOrder
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.ServiceType,
			&i.OrderDate,
			&i.Status,
			&i.TotalCost,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateFlyerOrderStatus = `-- name: UpdateFlyerOrderStatus :exec
UPDATE flyer_orders
SET status = $2, updated_at = now()
WHERE id = $1
`

type UpdateFlyerOrderStatusParams struct {
	ID     uuid.UUID      `json:"id"`
	Status sql.NullString `json:"status"`
}

func (q *Queries) UpdateFlyerOrderStatus(ctx context.Context, arg UpdateFlyerOrderStatusParams) error {
	_, err := q.exec(ctx, q.updateFlyerOrderStatusStmt, updateFlyerOrderStatus, arg.ID, arg.Status)
	return err
}
