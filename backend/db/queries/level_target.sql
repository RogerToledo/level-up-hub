-- name: CreateLevelTarget :one
INSERT INTO level_target (target, year)
VALUES ($1, $2)
RETURNING *;

-- name: FindLevelTargetByID :one
SELECT id, target, year
FROM level_target
WHERE id = $1;

-- name: FindLevelTargetByYear :one
SELECT id, target, year
FROM level_target
WHERE year = $1;

-- name: ListLevelTargets :many
SELECT id, target, year
FROM level_target
ORDER BY year DESC;

-- name: UpdateLevelTarget :one
UPDATE level_target
SET target = $2, year = $3
WHERE id = $1
RETURNING *;

-- name: DeleteLevelTarget :exec
DELETE FROM level_target WHERE id = $1;
