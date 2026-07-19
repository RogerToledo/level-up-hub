-- name: FindCurrentCycle :one
SELECT * FROM evaluation_cycles 
WHERE CURRENT_DATE BETWEEN start_date AND end_date 
AND is_active = true LIMIT 1;

-- name: FindPreviousCycle :one
SELECT * FROM evaluation_cycles 
WHERE end_date < $1 -- Busca o ciclo que terminou antes do início do atual
ORDER BY end_date DESC LIMIT 1;

-- name: FindPerformanceByPeriod :many
SELECT 
    cl.level,
    SUM(cl.xp_reward)::int as total_xp,
    COUNT(t.id)::int as activity_count
FROM tasks t
JOIN initiatives i ON t.initiative_id = i.id
JOIN career_ladder cl ON t.ladder_id = cl.id
WHERE i.user_id = $1 
  AND t.completed_at BETWEEN $2 AND $3 
  AND t.progress_percentage = 100
GROUP BY cl.level;