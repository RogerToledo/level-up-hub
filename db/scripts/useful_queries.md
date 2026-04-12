# Database Useful Queries

Scripts úteis para desenvolvimento e debugging.

## Verificar Conexões Ativas

```sql
SELECT 
    pid,
    usename,
    application_name,
    client_addr,
    state,
    query_start,
    state_change,
    query
FROM pg_stat_activity
WHERE datname = 'leveluphub_dev'
ORDER BY query_start DESC;
```

## Ver Tamanho das Tabelas

```sql
SELECT 
    schemaname,
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size,
    pg_total_relation_size(schemaname||'.'||tablename) AS size_bytes
FROM pg_tables
WHERE schemaname NOT IN ('pg_catalog', 'information_schema')
ORDER BY size_bytes DESC;
```

## Verificar Performance dos Índices

```sql
SELECT 
    schemaname,
    tablename,
    indexname,
    idx_scan AS index_scans,
    idx_tup_read AS tuples_read,
    idx_tup_fetch AS tuples_fetched
FROM pg_stat_user_indexes
ORDER BY idx_scan DESC;
```

## Índices Não Utilizados

```sql
SELECT 
    schemaname,
    tablename,
    indexname,
    idx_scan
FROM pg_stat_user_indexes
WHERE idx_scan = 0
  AND indexname NOT LIKE '%_pkey';
```

## Queries Mais Lentas

```sql
SELECT 
    calls,
    total_exec_time,
    mean_exec_time,
    max_exec_time,
    query
FROM pg_stat_statements
ORDER BY mean_exec_time DESC
LIMIT 20;
```

## Estatísticas da Tabela Users

```sql
SELECT 
    COUNT(*) AS total_users,
    COUNT(*) FILTER (WHERE active = true) AS active_users,
    COUNT(*) FILTER (WHERE is_admin = true) AS admin_users,
    COUNT(*) FILTER (WHERE created_at > NOW() - INTERVAL '30 days') AS new_users_last_30d
FROM users;
```

## Estatísticas de Atividades

```sql
SELECT 
    status,
    COUNT(*) AS count,
    ROUND(AVG(progress_percentage), 2) AS avg_progress
FROM activities
GROUP BY status
ORDER BY count DESC;
```

## Usuários Mais Ativos

```sql
SELECT 
    u.id,
    u.username,
    u.email,
    COUNT(a.id) AS total_activities,
    COUNT(a.id) FILTER (WHERE a.progress_percentage = 100) AS completed_activities,
    ROUND(AVG(a.progress_percentage), 2) AS avg_progress
FROM users u
LEFT JOIN activities a ON u.id = a.user_id
GROUP BY u.id, u.username, u.email
ORDER BY total_activities DESC
LIMIT 20;
```

## Verificar Bloqueios (Locks)

```sql
SELECT 
    l.pid,
    l.locktype,
    l.relation::regclass,
    l.mode,
    l.granted,
    a.usename,
    a.query,
    a.query_start
FROM pg_locks l
JOIN pg_stat_activity a ON l.pid = a.pid
WHERE NOT l.granted
ORDER BY a.query_start;
```

## Cache Hit Ratio (deve ser > 99%)

```sql
SELECT 
    'cache hit rate' AS metric,
    ROUND(
        sum(blks_hit) * 100.0 / NULLIF(sum(blks_hit) + sum(blks_read), 0),
        2
    ) AS percentage
FROM pg_stat_database;
```

## Limpar Dados de Teste

```sql
-- ⚠️ CUIDADO: Apenas em desenvolvimento!
-- Remove todos os dados mas mantém a estrutura

TRUNCATE TABLE activities CASCADE;
TRUNCATE TABLE users CASCADE;
TRUNCATE TABLE career_ladder CASCADE;

-- Resetar sequences
ALTER SEQUENCE IF EXISTS users_id_seq RESTART WITH 1;
```

## Criar Usuário de Teste

```sql
-- Senha: test123 (bcrypt hash exemplo)
INSERT INTO users (id, username, email, password, active, is_admin, created_at, updated_at)
VALUES (
    gen_random_uuid(),
    'testuser',
    'test@example.com',
    '$2a$10$rXKj9xXJ8L9nxN7xnJ8L9e8L9nxN7xnJ8L9nxN7xnJ8L9nxN7xnJ8',
    true,
    false,
    NOW(),
    NOW()
)
ON CONFLICT (email) DO NOTHING;
```

## Verificar Missing Indexes

```sql
SELECT 
    schemaname,
    tablename,
    seq_scan,
    seq_tup_read,
    idx_scan,
    seq_tup_read / seq_scan AS avg_seq_read
FROM pg_stat_user_tables
WHERE seq_scan > 0
ORDER BY seq_tup_read DESC
LIMIT 25;
```

## Analisar Plano de Execução

```sql
-- Exemplo: Buscar atividades de um usuário
EXPLAIN ANALYZE
SELECT *
FROM activities
WHERE user_id = 'uuid-aqui'
  AND created_at > NOW() - INTERVAL '30 days'
ORDER BY created_at DESC
LIMIT 10;
```

## Monitorar Crescimento do Banco

```sql
SELECT 
    pg_size_pretty(pg_database_size('leveluphub_dev')) AS db_size,
    (pg_database_size('leveluphub_dev') / (1024*1024*1024))::numeric(10,2) AS db_size_gb;
```

## Backup e Restore

```bash
# Backup
pg_dump -U postgres -d leveluphub_dev -F c -f backup_$(date +%Y%m%d).dump

# Restore
pg_restore -U postgres -d leveluphub_dev -c backup_20250101.dump

# Backup apenas schema
pg_dump -U postgres -d leveluphub_dev --schema-only > schema.sql

# Backup apenas dados
pg_dump -U postgres -d leveluphub_dev --data-only > data.sql
```

## Reindexar (Manutenção)

```sql
-- Reindexar uma tabela específica
REINDEX TABLE activities;

-- Reindexar toda o banco (pode demorar!)
REINDEX DATABASE leveluphub_dev;

-- Atualizar estatísticas
ANALYZE activities;

-- Vacuum + analyze
VACUUM ANALYZE activities;
```

## Monitorar Replication Lag (se houver réplicas)

```sql
SELECT 
    client_addr,
    state,
    sync_state,
    pg_size_pretty(pg_wal_lsn_diff(pg_current_wal_lsn(), sent_lsn)) AS send_lag,
    pg_size_pretty(pg_wal_lsn_diff(pg_current_wal_lsn(), write_lsn)) AS write_lag,
    pg_size_pretty(pg_wal_lsn_diff(pg_current_wal_lsn(), flush_lsn)) AS flush_lag,
    pg_size_pretty(pg_wal_lsn_diff(pg_current_wal_lsn(), replay_lsn)) AS replay_lag
FROM pg_stat_replication;
```
