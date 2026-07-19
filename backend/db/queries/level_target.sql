-- name: CreateLevelTarget :one
INSERT INTO level_target (user_id, target, year)
VALUES ($1, $2, $3)
RETURNING *;

-- name: FindLevelTargetByUserAndYear :one
SELECT id, user_id, target, year
FROM level_target
WHERE user_id = $1 AND year = $2;

-- name: ListLevelTargetsByUser :many
SELECT id, user_id, target, year
FROM level_target
WHERE user_id = $1
ORDER BY year DESC;

-- name: UpsertLevelTarget :one
INSERT INTO level_target (user_id, target, year)
VALUES ($1, $2, $3)
ON CONFLICT (user_id, year) DO UPDATE SET target = $2
RETURNING *;

-- name: DeleteLevelTarget :exec
DELETE FROM level_target WHERE id = $1 AND user_id = $2;
