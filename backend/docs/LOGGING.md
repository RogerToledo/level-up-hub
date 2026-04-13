# Logging Estruturado - Guia de Uso

## 📌 Como Usar o Logger

### 1. Logger Automático em Handlers

Todos os requests já têm um logger com `request_id` no contexto:

```go
func (h *ActivityHandler) Create(c *gin.Context) {
    // Logger está disponível via contexto
    if log, ok := c.Get("logger"); ok {
        logger := log.(*slog.Logger)
        logger.Info("creating activity", 
            slog.String("user_id", userID.String()),
        )
    }
    
    // Ou use o logger default
    slog.Info("creating activity")
}
```

### 2. Logs Estruturados em Services

```go
func (s *Service) CreateActivity(ctx context.Context, params ...) error {
    slog.Info("creating activity",
        slog.String("user_id", params.UserID.String()),
        slog.Int("progress", int(params.Progress)),
    )
    
    // Em caso de erro
    if err != nil {
        slog.Error("failed to create activity",
            slog.String("error", err.Error()),
            slog.String("user_id", params.UserID.String()),
        )
        return err
    }
}
```

### 3. Níveis de Log

- **Debug**: Informações de desenvolvimento
- **Info**: Operações normais
- **Warn**: Avisos que não impedem funcionamento
- **Error**: Erros que precisam atenção

### 4. Formato de Saída

**Desenvolvimento (env=dev):**
```
time=2026-04-11T20:00:00.000Z level=INFO msg="incoming request" method=POST path=/activities request_id=abc-123
```

**Produção (env=prod):**
```json
{
  "time": "2026-04-11T20:00:00.000Z",
  "level": "INFO",
  "msg": "incoming request",
  "method": "POST",
  "path": "/activities",
  "request_id": "abc-123",
  "user_id": "uuid-456"
}
```

### 5. Request Tracking

Cada request tem um `request_id` único. Use-o para rastrear erros:

```bash
# Buscar todos os logs de um request específico
grep "request_id=abc-123" logs.txt

# Em JSON (produção)
jq 'select(.request_id == "abc-123")' logs.json
```

## ✅ Benefícios

- ✅ Logs estruturados e pesquisáveis
- ✅ Request tracking automático
- ✅ User tracking quando autenticado
- ✅ Performance metrics (latency)
- ✅ JSON em produção para parsing
- ✅ Texto legível em desenvolvimento
