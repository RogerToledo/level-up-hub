# Level Up Hub 🚀

API para gerenciamento de carreira e desenvolvimento profissional com sistema de XP, atividades e relatórios detalhados.

## 📚 Documentação Interativa

**Swagger UI:** http://localhost:8081/swagger/index.html

Documentação completa e interativa da AP com exemplos e possibilidade de testar todos os endpoints diretamente no navegador.

## 🏗️ Tecnologias

- **Go 1.26+**
- **Gin** - Framework web
- **PostgreSQL** - Banco de dados
- **pgx/v5** - Driver PostgreSQL
- **sqlc** - Geração de código type-safe
- **JWT** - Autenticação
- **Swagger/OpenAPI** - Documentação
- **slog** - Logging estruturado

## 🚀 Quick Start

```bash
# 1. Clone o repositório
git clone https://github.com/me/level-up-hub.git
cd level-up-hub

# 2. Configure as variáveis de ambiente
cp .env.example .env

# 3. Instale dependências
go mod download

# 4. Execute migrations
psql -f db/migrations/001_create_user.sql
psql -f db/migrations/002_create_career_ladder.sql
# ...

# 5. Gere código sqlc
sqlc generate

# 6. Rode a aplicação
make run

# 7. Acesse a documentação
open http://localhost:8081/swagger/index.html
```

## 📖 Documentação

- [Swagger/OpenAPI](docs/SWAGGER.md) - Como usar e documentar endpoints
- [Logging](docs/LOGGING.md) - Sistema de logs estruturados
- [Connection Pool](docs/CONNECTION_POOL.md) - Configuração e otimização
- [Database Indexes](docs/DATABASE_INDEXES.md) - Performance do banco
- [Pagination](docs/PAGINATION.md) - Paginação de APIs

## 🛠️ Comandos Make

```bash
make help      # Lista todos os comandos disponíveis
make build     # Compila o binário
make run       # Executa a aplicação
make test      # Roda os testes
make swagger   # Gera documentação Swagger
make dev       # Roda em modo desenvolvimento
make clean     # Limpa arquivos gerados
```

## 🔐 Autenticação

A API usa **JWT (Bearer Token)**:

```bash
# 1. Registrar
POST /v1/register
{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "senha123",
  "active": true
}

# 2. Login
POST /v1/login
{
  "email": "john@example.com", 
  "password": "senha123"
}

#  3. Usar token nas requisições
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

## 📊 Features

- ✅ Autenticação JWT
- ✅ CRUD de usuários e atividades
- ✅ Sistema de XP e níveis
- ✅ Dashboard de progresso
- ✅ Relatórios detalhados
- ✅ Análise de gap
- ✅ Evidências de atividades
- ✅ Paginação
- ✅ Logging estruturado
- ✅ Health check com métricas
- ✅ Documentação Swagger
- ✅ Connection pooling configurável

## 🗂️ Estrutura do Projeto

```
.
├── cmd/api/               # Ponto de entrada da aplicação
├── config/                # Configurações
├── db/
│   ├── migrations/        # Migrations SQL
│   ├── queries/           # Queries SQLC
│   └── scripts/           # Scripts utilitários
├── docs/                  # Documentação (Swagger auto-gerado)
├── internal/
│   ├── account/           # Módulo de usuários
│   ├── activity/          # Módulo de atividades
│   ├── api/               # Middlewares
│   ├── database/          # Configuração do banco
│   ├── ladder/            # Career ladder
│   ├── logger/            # Logging estruturado
│   ├── pagination/        # Sistema de paginação
│   ├── pkg/identity/      # Utilitários de identidade
│   ├── repository/        # Código gerado pelo SQLC
│   └── rest/              # Respostas padronizadas
└── routes/                # Definição de rotas

```

## ⚙️ Configuração

Variáveis de ambiente disponíveis em `.env`:

```bash
# Application
ENV=dev
PORT=8081

# Database
DB_URL_DEV=postgres://...
MAX_CONNS=25
MIN_CONNS=5
MAX_CONN_LIFETIME=3600
MAX_CONN_IDLE_TIME=1800

# Security
JWT_SECRET=your-secret-key
```

Ver [.env.example](.env.example) para todas as opções.

## 🧪 Testes

```bash
# Rodar todos os testes
make test

# Testes com cobertura
make cover

# Testes de um pacote específico
go test ./internal/account/... -v
```

## 📈 Performance

- **Connection Pooling**: Otimizado para alta carga
- **Índices no BD**: Queries 5-10x mais rápidas
- **Paginação**: Suporta milhões de registros
- **Logging estruturado**: Debugging eficiente

## 🤝 Contribuindo

1. Fork o projeto
2. Crie uma branch (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add some AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request

## 📝 Licença

Este projeto está sob a licença MIT. Veja o arquivo [LICENSE](LICENSE) para mais detalhes.

## 👥 Autores

- **Seu Nome** - *Trabalho Inicial*

## 🙏 Agradecimentos

- Comunidade Go
- Mantenedores do SQLC
- Time do Gin Framework