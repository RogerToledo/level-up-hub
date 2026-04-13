-- Migration: Adicionar índices para otimização de performance
-- Data: 2026-04-12
-- Descrição: Índices baseados em análise de queries frequentes

-- ============================================================================
-- ACTIVITIES - Queries principais
-- ============================================================================

-- Índice composto para queries ordenadas por data (user_id + created_at DESC)
-- Usado em: ListUserActivities, FindUserActivities, FindDetailedActivityReport
-- Ganho esperado: 5-10x mais rápido em listagens
CREATE INDEX IF NOT EXISTS idx_activities_user_created 
ON activities(user_id, created_at DESC);

-- Índice para atividades completadas (muito usado em dashboards e reports)
-- Usado em: FindPdiDashboard, FindActivityComposition
-- Ganho esperado: 3-5x mais rápido em cálculos de XP
CREATE INDEX IF NOT EXISTS idx_activities_user_completed 
ON activities(user_id, progress_percentage) 
WHERE progress_percentage = 100;

-- Índice para ordenação por progresso
-- Usado em: FindDetailedActivityReport (ORDER BY progress_percentage DESC)
CREATE INDEX IF NOT EXISTS idx_activities_user_progress 
ON activities(user_id, progress_percentage DESC, created_at DESC);

-- ============================================================================
-- ACTIVITY_EVIDENCES - Relação com activities
-- ============================================================================

-- Índice para buscar evidências por atividade
-- Usado em: FindEvidencesByActivity, ListUserActivitiesWithEvidences
-- Ganho esperado: 10x mais rápido ao buscar evidências
CREATE INDEX IF NOT EXISTS idx_activity_evidences_activity 
ON activity_evidences(activity_id, created_at DESC);

-- ============================================================================
-- ACTIVITY_PILLARS - Joins frequentes
-- ============================================================================

-- O índice da PRIMARY KEY (activity_id, pillar) já cobre queries por activity_id
-- Adicionar índice reverso para queries específicas por pillar
CREATE INDEX IF NOT EXISTS idx_activity_pillars_pillar 
ON activity_pillars(pillar, activity_id);

-- ============================================================================
-- CAREER_LADDER - Joins e lookups
-- ============================================================================

-- Índice para queries por nível (usado em joins e filtros)
CREATE INDEX IF NOT EXISTS idx_career_ladder_level 
ON career_ladder(level);

-- ============================================================================
-- USERS - Lookups frequentes
-- ============================================================================

-- Índice para login por email (já existe UNIQUE, mas explícito para performance)
-- O UNIQUE já cria índice, mas garantimos que está otimizado
CREATE INDEX IF NOT EXISTS idx_users_email 
ON users(email) WHERE active = true;

-- Índice para queries por role (usado em listagens admin)
CREATE INDEX IF NOT EXISTS idx_users_role_active 
ON users(role, active);

-- ============================================================================
-- ANÁLISE E MONITORAMENTO
-- ============================================================================

-- Para verificar uso dos índices, execute:
-- SELECT schemaname, tablename, indexname, idx_scan, idx_tup_read, idx_tup_fetch
-- FROM pg_stat_user_indexes
-- WHERE schemaname = 'public'
-- ORDER BY idx_scan DESC;

-- Para verificar índices não utilizados:
-- SELECT schemaname, tablename, indexname, idx_scan
-- FROM pg_stat_user_indexes
-- WHERE schemaname = 'public' AND idx_scan = 0
-- ORDER BY pg_relation_size(indexrelid) DESC;
