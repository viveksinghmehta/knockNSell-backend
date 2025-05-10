-- db/queries/otp_verification.sql

-- name: CreateOTPVerification :one
INSERT INTO otp_verification (phone_number, otp)
VALUES ($1, $2)
RETURNING id, phone_number, otp, created_at, updated_at;

-- name: GetOTPByPhoneNumber :one
SELECT id, phone_number, otp, created_at, updated_at
FROM otp_verification
WHERE phone_number = $1;

-- name: UpdateOTP :exec
UPDATE otp_verification
SET otp = $2
WHERE phone_number = $1;

-- name: DeleteOTPByPhoneNumber :exec
DELETE FROM otp_verification
WHERE phone_number = $1;

-- name: ListOTPVerifications :many
SELECT id, phone_number, otp, created_at, updated_at
FROM otp_verification
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;
