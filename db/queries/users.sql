-- name: CreateUser :one
INSERT INTO users (account_type, phone_number, email, name, photo, gender, aadhar_number, aadhar_photo_front, aadhar_photo_back, vehicle_type, age, gst_number, admin_role)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByPhoneNumber :one
SELECT *
FROM users
WHERE phone_number = $1;


-- name: UpdateUserByPhoneNumber :one
UPDATE users
SET
    account_type = COALESCE(NULLIF(sqlc.arg('account_type'), '')::account_type_enum, account_type),
    email = COALESCE(NULLIF(sqlc.arg('email'), ''), email),
    name = COALESCE(NULLIF(sqlc.arg('name'), ''), name),
    photo = COALESCE(NULLIF(sqlc.arg('photo'), ''), photo),
    gender = COALESCE(NULLIF(sqlc.arg('gender'), '')::gender_enum, gender),
    aadhar_number = COALESCE(NULLIF(sqlc.arg('aadhar_number'), ''), aadhar_number),
    aadhar_photo_front = COALESCE(NULLIF(sqlc.arg('aadhar_photo_front'), ''), aadhar_photo_front),
    aadhar_photo_back = COALESCE(NULLIF(sqlc.arg('aadhar_photo_back'), ''), aadhar_photo_back),
    vehicle_type = COALESCE(NULLIF(sqlc.arg('vehicle_type'), ''), vehicle_type),
    age = COALESCE(sqlc.arg('age'), age), -- Age is INT, no NULLIF needed unless you want special handling
    gst_number = COALESCE(NULLIF(sqlc.arg('gst_number'), ''), gst_number),
    updated_at = NOW()
WHERE phone_number = sqlc.arg('phone_number')
RETURNING *;
