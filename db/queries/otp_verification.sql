-- db/queries/otp_verification.sql

-- name: UpsertOTP :one
INSERT INTO otp_verification (phone_number, otp)
VALUES (
  @phone_number,
  ARRAY[@otp]::CHAR(6)[]
)
ON CONFLICT (phone_number) DO UPDATE
SET
  otp        = array_append(otp_verification.otp, @otp),
  updated_at = CURRENT_TIMESTAMP
WHERE array_length(otp_verification.otp, 1) < 10
RETURNING
  otp_verification.id,
  otp_verification.phone_number,
  otp_verification.otp,
  otp_verification.created_at,
  otp_verification.updated_at;


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
