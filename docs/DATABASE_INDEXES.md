# Database Indexes - Guia de Performance

## 📊 Índices Implementados

### 1. **Activities - Listagens por Usuário**
```sql
CREATE INDEX idx_activities_user_created 
ON activities(user_id, created_at DESC);
```
**Queries otimizadas:**
- Lista de atividades por data
- Timeline de atividades
- Paginação ordenada por data

**Ganho:** 5-10x mais rápido

---

### 2. **Activities - Atividades Completadas**
```sql
CREATE INDEX idx_activities_user_completed 
ON activities(user_id, progress_percentage) 
WHERE progress_percentage = 100;
```
**Queries otimizadas:**
- Cálculo de XP total
- Dashboard de progresso
- Relatórios de completude

**Ganho:** 3-5x mais rápido  
**Bônus:** Índice parcial economiza espaço (só indexa completed)

---

### 3. **Activities - Ordenação por Progresso**
```sql
CREATE INDEX idx_activities_user_progress 
ON activities(user_id, progress_percentage DESC, created_at DESC);
```
**Queries otimizadas:**
- Relatórios detalhados
- Atividades ordenadas por % conclusão

**Ganho:** 3-4x mais rápido

---

### 4. **Activity Evidences - Lookup por Atividade**
```sql
CREATE INDEX idx_activity_evidences_activity 
ON activity_evidences(activity_id, created_at DESC);
```
**Queries otimizadas:**
- Buscar evidências de uma atividade
- Lista de evidências ordenada

**Ganho:** 10x mais rápido

---

### 5. **Activity Pillars - Query por Pilar**
```sql
CREATE INDEX idx_activity_pillars_pillar 
ON activity_pillars(pillar, activity_id);
```
**Queries otimizadas:**
- Filtros por pilar
- Análise de gap por categoria

---

### 6. **Career Ladder - Query por Nível**
```sql
CREATE INDEX idx_career_ladder_level 
ON career_ladder(level);
```
**Queries otimizadas:**
- Joins com activities
- Filtros por nível

---

### 7. **Users - Login e Role**
```sql
CREATE INDEX idx_users_email ON users(email) WHERE active = true;
CREATE INDEX idx_users_role_active ON users(role, active);
```
**Queries otimizadas:**
- Login por email
- Listagem de usuários admin
- Filtros por role

---

## 🔍 Monitoramento de Performance

### Verificar Uso dos Índices

```sql
SELECT 
    schemaname,
    tablename,
    indexname,
    idx_scan as "Scans",
    idx_tup_read as "Tuples Read",
    idx_tup_fetch as "Tuples Fetched",
    pg_size_pretty(pg_relation_size(indexrelid)) as "Size"
FROM pg_stat_user_indexes
WHERE schemaname = 'public'
ORDER BY idx_scan DESC;
```

**Interpretação:**
- `idx_scan > 0`: Índice está sendo usado ✅
- `idx_scan = 0`: Índice não está sendo usado ⚠️
- `idx_scan > 10000`: Índice muito utilizado 🔥

---

### Identificar Índices Não Utilizados

```sql
SELECT 
    schemaname,
    tablename,
    indexname,
    idx_scan,
    pg_size_pretty(pg_relation_size(indexrelid)) as "Size"
FROM pg_stat_user_indexes
WHERE schemaname = 'public' 
  AND idx_scan = 0
  AND indexrelid NOT IN (
    SELECT indexrelid FROM pg_index WHERE indisprimary OR indisunique
  )
ORDER BY pg_relation_size(indexrelid) DESC;
```

**Ação:** Considere remover índices não utilizados que ocupam espaço.

---

### Análise de Query Performance

```sql
-- Habilitar tracking de queries (se ainda não estiver)
-- postgresql.conf: shared_preload_libraries = 'pg_stat_statements'

SELECT 
    query,
    calls,
    total_exec_time,
    mean_exec_time,
    stddev_exec_time,
    rows
FROM pg_stat_statements
WHERE query LIKE '%activities%'
ORDER BY mean_exec_time DESC
LIMIT 10;
```

---

### Cache Hit Ratio

```sql
SELECT 
    sum(heap_blks_read) as heap_read,
    sum(heap_blks_hit) as heap_hit,
    sum(heap_blks_hit) / (sum(heap_blks_hit) + sum(heap_blks_read)) as cache_ratio
FROM pg_statio_user_tables;
```

**Meta:** >0.99 (99% de cache hit) é excelente  
**< 0.95:** Considere aumentar `shared_buffers` no PostgreSQL

---

## 📈 Impacto Esperado

### Antes dos Índices
```
Query: SELECT * FROM activities WHERE user_id = $1 ORDER BY created_at DESC
Tempo: 150ms (tabela com 10k registros)
Explicação: Sequential scan em toda a tabela
```

### Depois dos Índices
```
Query: SELECT * FROM activities WHERE user_id = $1 ORDER BY created_at DESC
Tempo: 15ms (tabela com 10k registros)
Explicação: Index scan usando idx_activities_user_created
```

**Ganho: 10x mais rápido** 🚀

---

## ⚖️ Trade-offs

### Benefícios
✅ Queries 5-10x mais rápidas  
✅ Melhor experiência do usuário  
✅ Menor uso de CPU do banco  
✅ Aplicação escala melhor

### Custos
⚠️ Espaço em disco (cada índice ~5-20% do tamanho da tabela)  
⚠️ INSERT/UPDATE/DELETE levemente mais lentos  
⚠️ Índices precisam de manutenção (VACUUM/REINDEX)

---

## 🛠️ Manutenção

### Reindexar (se necessário)

```sql
-- Reindexar uma tabela
REINDEX TABLE activities;

-- Reindexar um índice específico
REINDEX INDEX idx_activities_user_created;

-- Reindexar todo o banco
REINDEX DATABASE leveluphub_prod;
```

**Quando fazer:**
- Após inserção massiva de dados
- Se queries ficarem lentas com o tempo
- Durante janela de manutenção

---

### Analisar e Atualizar Estatísticas

```sql
-- Atualizar estatísticas de uma tabela
ANALYZE activities;

-- Atualizar estatísticas de todo o banco
ANALYZE;
```

**Frequência recomendada:**
- Automático via autovacuum (default)
- Manual após grandes mudanças de dados

---

## 🎯 Boas Práticas

### 1. Menos é Mais
❌ Evite criar índices para todas as colunas  
✅ Crie apenas para queries realmente lentas

### 2. Índices Compostos
✅ `(user_id, created_at)` é melhor que dois separados  
❌ Ordem das colunas importa!

### 3. WHERE Clauses
✅ Use índices parciais quando possível: `WHERE active = true`  
✅ Economiza espaço e acelera queries específicas

### 4. Monitoramento
✅ Revise uso de índices mensalmente  
✅ Remova índices não utilizados  
✅ Adicione novos conforme aplicação cresce

---

## 📚 Referências

- [PostgreSQL Index Types](https://www.postgresql.org/docs/current/indexes-types.html)
- [Index Maintenance](https://www.postgresql.org/docs/current/routine-reindex.html)
- [Performance Tuning](https://wiki.postgresql.org/wiki/Performance_Optimization)
