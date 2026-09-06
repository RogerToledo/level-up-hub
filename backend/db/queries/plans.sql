-- name: CreatePlan :one
INSERT INTO plans (
    user_id, title, description, initiative_id, level_target, status
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: UpdatePlan :one
UPDATE plans
SET
    title = $2,
    description = $3,
    initiative_id = $4,
    level_target = $5,
    status = $6,
    updated_at = NOW()
WHERE id = $1 AND user_id = $7
RETURNING *;

-- name: DeletePlan :exec
DELETE FROM plans
WHERE id = $1 AND user_id = $2;

-- name: FindPlanByID :one
SELECT * FROM plans
WHERE id = $1 AND user_id = $2;

-- name: ListUserPlans :many
SELECT
    p.id,
    p.user_id,
    p.title,
    p.description,
    p.initiative_id,
    p.level_target,
    p.status,
    p.created_at,
    p.updated_at,
    i.title as initiative_title
FROM plans p
LEFT JOIN initiatives i ON p.initiative_id = i.id
WHERE p.user_id = $1
ORDER BY p.created_at DESC;

-- name: CountUserPlans :one
SELECT COUNT(*) FROM plans WHERE user_id = $1;

-- name: FindPlanByTitle :one
SELECT * FROM plans
WHERE user_id = $1 AND LOWER(title) = LOWER($2)
LIMIT 1;

-- name: ReorderPlan :exec
UPDATE plans
SET position = $3, updated_at = NOW()
WHERE id = $1 AND user_id = $2;

-- name: FindPlanByPosition :one
SELECT * FROM plans
WHERE user_id = $1 AND position = $2;