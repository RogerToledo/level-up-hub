# Variáveis para facilitar a manutenção
APP_NAME=level-up-hub-api
DOCKER_COMPOSE=docker-compose.yml

.PHONY: help build run test migrate-up docker-up clean

# Comando padrão ao digitar apenas 'make'
help:
	@echo "Comandos disponíveis:"
	@echo "  make build      - Compila o binário da aplicação"
	@echo "  make run        - Executa a aplicação localmente"
	@echo "  make test       - Executa todos os testes unitários"
	@echo "  make cover      - Gera relatório de cobertura de testes"
	@echo "  make docker-up  - Sobe o banco de dados e dependências (Docker)"
	@echo "  make tidy       - Limpa e organiza as dependências do Go"

build:
	go build -o bin/$(APP_NAME) main.go

run:
	go run main.go

test:
	go test ./... -v

cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

docker-up:
	docker-compose -f $(DOCKER_COMPOSE) up -d

tidy:
	go mod tidy
	go mod vendor

clean:
	rm -rf bin/
	rm coverage.out