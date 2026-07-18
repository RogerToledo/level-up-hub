-- Migration: Add indexes for performance optimization
-- Date: 2026-04-12
-- Description: Indexes based on frequent query analysis

-- ============================================================================
-- ACTIVITIES - Main queries
-- ============================================================================

-- Composite index for date-ordered queries (user_id + created_at DESC)
-- Used in: ListUserActivities, FindUserActivities, FindDetailedActivityReport
-- Expected gain: 5-10x faster in listings
CREATE INDEX IF NOT EXISTS idx_activities_user_created 
ON activities(user_id, created_at DESC);

-- Index for completed activities (heavily used in dashboards and reports)
-- Used in: FindPdiDashboard, FindActivityComposition
-- Expected gain: 3-5x faster in XP calculations
CREATE INDEX IF NOT EXISTS idx_activities_user_completed 
ON activities(user_id, progress_percentage) 
WHERE progress_percentage = 100;

-- Index for progress ordering
-- Used in: FindDetailedActivityReport (ORDER BY progress_percentage DESC)
CREATE INDEX IF NOT EXISTS idx_activities_user_progress 
ON activities(user_id, progress_percentage DESC, created_at DESC);

-- ============================================================================
-- ACTIVITY_EVIDENCES - Relation with activities
-- ============================================================================

-- Index to search evidence by activity
-- Used in: FindEvidencesByActivity, ListUserActivitiesWithEvidences
-- Expected gain: 10x faster when fetching evidence
CREATE INDEX IF NOT EXISTS idx_activity_evidences_activity 
ON activity_evidences(activity_id, created_at DESC);

-- ============================================================================
-- ACTIVITY_PILLARS - Frequent joins
-- ============================================================================

-- PRIMARY KEY index (activity_id, pillar) already covers queries by activity_id
-- Add reverse index for specific pillar queries
CREATE INDEX IF NOT EXISTS idx_activity_pillars_pillar 
ON activity_pillars(pillar, activity_id);

-- ============================================================================
-- CAREER_LADDER - Joins and lookups
-- ============================================================================

-- Index for queries by level (used in joins and filters)
CREATE INDEX IF NOT EXISTS idx_career_ladder_level 
ON career_ladder(level);

-- ============================================================================
-- USERS - Frequent lookups
-- ============================================================================

-- Index for email login (UNIQUE already exists, but explicit for performance)
-- UNIQUE already creates index, but we ensure it's optimized
CREATE INDEX IF NOT EXISTS idx_users_email 
ON users(email) WHERE active = true;

-- Index for role queries (used in admin listings)
CREATE INDEX IF NOT EXISTS idx_users_role_active 
ON users(role, active);

-- ============================================================================
-- ANALYSIS AND MONITORING
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
