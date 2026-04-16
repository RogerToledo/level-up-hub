-- name: CreateActivity :one
INSERT INTO activities (
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

-- name: AddEvidence :one
INSERT INTO activity_evidences (activity_id, evidence_url, description)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateActivityProgress :one
UPDATE activities 
SET 
    progress_percentage = $2,
    updated_at = NOW()
WHERE id = $1 AND user_id = $3
RETURNING *;

-- name: UpdateActivity :one
UPDATE activities 
SET 
    title = $2,
    description = $3,
    progress_percentage = $4,
    impact_summary = $5,
    is_pdi_target = $6,
    updated_at = NOW()
WHERE id = $1 AND user_id = $7
RETURNING *;

-- name: DeleteActivity :exec
DELETE FROM activities 
WHERE id = $1 AND user_id = $2;

-- name: FindActivityDetail :one
SELECT 
    a.id, a.title, a.progress_percentage, a.impact_summary,
    cl.level
FROM activities a
JOIN career_ladder cl ON a.ladder_id = cl.id
WHERE a.id = $1 AND a.user_id = $2;

-- name: FindActivityWithLadder :one
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
FROM activities a
JOIN career_ladder cl ON a.ladder_id = cl.id
WHERE a.id = $1 AND a.user_id = $2;

-- name: ListUserActivities :many
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
    (SELECT COUNT(*) FROM activity_evidences WHERE activity_id = a.id)::int as evidence_count
FROM activities a
WHERE a.user_id = $1 
ORDER BY a.created_at DESC;

-- name: ListUserActivitiesPaginated :many
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
FROM activities 
WHERE user_id = $1 
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountUserActivities :one
SELECT COUNT(*) FROM activities WHERE user_id = $1;

-- name: FindPdiDashboard :many
SELECT 
    cl.level,
    ap.pillar::text as pillar,
    SUM(CASE WHEN a.is_pdi_target = true THEN cl.xp_reward ELSE 0 END)::int as total_pdi_planned,
    SUM(CASE WHEN a.progress_percentage = 100 THEN cl.xp_reward ELSE 0 END)::int as total_achieved,
    SUM(CASE WHEN a.progress_percentage = 100 AND a.is_pdi_target = false THEN cl.xp_reward ELSE 0 END)::int as overdelivery_xp,
    COUNT(a.id)::int as activity_count
FROM activities a
JOIN activity_pillars ap ON a.id = ap.activity_id
JOIN career_ladder cl ON a.ladder_id = cl.id
WHERE a.user_id = $1
GROUP BY cl.level, ap.pillar
ORDER BY cl.level ASC;

-- name: FindUserActivities :many
SELECT 
    a.id,
    a.title,
    a.progress_percentage,
    a.is_pdi_target,
    cl.level,
    COALESCE(string_agg(ap.pillar::text, ', '), '') as pillars
FROM activities a
JOIN career_ladder cl ON a.ladder_id = cl.id
LEFT JOIN activity_pillars ap ON a.id = ap.activity_id
WHERE a.user_id = $1
GROUP BY a.id, cl.level
ORDER BY a.created_at DESC;

-- name: FindActivityByID :one
SELECT
	a.id,
	a.title,
	a.description,
	a.impact_summary,
	a.is_pdi_target,
	a.progress_percentage,
	a.ladder_id,
	a.user_id 
FROM activities a 
WHERE a.id = $1 AND a.user_id = $2;

-- name: FindEvidencesByActivity :many
SELECT * FROM activity_evidences 
WHERE activity_id = $1 
ORDER BY created_at DESC;

-- name: ListUserActivitiesWithEvidences :many
SELECT 
    a.id,
    a.title,
    a.progress_percentage,
    cl.level,
    COALESCE(
        (SELECT json_agg(ed) FROM (
            SELECT id, evidence_url, description FROM activity_evidences WHERE activity_id = a.id
        ) ed), 
        '[]'
    )::json as evidences
FROM activities a
JOIN career_ladder cl ON a.ladder_id = cl.id
WHERE a.user_id = $1
GROUP BY a.id, cl.level
ORDER BY a.created_at DESC;

-- name: ListUserActivitiesWithEvidencesPaginated :many
SELECT 
    a.id,
    a.title,
    a.progress_percentage,
    cl.level,
    COALESCE(
        (SELECT json_agg(ed) FROM (
            SELECT id, evidence_url, description FROM activity_evidences WHERE activity_id = a.id
        ) ed), 
        '[]'
    )::json as evidences
FROM activities a
JOIN career_ladder cl ON a.ladder_id = cl.id
WHERE a.user_id = $1
GROUP BY a.id, cl.level
ORDER BY a.created_at DESC
LIMIT $2 OFFSET $3;

-- name: FindDetailedActivityReport :many
SELECT 
    a.id,
    a.title,
    a.progress_percentage,
    a.is_pdi_target,
    cl.level,
    cl.xp_reward,
    (
        SELECT array_agg(ap.pillar::text)
        FROM activity_pillars ap
        WHERE ap.activity_id = a.id
    ) as pillars,
    COALESCE(
        (
            SELECT json_agg(json_build_object(
                'url', ae.evidence_url,
                'description', ae.description,
                'created_at', ae.created_at
            ))
            FROM activity_evidences ae
            WHERE ae.activity_id = a.id
        ), 
        '[]'
    )::json as evidences
FROM activities a
JOIN career_ladder cl ON a.ladder_id = cl.id
WHERE a.user_id = $1
ORDER BY a.progress_percentage DESC, a.created_at DESC;

-- name: FindGapAnalysis :many
-- Gap analysis based on PDI activities (is_pdi_target = true) vs completed activities
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
FROM activities a
JOIN career_ladder cl ON a.ladder_id = cl.id
JOIN activity_pillars ap ON a.id = ap.activity_id
WHERE a.user_id = $1 
  AND EXTRACT(YEAR FROM a.created_at)::int = $2::int
  AND a.is_pdi_target = true
GROUP BY cl.level, ap.pillar
HAVING SUM(CASE WHEN a.is_pdi_target = true THEN cl.xp_reward ELSE 0 END) > 0
ORDER BY cl.level, ap.pillar;

-- name: FindActivityComposition :many
SELECT 
    cl.level,
    COUNT(a.id)::int as total_activities,
    SUM(cl.xp_reward)::int as total_xp
FROM activities a
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