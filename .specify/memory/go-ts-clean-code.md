# Skill: Idiomatic Go, TypeScript & Clean Code

When writing or refactoring code, you MUST adhere to the following language standards and Clean Code principles.

## Clean Code Fundamentals
- **Single Responsibility**: Keep functions small (< 25 lines) and focused on a single task.
- **Guard Clauses**: Avoid deep `if/else` nesting. Return early when handling edge cases or errors.
- **No Magic Values**: Replace literal numbers and strings with named constants.
- **Self-Documenting**: Use descriptive variable and function names. Avoid unnecessary comments that explain *what* the code does instead of *why*.

## Idiomatic Go Standards
- **Explicit Error Handling**: Never ignore errors using `_`. Always handle or wrap them with context: `fmt.Errorf("failed to process task: %w", err)`.
- **Context Propagation**: Always pass `context.Context` as the first argument in handlers, repositories, and I/O functions.
- **No Panics**: Never use `panic()` in production logic. Return errors gracefully.
- **Structured Logging**: Use `slog` exclusively — raw `fmt.Println` or `log.Printf` are strictly forbidden.
- **Table-Driven Tests**: Write unit tests using Go's table-driven pattern (`*_test.go`).

## Strict TypeScript & React Standards
- **Strict Typing**: Explicit `any` is strictly forbidden. Use `unknown` with type guards if the type is unknown.
- **Immutability**: Prefer `const` over `let`. Avoid mutating objects or arrays directly.
- **Async Handling**: Always use `async/await` instead of promise chaining (`.then()`). Handle potential rejections explicitly.
- **React Component Isolation**: Keep components pure and presentation-focused. Move business logic, data fetching, and state management into Custom Hooks or service modules.