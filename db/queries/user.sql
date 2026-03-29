-- name: CreateUser :exec
INSERT INTO users (
    username, email, password, active
) VALUES (
    $1, $2, $3, $4
);

-- name: UpdateUser :exec
UPDATE users SET
    username = $2,
    email = $3,
    password = $4,
    active = $5
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;

-- name: FindUserByEmail :one
SELECT id, username, email, password, active, role FROM users WHERE email = $1;

-- name: FindUserByID :one
SELECT id, username, email, active FROM users WHERE id = $1;

-- name: FindAllUsers :many
SELECT id, username, email, active FROM users;