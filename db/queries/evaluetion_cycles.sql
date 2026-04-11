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
    COUNT(a.id)::int as activity_count
FROM activities a
JOIN career_ladder cl ON a.ladder_id = cl.id
WHERE a.user_id = $1 
  AND a.completed_at BETWEEN $2 AND $3 
  AND a.progress_percentage = 100
GROUP BY cl.level;