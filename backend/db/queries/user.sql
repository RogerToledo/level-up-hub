-- name: CreateUser :exec
INSERT INTO users (
    username, email, password, active, current_level
) VALUES (
    $1, $2, $3, $4, $5
);

-- name: UpdateUser :exec
UPDATE users SET
    username = $2,
    email = $3,
    password = $4,
    active = $5,
    current_level = $6,
    manager_name = $7,
    manager_email = $8
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;

-- name: FindUserByEmail :one
SELECT id, username, email, password, active, role, current_level FROM users WHERE email = $1;

-- name: FindUserByID :one
SELECT id, username, email, password, active, current_level, manager_name, manager_email FROM users WHERE id = $1;

-- name: FindAllUsers :many
SELECT id, username, email, active, current_level FROM users;

-- name: FindAllUsersPaginated :many
SELECT id, username, email, active, role, current_level, created_at 
FROM users
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountAllUsers :one
SELECT COUNT(*) FROM users;