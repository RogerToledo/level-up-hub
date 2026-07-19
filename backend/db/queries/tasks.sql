-- name: CreateTask :one
INSERT INTO tasks (
    initiative_id, ladder_id, title, execution, impact_summary, progress_percentage
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: UpdateTask :one
UPDATE tasks
SET
    title = $2,
    ladder_id = $3,
    execution = $4,
    impact_summary = $5,
    progress_percentage = $6,
    completed_at = CASE WHEN $6 = 100 THEN CURRENT_DATE ELSE NULL END,
    updated_at = CURRENT_DATE
WHERE id = $1
RETURNING *;

-- name: DeleteTask :exec
DELETE FROM tasks WHERE id = $1;

-- name: FindTaskByID :one
SELECT id, initiative_id, ladder_id, title, execution, impact_summary, progress_percentage, completed_at, created_at, updated_at
FROM tasks
WHERE id = $1;

-- name: ListTasksByInitiative :many
SELECT id, initiative_id, ladder_id, title, execution, impact_summary, progress_percentage, completed_at, created_at, updated_at,
    (SELECT COUNT(*) FROM task_evidences WHERE task_id = tasks.id)::int as evidence_count
FROM tasks
WHERE initiative_id = $1
ORDER BY created_at DESC;

-- name: CalculateInitiativeProgress :one
SELECT COALESCE(AVG(progress_percentage), 0)::int as avg_progress
FROM tasks
WHERE initiative_id = $1;
