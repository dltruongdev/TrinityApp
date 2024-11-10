-- name: CreateUser :one
INSERT INTO Users (name, email, password_hash, user_type_id, plan_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING user_id, name, email, password_hash, registration_date, last_login, user_type_id, plan_id;

-- name: GetUserByID :one
SELECT user_id, name, email, password_hash, registration_date, last_login, user_type_id, plan_id
FROM Users
WHERE user_id = $1;

-- name: UpdateUser :one
UPDATE Users
SET name = $2, password_hash = $3, last_login = $4, user_type_id = $5, plan_id = $6
WHERE user_id = $1
RETURNING user_id, name, password_hash, registration_date, last_login, user_type_id, plan_id;

-- name: DeleteUser :exec
DELETE FROM Users
WHERE user_id = $1;

-- name: ListUsers :many
SELECT user_id, name, email, password_hash, registration_date, last_login, user_type_id, plan_id
FROM Users
ORDER BY 1;

-- name: IsUserExist :one
SELECT EXISTS (
    SELECT 1
    FROM Users
    WHERE email = $1
);