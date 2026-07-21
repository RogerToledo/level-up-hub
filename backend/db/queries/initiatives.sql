-- name: CreateInitiative :one
INSERT INTO initiatives (
    user_id, title, description, is_pdi_target
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: UpdateInitiative :one
UPDATE initiatives
SET
    title = $2,
    description = $3,
    is_pdi_target = $4,
    updated_at = NOW()
WHERE id = $1 AND user_id = $5
RETURNING *;

-- name: DeleteInitiative :exec
DELETE FROM initiatives
WHERE id = $1 AND user_id = $2;

-- name: FindInitiativeByID :one
SELECT id, title, description, is_pdi_target, progress_percentage, user_id, created_at
FROM initiatives
WHERE id = $1 AND user_id = $2;

-- name: ListUserInitiatives :many
SELECT
    i.id,
    i.user_id,
    i.title,
    i.description,
    i.progress_percentage,
    i.is_pdi_target,
    i.completed_at,
    i.created_at,
    (SELECT COUNT(*) FROM tasks WHERE initiative_id = i.id)::int as task_count,
    (SELECT COUNT(*) > 0 FROM tasks WHERE initiative_id = i.id AND is_extra = true)::bool as has_extra
FROM initiatives i
WHERE i.user_id = $1
ORDER BY i.progress_percentage ASC, i.created_at DESC;

-- name: CountUserInitiatives :one
SELECT COUNT(*) FROM initiatives WHERE user_id = $1;

-- name: FindPdiDashboard :many
SELECT
    cl.level,
    tp.pillar::text as pillar,
    SUM(CASE WHEN i.is_pdi_target = true THEN cl.xp_reward ELSE 0 END)::int as total_pdi_planned,
    SUM(CASE WHEN t.progress_percentage = 100 THEN cl.xp_reward ELSE 0 END)::int as total_achieved,
    SUM(CASE WHEN t.progress_percentage = 100 AND i.is_pdi_target = false THEN cl.xp_reward ELSE 0 END)::int as overdelivery_xp,
    COUNT(t.id)::int as activity_count
FROM initiatives i
JOIN tasks t ON t.initiative_id = i.id
JOIN task_pillars tp ON tp.task_id = t.id
JOIN career_ladder cl ON t.ladder_id = cl.id
WHERE i.user_id = $1
GROUP BY cl.level, tp.pillar
ORDER BY cl.level ASC;

-- name: FindDetailedInitiativeReport :many
SELECT
    t.id,
    t.title,
    t.execution,
    t.impact_summary,
    t.progress_percentage,
    i.is_pdi_target,
    cl.level,
    cl.xp_reward,
    (
        SELECT array_agg(tp.pillar::text)
        FROM task_pillars tp
        WHERE tp.task_id = t.id
    ) as pillars
FROM tasks t
JOIN initiatives i ON t.initiative_id = i.id
JOIN career_ladder cl ON t.ladder_id = cl.id
WHERE i.user_id = $1
ORDER BY t.progress_percentage DESC, t.created_at DESC;

-- name: FindGapAnalysis :many
SELECT
    cl.level,
    tp.pillar::text as pillar,
    SUM(CASE WHEN i.is_pdi_target = true THEN cl.xp_reward ELSE 0 END)::int as target_xp,
    SUM(CASE WHEN t.progress_percentage = 100 THEN cl.xp_reward ELSE 0 END)::int as achieved_xp,
    (SUM(CASE WHEN i.is_pdi_target = true THEN cl.xp_reward ELSE 0 END) -
     SUM(CASE WHEN t.progress_percentage = 100 THEN cl.xp_reward ELSE 0 END))::int as gap_xp,
    CASE
        WHEN SUM(CASE WHEN i.is_pdi_target = true THEN cl.xp_reward ELSE 0 END) = 0 THEN 0
        ELSE ROUND(
            (SUM(CASE WHEN t.progress_percentage = 100 THEN cl.xp_reward ELSE 0 END)::float /
             SUM(CASE WHEN i.is_pdi_target = true THEN cl.xp_reward ELSE 0 END)::float) * 100
        )::int
    END as completion_percentage
FROM tasks t
JOIN initiatives i ON t.initiative_id = i.id
JOIN career_ladder cl ON t.ladder_id = cl.id
JOIN task_pillars tp ON tp.task_id = t.id
WHERE i.user_id = $1
  AND EXTRACT(YEAR FROM t.created_at)::int = $2::int
  AND i.is_pdi_target = true
GROUP BY cl.level, tp.pillar
HAVING SUM(CASE WHEN i.is_pdi_target = true THEN cl.xp_reward ELSE 0 END) > 0
ORDER BY cl.level, tp.pillar;

-- name: FindInitiativeComposition :many
SELECT
    cl.level,
    COUNT(t.id)::int as total_activities,
    SUM(cl.xp_reward)::int as total_xp
FROM tasks t
JOIN initiatives i ON t.initiative_id = i.id
JOIN career_ladder cl ON t.ladder_id = cl.id
WHERE i.user_id = $1 AND t.progress_percentage = 100
GROUP BY cl.level
ORDER BY cl.level ASC;

-- name: FindCurrentTargetLevel :one
SELECT
    lt.target as level,
    cl.id as ladder_id
FROM level_target lt
JOIN career_ladder cl ON cl.level = lt.target
WHERE lt.user_id = $1 AND lt.year = $2
LIMIT 1;
