-- name: CreateTaskEvidence :one
INSERT INTO task_evidences (task_id, evidence_url, description)
VALUES ($1, $2, $3)
RETURNING *;

-- name: ListEvidencesByTask :many
SELECT id, task_id, evidence_url, description, created_at
FROM task_evidences
WHERE task_id = $1
ORDER BY created_at DESC;

-- name: UpdateTaskEvidence :one
UPDATE task_evidences
SET evidence_url = $2, description = $3
WHERE id = $1
RETURNING *;

-- name: DeleteTaskEvidence :exec
DELETE FROM task_evidences WHERE id = $1;
