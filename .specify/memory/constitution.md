# Level Up Hub Constitution

## Core Principles

### I. Modular Architecture
Each domain (`account`, `ladder`, `initiative`, `task`, `leveltarget`, `ai`) is a self-contained module within `internal/`. Modules expose Handler, Service, and DTO layers. Cross-module imports are forbidden except via interfaces in `repository/`. Each module owns its HTTP routes, business logic, and data contracts.

### II. Type-Safe Data Access
All database interactions use `sqlc`-generated code from `db/queries/`. Raw SQL is strictly forbidden in business logic. Migrations live in `db/migrations/` and are applied automatically on startup. The `repository/` package aggregates `sqlc` queries into a single `Queries` struct injected via dependency injection.

### III. Decoupled & Typed Frontend (React + TS)
The user interface is isolated in the `/frontend` directory. It must use TypeScript in strict mode (`strict: true`), forbidding explicit `any`. Components must focus solely on rendering, isolating business logic and API interactions inside Custom Hooks or service layers.

### IV. Test-First for Critical Paths
New endpoints and business logic must include unit tests (`*_test.go` in Go and `.test.tsx` in React) alongside implementation. Integration tests verify database interactions. Test coverage is tracked via `make test-coverage`. Mocks live in `internal/mocks/` and implement repository interfaces.

### V. API Contract Stability
All endpoints follow REST conventions under the `/v1/` prefix. Request and response DTOs are defined per module in `dto.go`. Swagger annotations must be maintained on handlers. Breaking changes require a version bump and a migration guide.

### VI. Configuration via Environment
All configuration flows through `config/` using `caarlos0/env`. No hardcoded values. Secrets are read via `.env` (gitignored). Database pool settings, JWT secrets, SMTP credentials, and server port must all be configurable.

## Security Requirements

- JWT authentication on all protected endpoints via `api.AuthMiddleware`
- Admin-only routes gated by `api.AdminOnly()` middleware
- CORS restricted strictly to configured origins
- Password hashing with bcrypt (`golang.org/x/crypto`)
- No secrets or credentials allowed in code or version control

## Development Workflow & SDD

- **Spec-Driven Development (Spec-Kit)**: No feature or refactoring is implemented without a spec, plan, and task list mapped inside `specs/`.
- `make check`: Runs `fmt` + `lint` + `test` (must pass before merging).
- `make sqlc`: Execute after any query changes in `db/queries/`.
- `make swagger`: Execute after endpoint documentation changes.
- `make migrate-up`: Applies new migrations locally.
- **Structured Logging**: Exclusive use of `slog` — `fmt.Println` is forbidden in production code.

## Code Style

- **Backend (Go)**: Standard Go formatting via `go fmt` and static analysis via `golangci-lint`.
  - Handler methods: `<Module>.<Action>` (e.g., `InitiativeHandler.Create`)
  - Service methods: `<Module>.<Action>` (e.g., `TaskService.ListByInitiative`)
  - DTOs: `<Module><Action>Request` / `<Module><Action>Response` structs (e.g., `CreateInitiativeRequest`)
- **Frontend (React)**: Standard formatting via ESLint and Prettier. Modular or utility-based styling.

## Governance

This constitution reflects the target architecture for the project. Amendments require:
1. Documentation of the proposed change
2. Impact analysis on existing modules and frontend
3. Updates to relevant `Makefile` targets and CI/CD workflows

**Version**: 1.1.0 | **Ratified**: 2026-09-01 | **Last Amended**: 2026-09-04