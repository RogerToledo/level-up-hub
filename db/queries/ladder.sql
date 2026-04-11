-- name: CreateLadderLevel :one
INSERT INTO career_ladder (
    level, xp_reward, technical, expected_results, leadership_scope
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: FindLadderLevel :one
SELECT id, level, xp_reward, technical, expected_results, leadership_scope
FROM career_ladder
WHERE id = $1;

-- name: ListEntireLadder :many
SELECT id, level, xp_reward, technical, expected_results, leadership_scope
FROM career_ladder
ORDER BY level;

-- name: UpdateLadderXP :exec
UPDATE career_ladder
SET xp_reward = $2
WHERE id = $1;

-- name: FindLadderByLevel :one
SELECT 
    id, 
    level, 
    xp_reward,
    technical, 
    expected_results, 
    leadership_scope
FROM career_ladder 
WHERE level = $1;

