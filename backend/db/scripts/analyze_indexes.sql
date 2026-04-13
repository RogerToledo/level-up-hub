-- Script de Análise de Performance de Índices
-- Execute este script no psql ou DBeaver para verificar o estado dos índices

\echo '========================================='
\echo 'ANÁLISE DE PERFORMANCE - ÍNDICES'
\echo '========================================='
\echo ''

-- 1. ÍNDICES MAIS UTILIZADOS
\echo '1. ÍNDICES MAIS UTILIZADOS (Top 10)'
\echo '-----------------------------------'
SELECT 
    schemaname,
    tablename,
    indexname,
    idx_scan as scans,
    pg_size_pretty(pg_relation_size(indexrelid)) as size
FROM pg_stat_user_indexes
WHERE schemaname = 'public'
ORDER BY idx_scan DESC
LIMIT 10;

\echo ''
\echo '2. ÍNDICES NÃO UTILIZADOS'
\echo '-------------------------'
SELECT 
    schemaname,
    tablename,
    indexname,
    pg_size_pretty(pg_relation_size(indexrelid)) as size
FROM pg_stat_user_indexes
WHERE schemaname = 'public' 
  AND idx_scan = 0
  AND indexrelid NOT IN (
    SELECT indexrelid FROM pg_index WHERE indisprimary OR indisunique
  )
ORDER BY pg_relation_size(indexrelid) DESC;

\echo ''
\echo '3. TAMANHO DAS TABELAS E ÍNDICES'
\echo '--------------------------------'
SELECT 
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) as total_size,
    pg_size_pretty(pg_relation_size(schemaname||'.'||tablename)) as table_size,
    pg_size_pretty(pg_indexes_size(schemaname||'.'||tablename)) as indexes_size
FROM pg_tables
WHERE schemaname = 'public'
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;

\echo ''
\echo '4. CACHE HIT RATIO (Meta: > 99%)'
\echo '--------------------------------'
SELECT 
    sum(heap_blks_read) as heap_read,
    sum(heap_blks_hit) as heap_hit,
    ROUND(
        sum(heap_blks_hit)::numeric / 
        NULLIF(sum(heap_blks_hit) + sum(heap_blks_read), 0) * 100, 
        2
    )::text || '%' as cache_hit_ratio
FROM pg_statio_user_tables;

\echo ''
\echo '5. ESTATÍSTICAS DAS TABELAS PRINCIPAIS'
\echo '--------------------------------------'
SELECT 
    schemaname,
    tablename,
    n_tup_ins as inserts,
    n_tup_upd as updates,
    n_tup_del as deletes,
    n_live_tup as live_rows,
    n_dead_tup as dead_rows,
    last_vacuum,
    last_autovacuum
FROM pg_stat_user_tables
WHERE schemaname = 'public'
ORDER BY n_live_tup DESC;

\echo ''
\echo '6. QUERIES MAIS LENTAS (requer pg_stat_statements)'
\echo '-------------------------------------------------'
-- Descomente se pg_stat_statements estiver habilitado
/*
SELECT 
    LEFT(query, 80) as query_preview,
    calls,
    ROUND(total_exec_time::numeric, 2) as total_time_ms,
    ROUND(mean_exec_time::numeric, 2) as avg_time_ms,
    ROUND((stddev_exec_time)::numeric, 2) as stddev_ms
FROM pg_stat_statements
WHERE query NOT LIKE '%pg_stat%'
ORDER BY mean_exec_time DESC
LIMIT 10;
*/

\echo ''
\echo '========================================='
\echo 'FIM DA ANÁLISE'
\echo '========================================='
