-- db/queries/otp_verification.sql

-- name: UpsertOtpVerification :one
INSERT INTO otp_verification (
    phone_number,
    otp
) VALUES (
    sqlc.arg(phone_number),
    ARRAY[sqlc.arg(otp)::text]
)
ON CONFLICT (phone_number) DO UPDATE
SET otp = array_append(otp_verification.otp, sqlc.arg(otp)::text)
RETURNING *;


-- name: GetOTPByPhoneNumber :one
SELECT id, phone_number, otp, created_at, updated_at
FROM otp_verification
WHERE phone_number = $1;

-- name: DeleteOTPByPhoneNumber :exec
DELETE FROM otp_verification
WHERE phone_number = $1;

-- name: ListOTPVerifications :many
SELECT id, phone_number, otp, created_at, updated_at
FROM otp_verification
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;
