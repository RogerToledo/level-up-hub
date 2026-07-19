-- name: CreateTaskPillar :one
INSERT INTO task_pillars (task_id, pillar)
VALUES ($1, $2)
RETURNING *;

-- name: GetTaskPillars :many
SELECT pillar FROM task_pillars WHERE task_id = $1;

-- name: DeleteTaskPillars :exec
DELETE FROM task_pillars WHERE task_id = $1;
