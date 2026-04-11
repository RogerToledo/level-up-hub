-- name: CreateActivityPillar :one
INSERT INTO public.activity_pillars
(activity_id, pillar)
VALUES($1, $2)
RETURNING *;