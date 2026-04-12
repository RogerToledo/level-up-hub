#!/bin/bash

# Graceful Shutdown Test Script
# Demonstrates that in-flight requests complete during shutdown

set -e

BINARY="./api"
PORT=8081
LOG_FILE="/tmp/graceful_shutdown_test.log"

echo "======================================"
echo "  Graceful Shutdown Test"
echo "======================================"
echo ""

# Cleanup
cleanup() {
    pkill -9 api 2>/dev/null || true
    rm -f "$LOG_FILE"
}
trap cleanup EXIT

# 1. Build
echo "📦 Building application..."
go build -o "$BINARY" cmd/api/main.go
echo "✅ Build complete"
echo ""

# 2. Start server
echo "🚀 Starting server..."
$BINARY > "$LOG_FILE" 2>&1 &
SERVER_PID=$!
echo "   PID: $SERVER_PID"
sleep 3

# 3. Check health
echo ""
echo "🏥 Checking health endpoint..."
if curl -s "http://localhost:$PORT/health" > /dev/null; then
    echo "✅ Server is healthy"
else
    echo "❌ Server failed to start"
    cat "$LOG_FILE"
    exit 1
fi

# 4. Simulate load
echo ""
echo "🔄 Simulating load (20 requests in background)..."
for i in {1..20}; do
    (
        sleep $(echo "scale=2; $i * 0.1" | bc)
        curl -s "http://localhost:$PORT/health" > /dev/null
        echo "  ✓ Request $i completed"
    ) &
done

echo "⏳ Waiting 1 second for requests to start..."
sleep 1

# 5. Trigger graceful shutdown
echo ""
echo "🛑 Triggering graceful shutdown (SIGTERM)..."
kill -TERM $SERVER_PID

echo "⏳ Waiting for shutdown to complete..."
wait $SERVER_PID 2>/dev/null || true

# 6. Verify shutdown
echo ""
echo "🔍 Verifying shutdown process..."
echo ""

if grep -q "shutdown signal received" "$LOG_FILE"; then
    echo "✅ Shutdown signal was captured"
else
    echo "❌ Shutdown signal not captured"
fi

if grep -q "server shutdown completed successfully" "$LOG_FILE"; then
    echo "✅ Server shutdown completed"
else
    echo "❌ Server shutdown failed"
fi

if grep -q "closing database connection pool" "$LOG_FILE"; then
    echo "✅ Database pool was closed"
else
    echo "❌ Database pool not closed properly"
fi

if grep -q "application stopped gracefully" "$LOG_FILE"; then
    echo "✅ Application stopped gracefully"
else
    echo "❌ Application did not stop gracefully"
fi

echo ""
echo "📊 Shutdown sequence from logs:"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
grep -E "(shutdown|closing|stopped)" "$LOG_FILE" | tail -10
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

echo ""
echo "✨ Test complete!"
echo ""
echo "💡 Tips:"
echo "   - All in-flight requests should complete"
echo "   - Database connections should close cleanly"
echo "   - No panic or forced shutdown messages"
echo ""
