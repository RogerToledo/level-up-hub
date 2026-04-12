# Graceful Shutdown

## Overview

The application implements graceful shutdown to ensure all in-flight requests complete before the server stops, and all resources (database connections, goroutines, etc.) are properly cleaned up.

## How It Works

### 1. Signal Capture

The application listens for OS signals:
- **SIGINT** (Ctrl+C) - User interruption
- **SIGTERM** (kill command) - Termination request

```go
quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
<-quit
```

### 2. Shutdown Sequence

When a signal is received:

1. **Stop accepting new connections**
   - Server stops accepting new HTTP requests
   - Existing connections remain active

2. **Complete in-flight requests**
   - Wait for all active requests to complete
   - 30-second timeout for completion

3. **Close resources**
   - Database connection pool
   - Any open file handles
   - Background goroutines

4. **Exit cleanly**
   - Log completion
   - Return exit code 0

## Configuration

### Shutdown Timeout

Currently set to **30 seconds**:

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
```

To change the timeout, modify this value in `cmd/api/main.go`.

### HTTP Server Timeouts

Additional server timeout configurations:

```go
srv := &http.Server{
    ReadTimeout:       10 * time.Second,  // Time to read request
    WriteTimeout:      10 * time.Second,  // Time to write response
    IdleTimeout:       120 * time.Second, // Keep-alive timeout
    ReadHeaderTimeout: 5 * time.Second,   // Time to read headers
}
```

## Usage

### Development

```bash
# Start server
make run

# Graceful stop (Ctrl+C)
^C

# Output:
# level=INFO msg="shutdown signal received, initiating graceful shutdown..."
# level=INFO msg="server shutdown completed successfully"
# level=INFO msg="closing database connection pool..."
# level=INFO msg="database connections closed"
# level=INFO msg="application stopped gracefully"
```

### Production

```bash
# Start server
./api

# Graceful stop with SIGTERM
kill -SIGTERM $(pgrep api)

# Or with systemd
systemctl stop level-up-hub
```

## Systemd Service Example

```ini
[Unit]
Description=Level Up Hub API
After=network.target postgresql.service

[Service]
Type=notify
User=leveluphub
Group=leveluphub
WorkingDirectory=/opt/level-up-hub
ExecStart=/opt/level-up-hub/api
ExecReload=/bin/kill -HUP $MAINPID
KillMode=mixed
KillSignal=SIGTERM
TimeoutStopSec=30
Restart=on-failure

[Install]
WantedBy=multi-user.target
```

## Docker Example

```dockerfile
FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o api cmd/api/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/api .
COPY --from=builder /app/.env .

# Add signal handling
STOPSIGNAL SIGTERM

EXPOSE 8081
CMD ["./api"]
```

## Kubernetes Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: level-up-hub
spec:
  replicas: 3
  template:
    spec:
      containers:
      - name: api
        image: level-up-hub:latest
        ports:
        - containerPort: 8081
        
        # Graceful shutdown configuration
        lifecycle:
          preStop:
            exec:
              command: ["/bin/sh", "-c", "sleep 5"]
        
        # Termination grace period (must be >= shutdown timeout)
        terminationGracePeriodSeconds: 35
        
        # Liveness and readiness probes
        livenessProbe:
          httpGet:
            path: /health
            port: 8081
          initialDelaySeconds: 10
          periodSeconds: 10
        
        readinessProbe:
          httpGet:
            path: /health
            port: 8081
          initialDelaySeconds: 5
          periodSeconds: 5
```

## Testing

### Manual Test

```bash
# Terminal 1: Start server
make run

# Terminal 2: Send test requests in a loop
while true; do
  curl -s http://localhost:8081/health
  sleep 0.5
done

# Terminal 1: Stop server (Ctrl+C)
# Observe: Active requests complete before shutdown
```

### Load Test During Shutdown

```bash
# Install hey (HTTP load generator)
go install github.com/rakyll/hey@latest

# Terminal 1: Start server
make run

# Terminal 2: Start load test (10 seconds)
hey -z 10s -c 50 http://localhost:8081/health

# Terminal 1: Stop server after 5 seconds
# Wait and observe: All requests complete within timeout
```

## Monitoring Shutdown

### Logs to Watch

```bash
# During shutdown, you should see:
level=INFO msg="shutdown signal received, initiating graceful shutdown..."

# If successful:
level=INFO msg="server shutdown completed successfully"
level=INFO msg="closing database connection pool..."
level=INFO msg="database connections closed"
level=INFO msg="application stopped gracefully"

# If forced (timeout exceeded):
level=ERROR msg="server forced to shutdown" error="context deadline exceeded"
```

### Metrics

Track these metrics in production:

- **Shutdown duration** - Time from signal to complete stop
- **Active connections at shutdown** - How many requests were in-flight
- **Failed requests during shutdown** - Requests that didn't complete
- **Database pool drain time** - Time to close all DB connections

## Best Practices

### ✅ DO

- Always use graceful shutdown in production
- Set appropriate timeout based on your longest request
- Monitor shutdown metrics
- Test shutdown under load
- Use health checks for orchestration systems
- Log each shutdown phase

### ❌ DON'T

- Don't set timeout too short (requests may fail)
- Don't shutdown without draining connections
- Don't ignore shutdown errors in logs
- Don't forget to close all resources
- Don't use `os.Exit()` or `panic()` for normal shutdown

## Troubleshooting

### Shutdown Takes Too Long

**Problem:** Server takes full 30 seconds to stop

**Solution:**
1. Check for slow database queries
2. Review long-running requests
3. Ensure background tasks are cancellable
4. Consider reducing timeout if appropriate

```go
// Add request timeout middleware
r.Use(func(c *gin.Context) {
    ctx, cancel := context.WithTimeout(c.Request.Context(), 20*time.Second)
    defer cancel()
    c.Request = c.Request.WithContext(ctx)
    c.Next()
})
```

### Database Connections Not Closing

**Problem:** Database pool doesn't close properly

**Solution:**
1. Ensure `defer dbPool.Close()` is called
2. Check for leaked connections
3. Review pool stats during shutdown

```go
// Log pool stats before closing
stats := dbPool.Stat()
log.Info("draining connection pool",
    slog.Int("total_conns", int(stats.TotalConns())),
    slog.Int("idle_conns", int(stats.IdleConns())),
)
dbPool.Close()
```

### Requests Fail During Shutdown

**Problem:** 502/503 errors during deployment

**Solution:**
1. Implement health check endpoints
2. Use readiness probes (Kubernetes)
3. Implement retry logic in clients
4. Use load balancer health checks

## References

- [Go Server Shutdown Documentation](https://pkg.go.dev/net/http#Server.Shutdown)
- [Graceful Shutdown Best Practices](https://github.com/gotomicro/daemon)
- [Signal Handling in Go](https://gobyexample.com/signals)
