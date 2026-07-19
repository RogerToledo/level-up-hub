-- name: CreateInitiativePillar :one
INSERT INTO initiative_pillars (initiative_id, pillar)
VALUES ($1, $2)
RETURNING *;
