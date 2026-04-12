# Variables for easier maintenance
APP_NAME=level-up-hub-api
DOCKER_COMPOSE=docker-compose.yml
BINARY_NAME=api

.PHONY: help build run test migrate-up docker-up clean swagger dev install-tools fmt lint check sqlc stop test-shutdown

# Default command when typing just 'make'
help:
	@echo "════════════════════════════════════════════════════════"
	@echo "  Level Up Hub - Comandos Make Disponíveis"
	@echo "════════════════════════════════════════════════════════"
	@echo ""
	@echo "📦 Build & Run:"
	@echo "  make build         - Compila o binário da aplicação"
	@echo "  make run           - Executa a aplicação localmente"
	@echo "  make dev           - Roda em modo desenvolvimento (hot-reload)"
	@echo "  make stop          - Para a aplicação gracefully (SIGTERM)"
	@echo ""
	@echo "🧪 Testes:"
	@echo "  make test          - Executa todos os testes unitários"
	@echo "  make cover         - Gera relatório de cobertura HTML"
	@echo "  make test-verbose  - Testes com output detalhado"
	@echo "  make test-shutdown - Testa graceful shutdown"
	@echo ""
	@echo "📚 Documentação:"
	@echo "  make swagger       - Gera documentação Swagger/OpenAPI"
	@echo "  make docs          - Abre documentação no navegador"
	@echo ""
	@echo "🔧 Desenvolvimento:"
	@echo "  make install-tools - Instala ferramentas necessárias"
	@echo "  make fmt           - Formata código Go"
	@echo "  make lint          - Executa linter (golangci-lint)"
	@echo "  make check         - fmt + lint + test"
	@echo "  make sqlc          - Gera código sqlc"
	@echo ""
	@echo "🐳 Docker:"
	@echo "  make docker-up     - Sobe banco de dados (Docker)"
	@echo "  make docker-down   - Para containers Docker"
	@echo "  make docker-logs   - Mostra logs dos containers"
	@echo ""
	@echo "🧹 Limpeza:"
	@echo "  make clean         - Remove arquivos gerados"
	@echo "  make tidy          - Limpa e organiza dependências"
	@echo ""
	@echo "════════════════════════════════════════════════════════"

build:
	@echo "🔨 Compilando aplicação..."
	go build -o bin/$(BINARY_NAME) cmd/api/main.go
	@echo "✅ Build concluído: bin/$(BINARY_NAME)"

run:
	@echo "🚀 Iniciando aplicação..."
	go run cmd/api/main.go

dev:
	@echo "🔧 Iniciando em modo desenvolvimento..."
	ENV=dev go run cmd/api/main.go

test:
	@echo "🧪 Executando testes..."
	go test ./... -v

test-verbose:
	@echo "🧪 Executando testes (modo verbose)..."
	go test ./... -v -race -timeout 30s

cover:
	@echo "📊 Gerando relatório de cobertura..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Relatório gerado: coverage.html"

swagger:
	@echo "📚 Gerando documentação Swagger..."
	swag init -g cmd/api/main.go --output docs
	@echo "✅ Swagger docs gerados em ./docs"
	@echo "📖 Acesse: http://localhost:8081/swagger/index.html"

docs: swagger
	@echo "🌐 Abrindo documentação no navegador..."
	open http://localhost:8081/swagger/index.html || xdg-open http://localhost:8081/swagger/index.html

install-tools:
	@echo "🔧 Instalando ferramentas..."
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	@echo "✅ Ferramentas instaladas"

fmt:
	@echo "🎨 Formatando código..."
	go fmt ./...
	@echo "✅ Código formatado"

lint:
	@echo "🔍 Executando linter..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run ./... && echo "✅ Nenhum problema encontrado!"; \
	else \
		echo "⚠️  golangci-lint não instalado. Execute:"; \
		echo "    brew install golangci-lint"; \
	fi

check: fmt lint test
	@echo "✅ Todas as verificações passaram!"

sqlc:
	@echo "🔄 Gerando código sqlc..."
	sqlc generate
	@echo "✅ Código sqlc gerado"

docker-up:
	@echo "🐳 Subindo containers Docker..."
	docker-compose -f $(DOCKER_COMPOSE) up -d
	@echo "✅ Containers iniciados"

docker-down:
	@echo "🛑 Parando containers Docker..."
	docker-compose -f $(DOCKER_COMPOSE) down
	@echo "✅ Containers parados"

docker-logs:
	docker-compose -f $(DOCKER_COMPOSE) logs -f

tidy:
	@echo "🧹 Organizando dependências..."
	go mod tidy
	@echo "✅ Dependências organizadas"

clean:
	@echo "🧹 Limpando arquivos gerados..."
	rm -rf bin/
	rm -f coverage.out coverage.html
	rm -f $(BINARY_NAME)
	@echo "✅ Limpeza concluída"
	rm -rf docs/docs.go docs/swagger.json docs/swagger.yaml

stop:
	@echo "🛑 Enviando sinal de graceful shutdown..."
	@pkill -TERM $(BINARY_NAME) || echo "⚠️  Nenhum processo encontrado"
	@sleep 2
	@echo "✅ Shutdown completo"

test-shutdown:
	@echo "🧪 Testando graceful shutdown..."
	@echo "1️⃣  Iniciando servidor..."
	@(./$(BINARY_NAME) > /tmp/api.log 2>&1 &)
	@sleep 2
	@echo "2️⃣  Verificando health check..."
	@curl -s http://localhost:8081/health > /dev/null && echo "✅ Servidor respondendo" || echo "❌ Servidor não responde"
	@sleep 1
	@echo "3️⃣  Enviando SIGTERM para graceful shutdown..."
	@pkill -TERM $(BINARY_NAME)
	@sleep 2
	@echo "4️⃣  Verificando logs de shutdown..."
	@grep -q "application stopped gracefully" /tmp/api.log && echo "✅ Graceful shutdown OK" || echo "❌ Shutdown incompleto"
	@echo ""
	@echo "📝 Últimas linhas do log:"
	@tail -n 5 /tmp/api.log
	@rm -f /tmp/api.log