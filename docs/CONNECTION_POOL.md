# Connection Pool - Guia de Configuração

## 📊 Configurações Disponíveis

### Variáveis de Ambiente

| Variável | Default | Descrição |
|----------|---------|-----------|
| `MAX_CONNS` | 25 | Número máximo de conexões simultâneas |
| `MIN_CONNS` | 5 | Número mínimo de conexões mantidas |
| `MAX_CONN_LIFETIME` | 3600 | Vida máxima de uma conexão (segundos) |
| `MAX_CONN_IDLE_TIME` | 1800 | Tempo máximo idle antes de fechar (segundos) |
| `HEALTH_CHECK_PERIOD` | 60 | Intervalo entre health checks (segundos) |
| `CONNECT_TIMEOUT` | 5 | Timeout para conectar (segundos) |

## 🎯 Recomendações por Ambiente

### Desenvolvimento
```env
MAX_CONNS=10
MIN_CONNS=2
MAX_CONN_LIFETIME=1800
MAX_CONN_IDLE_TIME=900
```

**Por quê:** Desenvolvimento não precisa de muitas conexões, economiza recursos.

---

### Staging
```env
MAX_CONNS=25
MIN_CONNS=5
MAX_CONN_LIFETIME=3600
MAX_CONN_IDLE_TIME=1800
```

**Por quê:** Simula produção mas com carga menor.

---

### Produção (Baixa Carga)
```env
MAX_CONNS=50
MIN_CONNS=10
MAX_CONN_LIFETIME=3600
MAX_CONN_IDLE_TIME=1800
HEALTH_CHECK_PERIOD=30
```

**Por quê:** Até 1000 requests/min, ~10-20 conexões simultâneas esperadas.

---

### Produção (Alta Carga)
```env
MAX_CONNS=100
MIN_CONNS=25
MAX_CONN_LIFETIME=3600
MAX_CONN_IDLE_TIME=1800
HEALTH_CHECK_PERIOD=30
```

**Por quê:** >5000 requests/min, ~50-80 conexões simultâneas esperadas.

---

## 📈 Como Calcular MAX_CONNS

### Fórmula Base
```
MAX_CONNS = (requests_por_segundo × query_duration_médio) × fator_segurança
```

### Exemplo Prático

**Cenário:**
- 100 requests/segundo
- Cada request faz 2 queries
- Tempo médio de query: 50ms (0.05s)
- Fator de segurança: 2x

**Cálculo:**
```
Queries/seg = 100 × 2 = 200
Conexões necessárias = 200 × 0.05 = 10
Com segurança = 10 × 2 = 20
```

**Resultado:** `MAX_CONNS=25` (arredondado)

---

## 🔍 Monitoramento

### Endpoint de Health Check

```bash
curl http://localhost:8081/health
```

**Resposta:**
```json
{
  "status": "up",
  "database": "ok",
  "pool_stats": {
    "total_conns": 8,
    "idle_conns": 5,
    "acquiring_conns": 0,
    "max_conns": 25
  }
}
```

### Métricas Importantes

**Uso do Pool (%):**
```
usage = (total_conns / max_conns) × 100
```

**Alertas:**
- 🟢 < 60%: Saudável
- 🟡 60-80%: Monitorar
- 🔴 > 80%: Aumentar MAX_CONNS

### Logs Automáticos

O sistema loga avisos automaticamente:

```
time=2026-04-12T10:00:00Z level=WARN msg="connection pool usage high" 
  usage_percent=85.5 total_conns=86 max_conns=100 
  recommendation="consider increasing MAX_CONNS"
```

---

## ⚠️ Problemas Comuns

### 1. "Too many connections"
**Causa:** MAX_CONNS excede limite do PostgreSQL  
**Solução:** 
```sql
-- Verificar limite do banco
SHOW max_connections;

-- Ajustar no postgresql.conf
max_connections = 200
```

### 2. Queries lentas
**Causa:** Connection acquisition bloqueada  
**Solução:** Aumentar `MAX_CONNS` ou otimizar queries

### 3. Connection leaks
**Causa:** Transações não commitadas/rolled back  
**Solução:** Sempre usar `defer tx.Rollback()`

```go
tx, err := pool.Begin(ctx)
if err != nil {
    return err
}
defer tx.Rollback(ctx) // ✅ Sempre fazer rollback no defer

// ... código ...

return tx.Commit(ctx)
```

---

## 🚀 Performance Tips

### 1. Connection Pooling
✅ **BOM:** Reutilizar conexões do pool  
❌ **RUIM:** Criar nova conexão para cada request

### 2. Prepared Statements
✅ **BOM:** sqlc já gera prepared statements automaticamente  
❌ **RUIM:** String interpolation manual

### 3. Timeouts
```go
ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
defer cancel()
```

### 4. Batch Operations
```go
// ✅ BOM: Batch insert
batch := &pgx.Batch{}
for _, item := range items {
    batch.Queue("INSERT INTO ...", item)
}
results := pool.SendBatch(ctx, batch)
defer results.Close()
```

---

## 📊 Monitoramento Avançado

### Habilitar Monitor Periódico

No `main.go`:
```go
import "github.com/me/level-up-hub/internal/database"

func main() {
    // ... setup ...
    
    // Inicia monitor a cada 5 minutos
    stopMonitor := database.StartPoolMonitor(dbPool, 5*time.Minute)
    defer func() {
        stopMonitor <- true
    }()
    
    // ... resto do código ...
}
```

**Logs gerados:**
```
time=2026-04-12T10:05:00Z level=INFO msg="connection pool stats" 
  total_conns=12 idle_conns=8 acquiring_conns=1 max_conns=25 usage_percent=48.0
```

---

## 🎓 Referências

- [PostgreSQL Connection Pooling Best Practices](https://www.postgresql.org/docs/current/runtime-config-connection.html)
- [pgx Pool Documentation](https://pkg.go.dev/github.com/jackc/pgx/v5/pgxpool)
- [Connection Pool Sizing](https://github.com/brettwooldridge/HikariCP/wiki/About-Pool-Sizing)
