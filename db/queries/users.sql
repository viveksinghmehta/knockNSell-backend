-- name: CreateUser :one
INSERT INTO users (account_type, phone_number, email, name, photo, gender, aadhar_number, aadhar_photo_front, aadhar_photo_back, vehicle_type, age, gst_number, admin_role)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users WHERE id = $1;

-- name: ListUsers :many
SELECT * FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2;

-- name: GetUserByPhoneNumber :one
SELECT *
FROM users
WHERE phone_number = $1;
