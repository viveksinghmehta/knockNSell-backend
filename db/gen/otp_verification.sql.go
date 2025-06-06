// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: otp_verification.sql

package db

import (
	"context"
)

const deleteOTPByPhoneNumber = `-- name: DeleteOTPByPhoneNumber :exec
DELETE FROM otp_verification
WHERE phone_number = $1
`

func (q *Queries) DeleteOTPByPhoneNumber(ctx context.Context, phoneNumber string) error {
	_, err := q.exec(ctx, q.deleteOTPByPhoneNumberStmt, deleteOTPByPhoneNumber, phoneNumber)
	return err
}

const getOTPByPhoneNumber = `-- name: GetOTPByPhoneNumber :one
SELECT id, phone_number, otp, created_at, updated_at
FROM otp_verification
WHERE phone_number = $1
`

func (q *Queries) GetOTPByPhoneNumber(ctx context.Context, phoneNumber string) (OtpVerification, error) {
	row := q.queryRow(ctx, q.getOTPByPhoneNumberStmt, getOTPByPhoneNumber, phoneNumber)
	var i OtpVerification
	err := row.Scan(
		&i.ID,
		&i.PhoneNumber,
		&i.Otp,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listOTPVerifications = `-- name: ListOTPVerifications :many
SELECT id, phone_number, otp, created_at, updated_at
FROM otp_verification
ORDER BY created_at DESC
LIMIT $1 OFFSET $2
`

type ListOTPVerificationsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListOTPVerifications(ctx context.Context, arg ListOTPVerificationsParams) ([]OtpVerification, error) {
	rows, err := q.query(ctx, q.listOTPVerificationsStmt, listOTPVerifications, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []OtpVerification
	for rows.Next() {
		var i OtpVerification
		if err := rows.Scan(
			&i.ID,
			&i.PhoneNumber,
			&i.Otp,
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

const upsertOtpVerification = `-- name: UpsertOtpVerification :one

INSERT INTO otp_verification (
    phone_number,
    otp
) VALUES (
    $1,
    ARRAY[$2::text]
)
ON CONFLICT (phone_number) DO UPDATE
SET otp = array_append(otp_verification.otp, $2::text)
RETURNING id, phone_number, otp, created_at, updated_at
`

type UpsertOtpVerificationParams struct {
	PhoneNumber string `json:"phone_number"`
	Otp         string `json:"otp"`
}

// db/queries/otp_verification.sql
func (q *Queries) UpsertOtpVerification(ctx context.Context, arg UpsertOtpVerificationParams) (OtpVerification, error) {
	row := q.queryRow(ctx, q.upsertOtpVerificationStmt, upsertOtpVerification, arg.PhoneNumber, arg.Otp)
	var i OtpVerification
	err := row.Scan(
		&i.ID,
		&i.PhoneNumber,
		&i.Otp,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
