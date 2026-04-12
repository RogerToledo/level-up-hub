-- Rollback Migration: Remover índices de performance
-- Data: 2026-04-12
-- Descrição: Remove os índices criados na migration 006

-- ATENÇÃO: Execute APENAS se tiver certeza que precisa fazer rollback
-- Remover índices impacta negativamente a performance das queries

DROP INDEX IF EXISTS idx_activities_user_created;
DROP INDEX IF EXISTS idx_activities_user_completed;
DROP INDEX IF EXISTS idx_activities_user_progress;
DROP INDEX IF EXISTS idx_activity_evidences_activity;
DROP INDEX IF EXISTS idx_activity_pillars_pillar;
DROP INDEX IF EXISTS idx_career_ladder_level;
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_role_active;
