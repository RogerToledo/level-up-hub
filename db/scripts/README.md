# Database Scripts

Scripts utilitários para análise e manutenção do banco de dados.

## 📊 analyze_indexes.sql

Analisa a performance e uso dos índices no banco de dados.

### Como executar:

```bash
# Via psql
psql -U username -d leveluphub_dev -f db/scripts/analyze_indexes.sql

# Ou
psql leveluphub_dev < db/scripts/analyze_indexes.sql
```

### O que ele mostra:

1. **Índices mais utilizados** - Top 10 por número de scans
2. **Índices não utilizados** - Candidatos para remoção
3. **Tamanho das tabelas** - Espaço em disco usado
4. **Cache hit ratio** - Eficiência do cache (meta: >99%)
5. **Estatísticas das tabelas** - Inserts, updates, deletes
6. **Queries mais lentas** - Se pg_stat_statements estiver habilitado

### Interpretação dos resultados:

**Cache Hit Ratio:**
- ✅ > 99%: Excelente
- ⚠️ 95-99%: Bom, mas pode melhorar
- ❌ < 95%: Problema! Considere aumentar shared_buffers

**Índices não utilizados:**
- Se `idx_scan = 0` há muito tempo, considere remover
- Exceção: índices em tabelas pequenas ou recém criados

**Dead rows:**
- Se `n_dead_tup` é alto, execute `VACUUM ANALYZE`
- Autovacuum deve cuidar disso automaticamente

## 🔧 Manutenção

### Reindexar todas as tabelas:

```sql
REINDEX DATABASE leveluphub_dev;
```

### Atualizar estatísticas:

```sql
ANALYZE;
```

### Limpar dead tuples:

```sql
VACUUM ANALYZE;
```

## ⚠️ Aviso

Execute scripts de análise fora de horário de pico em produção.
