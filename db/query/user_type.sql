-- name: CreateUserType :one
INSERT INTO UserTypes (type_name)
VALUES ($1)
RETURNING user_type_id, type_name;

-- name: GetUserTypeByID :one
SELECT user_type_id, type_name
FROM UserTypes
WHERE user_type_id = $1;

-- name: UpdateUserType :one
UPDATE UserTypes
SET type_name = $2
WHERE user_type_id = $1
RETURNING user_type_id, type_name;

-- name: DeleteUserType :exec
DELETE FROM UserTypes
WHERE user_type_id = $1;

-- name: ListUserTypes :many
SELECT user_type_id, type_name
FROM UserTypes
ORDER BY 1;