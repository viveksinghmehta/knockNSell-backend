// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: flyer_distribution_details.sql

package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createDistributionDetail = `-- name: CreateDistributionDetail :one
INSERT INTO flyer_distribution_details (order_id, target_area, distribution_date, pavement_required)
VALUES ($1, $2, $3, $4)
RETURNING id, order_id, target_area, distribution_date, pavement_required, status, created_at
`

type CreateDistributionDetailParams struct {
	OrderID          uuid.NullUUID  `json:"order_id"`
	TargetArea       sql.NullString `json:"target_area"`
	DistributionDate time.Time      `json:"distribution_date"`
	PavementRequired sql.NullBool   `json:"pavement_required"`
}

func (q *Queries) CreateDistributionDetail(ctx context.Context, arg CreateDistributionDetailParams) (FlyerDistributionDetail, error) {
	row := q.queryRow(ctx, q.createDistributionDetailStmt, createDistributionDetail,
		arg.OrderID,
		arg.TargetArea,
		arg.DistributionDate,
		arg.PavementRequired,
	)
	var i FlyerDistributionDetail
	err := row.Scan(
		&i.ID,
		&i.OrderID,
		&i.TargetArea,
		&i.DistributionDate,
		&i.PavementRequired,
		&i.Status,
		&i.CreatedAt,
	)
	return i, err
}

const listDistributionByOrder = `-- name: ListDistributionByOrder :many
SELECT id, order_id, target_area, distribution_date, pavement_required, status, created_at FROM flyer_distribution_details
WHERE order_id = $1
`

func (q *Queries) ListDistributionByOrder(ctx context.Context, orderID uuid.NullUUID) ([]FlyerDistributionDetail, error) {
	rows, err := q.query(ctx, q.listDistributionByOrderStmt, listDistributionByOrder, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FlyerDistributionDetail
	for rows.Next() {
		var i FlyerDistributionDetail
		if err := rows.Scan(
			&i.ID,
			&i.OrderID,
			&i.TargetArea,
			&i.DistributionDate,
			&i.PavementRequired,
			&i.Status,
			&i.CreatedAt,
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
