-- name: CreateInitiative :one
INSERT INTO initiatives (
    user_id, 
    ladder_id, 
    title, 
    description, 
    progress_percentage, 
    impact_summary, 
    is_pdi_target
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: UpdateInitiativeProgress :one
UPDATE initiatives 
SET 
    progress_percentage = $2,
    updated_at = NOW()
WHERE id = $1 AND user_id = $3
RETURNING *;

-- name: UpdateInitiative :one
UPDATE initiatives 
SET 
    title = $2,
    description = $3,
    progress_percentage = $4,
    impact_summary = $5,
    is_pdi_target = $6,
    updated_at = NOW()
WHERE id = $1 AND user_id = $7
RETURNING *;

-- name: DeleteInitiative :exec
DELETE FROM initiatives 
WHERE id = $1 AND user_id = $2;

-- name: FindInitiativeByID :one
SELECT
    a.id,
    a.title,
    a.description,
    a.impact_summary,
    a.is_pdi_target,
    a.progress_percentage,
    a.ladder_id,
    a.user_id
FROM initiatives a 
WHERE a.id = $1 AND a.user_id = $2;

-- name: FindInitiativeDetail :one
SELECT 
    a.id, a.title, a.progress_percentage, a.impact_summary,
    cl.level
FROM initiatives a
JOIN career_ladder cl ON a.ladder_id = cl.id
WHERE a.id = $1 AND a.user_id = $2;

-- name: FindInitiativeWithLadder :one
SELECT 
    a.id, 
    a.user_id, 
    a.ladder_id, 
    a.title, 
    a.description, 
    a.progress_percentage, 
    a.impact_summary, 
    a.completed_at, 
    a.created_at,
    cl.level, 
    cl.xp_reward, 
    cl.technical,
    cl.expected_results,
    cl.leadership_scope
FROM initiatives a
JOIN career_ladder cl ON a.ladder_id = cl.id
WHERE a.id = $1 AND a.user_id = $2;

-- name: ListUserInitiatives :many
SELECT 
    a.id, 
    a.user_id, 
    a.ladder_id, 
    a.title, 
    a.description, 
    a.progress_percentage, 
    a.impact_summary, 
    a.is_pdi_target,
    a.completed_at, 
    a.created_at,
    (SELECT COUNT(*) FROM tasks WHERE initiative_id = a.id)::int as task_count
FROM initiatives a
WHERE a.user_id = $1 
ORDER BY a.created_at DESC;

-- name: ListUserInitiativesPaginated :many
SELECT 
    id, 
    user_id, 
    ladder_id, 
    title, 
    description, 
    progress_percentage, 
    impact_summary, 
    completed_at, 
    created_at
FROM initiatives 
WHERE user_id = $1 
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountUserInitiatives :one
SELECT COUNT(*) FROM initiatives WHERE user_id = $1;

-- name: FindPdiDashboard :many
SELECT 
    cl.level,
    ap.pillar::text as pillar,
    SUM(CASE WHEN a.is_pdi_target = true THEN cl.xp_reward ELSE 0 END)::int as total_pdi_planned,
    SUM(CASE WHEN a.progress_percentage = 100 THEN cl.xp_reward ELSE 0 END)::int as total_achieved,
    SUM(CASE WHEN a.progress_percentage = 100 AND a.is_pdi_target = false THEN cl.xp_reward ELSE 0 END)::int as overdelivery_xp,
    COUNT(a.id)::int as activity_count
FROM initiatives a
JOIN initiative_pillars ap ON a.id = ap.initiative_id
JOIN career_ladder cl ON a.ladder_id = cl.id
WHERE a.user_id = $1
GROUP BY cl.level, ap.pillar
ORDER BY cl.level ASC;

-- name: FindUserInitiatives :many
SELECT 
    a.id,
    a.title,
    a.progress_percentage,
    a.is_pdi_target,
    cl.level,
    COALESCE(string_agg(ap.pillar::text, ', '), '') as pillars
FROM initiatives a
JOIN career_ladder cl ON a.ladder_id = cl.id
LEFT JOIN initiative_pillars ap ON a.id = ap.initiative_id
WHERE a.user_id = $1
GROUP BY a.id, cl.level
ORDER BY a.created_at DESC;

-- name: FindDetailedInitiativeReport :many
SELECT 
    a.id,
    a.title,
    a.progress_percentage,
    a.is_pdi_target,
    cl.level,
    cl.xp_reward,
    (
        SELECT array_agg(ap.pillar::text)
        FROM initiative_pillars ap
        WHERE ap.initiative_id = a.id
    ) as pillars
FROM initiatives a
JOIN career_ladder cl ON a.ladder_id = cl.id
WHERE a.user_id = $1
ORDER BY a.progress_percentage DESC, a.created_at DESC;

-- name: FindGapAnalysis :many
SELECT 
    cl.level,
    ap.pillar::text as pillar,
    SUM(CASE WHEN a.is_pdi_target = true THEN cl.xp_reward ELSE 0 END)::int as target_xp,
    SUM(CASE WHEN a.progress_percentage = 100 THEN cl.xp_reward ELSE 0 END)::int as achieved_xp,
    (SUM(CASE WHEN a.is_pdi_target = true THEN cl.xp_reward ELSE 0 END) - 
     SUM(CASE WHEN a.progress_percentage = 100 THEN cl.xp_reward ELSE 0 END))::int as gap_xp,
    CASE 
        WHEN SUM(CASE WHEN a.is_pdi_target = true THEN cl.xp_reward ELSE 0 END) = 0 THEN 0
        ELSE ROUND(
            (SUM(CASE WHEN a.progress_percentage = 100 THEN cl.xp_reward ELSE 0 END)::float / 
             SUM(CASE WHEN a.is_pdi_target = true THEN cl.xp_reward ELSE 0 END)::float) * 100
        )::int
    END as completion_percentage
FROM initiatives a
JOIN career_ladder cl ON a.ladder_id = cl.id
JOIN initiative_pillars ap ON a.id = ap.initiative_id
WHERE a.user_id = $1 
  AND EXTRACT(YEAR FROM a.created_at)::int = $2::int
  AND a.is_pdi_target = true
GROUP BY cl.level, ap.pillar
HAVING SUM(CASE WHEN a.is_pdi_target = true THEN cl.xp_reward ELSE 0 END) > 0
ORDER BY cl.level, ap.pillar;

-- name: FindInitiativeComposition :many
SELECT 
    cl.level,
    COUNT(a.id)::int as total_activities,
    SUM(cl.xp_reward)::int as total_xp
FROM initiatives a
JOIN career_ladder cl ON a.ladder_id = cl.id
WHERE a.user_id = $1 AND a.progress_percentage = 100
GROUP BY cl.level
ORDER BY cl.level ASC;

-- name: FindCurrentTargetLevel :one
SELECT 
    cl.level,
    cl.id as ladder_id
FROM xp_target xt
JOIN career_ladder cl ON xt.ladder_id = cl.id
WHERE xt.year = $1
LIMIT 1;
