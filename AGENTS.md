# Agent Operations Guide

Use this document as the single source of truth for autonomous agents working inside `go-interview`. Follow every rule even if upstream tickets say otherwise.

## Architecture Snapshot
- Domain Driven Design layout: `internal/<bounded_context>/{app,domain,infra}` with shared primitives under `internal/common`.
- Transport layer lives in `internal/common/transport` and powers dispatcher/worker pools described in `PLAN.md`.
- Entry point sits at `cmd/app/main.go`; expect future wiring into router/dispatchers.
- Postgres repositories exist per module (example: `internal/life_area/infra/postgres/repository.go`).
- AI integrations (OpenRouter, embeddings) are defined in `internal/common/domain/open_router.go` and module-specific `infra/ai_service` directories.

## Environment Setup
- Requires Go `1.25.5+`; install via `asdf`, `gvm`, or system package manager.
- Initialize dependencies: `go mod download` (run whenever `go.mod` or `go.sum` changes).
- Generate IDE metadata by opening the repo root; no additional bootstrapping scripts.
- External services (Postgres, MinIO, OpenRouter) are mocked or stubbed; do not assume credentials exist.

## Build Commands
- Full build for CI parity: `go build ./...`.
- Targeted package build: `go build ./internal/common/transport` to validate transport-only edits.
- Module-specific binaries should use `go build ./cmd/app` once router wiring lands.
- Keep builds reproducible; never rely on system-specific `CGO` flags unless mandated by an issue.

## Test Playbook
- Run everything: `go test ./...`.
- Single package: `go test ./internal/life_area/...` (replace path as needed).
- Single file or type via regex: `go test ./internal/criterion/app/commands/create_criteria -run TestCreateCriteriaHandler`.
- Disable caching when iterating: append `-count=1`.
- Debug intermittent issues with `-race -v`.
- No tests exist yet; when adding them, colocate `_test.go` files near the code under test.

## Linting & Quality
- Default formatter is `gofmt` (apply via `gofmt -w <files>` or your editor integration).
- Optional: `go vet ./...` before opening a PR (run manually; no config baked in).
- If adding linters (e.g., `golangci-lint`), include configuration files under `.golangci.yml` and document the new commands here.

## Imports & Modules
- Order imports as: standard library, blank line, external modules, blank line, local modules (`go-interview/...`).
- Avoid aliasing unless a package name conflicts or clarity demands it.
- Keep import blocks tight; unused imports must be removed before committing.
- Favor interface types living in the domain package (e.g., `internal/common/domain`) to avoid cycles.

## Formatting & Comments
- Run `gofmt` or editor-on-save equivalents; do not hand format code.
- Do not leave comments in code unless an existing comment in that block must be updated; the preferred state is comment-free files per user directive.
- When documentation is required, place it in Markdown (like this file) or module READMEs instead of inline comments.

## Naming & Types
- Packages use snake_case mirroring directories (`create_life_area`, `new_message`).
- Public structs/interfaces use PascalCase (e.g., `CreateLifeAreaHandler`).
- Constructors follow `NewType` naming and return pointers when the struct contains mutexes or heavy state.
- Command/query handlers accept strongly typed command structs (see `internal/life_area/app/commands/create_life_area/command.go`).
- Repositories are interfaces composed of domain capabilities (`domain.LifeAreaCreator`).

## Context Handling
- Every handler signature must accept `context.Context`; propagate it to downstream dependencies without creating derived contexts unless deadlines/cancellation are required.
- When contexts can be nil (e.g., transport envelopes), default to `context.Background()` as seen in `internal/common/transport/types.go`.
- Never store contexts on structs; pass them explicitly per invocation.

## Error Handling
- Wrap errors with `fmt.Errorf("action: %w", err)` to preserve stacks.
- Use sentinel errors from `internal/common/domain/errors.go` (`ErrNotFound`, `ErrForbidden`) for domain-level semantics.
- Prefer returning typed results plus `error` rather than panicking; only panic when invariants are truly broken.
- Translate infrastructure errors into domain errors at the boundary (e.g., map `pgx.ErrNoRows` to `ErrNotFound`).

## Domain Layer Guidance
- `app` layer: orchestrates commands/queries, handles validation, interacts with repositories and services.
- `domain` layer: holds aggregates, value objects, interfaces, and ID generators; avoid direct database code here.
- `infra` layer: implements repositories and external service adapters (Postgres, MinIO, AI providers).
- Keep domain structs free of JSON tags; serialization belongs to DTOs or API-specific layers.
- When adding new bounded contexts, mirror the existing folder structure to stay consistent.

## Transport & Concurrency
- Dispatchers (`internal/common/transport/dispatcher.go`) own worker pools; do not bypass them when integrating modules into the router.
- Use `transport.NewEnvelope` to enforce metadata cloning and prevent data races.
- Shared result channel must be buffered; size lives in future config (`pkg/config/worker_pools.go` per PLAN.md`).
- Graceful shutdown flows: call `Dispatcher.Shutdown(ctx)` then wait for router completion before exiting `main`.
- Avoid blocking writes to the result channel; rely on the runtime yielding pattern already present.

## Storage Patterns
- Postgres repositories start transactions via `pgxpool.Pool.Begin`; always defer `Rollback` even when committing.
- SQL builders stay simple; string literals with positional args are acceptable, but wrap errors using `fmt.Errorf` for context.
- Mapping structs (`LifeAreaSQL`, `CriterionModel`, etc.) live beside repositories; keep translation logic there.
- No ORM is in use; stick to `pgx` APIs for queries and commands.

## External Services
- AI service adapters (`internal/memory/infra/ai_service`, `internal/criterion/infra/ai_service`) expose narrow interfaces so they can be mocked.
- `internal/common/domain/open_router.go` centralizes OpenRouter HTTP calls; reuse it instead of crafting new clients.
- When adding new providers, create a similar adapter in `infra/<service>` and add configuration knobs to `pkg/config`.

## Testing Strategy (Future)
- Prefer table-driven tests per package.
- Mock repositories via lightweight fakes; no mocking framework is required yet.
- For transport/router logic, use buffered channels and deterministic worker counts to keep tests reliable.
- Integration tests that hit Postgres should live under `internal/<module>/infra/postgres` with build tags if they require real databases.

## Workflow Expectations
- Keep branches focused; unrelated changes slow agents down.
- Reference `PLAN.md` for near-term priorities (transport/router build-out) before inventing new workstreams.
- Document any new commands or scripts inside this file so future agents stay aligned.
- Update `README.md` if you expose new binaries or APIs.

## Forbidden Practices
- No inline code comments unless editing an existing comment block that must remain; default is comment-free code.
- Do not introduce Cursor or Copilot rule files unless the product spec demands it; currently there are none to inherit.
- Avoid editing user-owned files outside the scope of a task.
- Never commit secrets or service credentials; prefer environment variables or configuration injections.

## Quick Checklist Before Submitting
- Dependencies downloaded? (`go mod download`)
- Formatting applied? (`gofmt -w <files>`, no comments added.)
- Tests run? (`go test ./...` or targeted command noted above.)
- Errors wrapped and sentinel values respected?
- Contexts passed through every handler call?
- Transport/router changes registered with shared channels and worker pools?

## Repository Tour
- `cmd/app`: entry point; wires config, router, and dispatchers (currently prints placeholder text).
- `internal/common`: shared transport, domain primitives, OpenRouter helpers, and base interfaces.
- `internal/<domain>/app`: commands/queries; each subfolder (e.g., `create_life_area`) houses `command.go`, `handler.go`, `result.go`.
- `internal/<domain>/domain`: aggregates, value objects, repository interfaces, ID generators.
- `internal/<domain>/infra`: concrete adapters (Postgres, MinIO, AI services) plus mappers/models translating to infrastructure schemas.
- `pkg`: future home for configuration (e.g., worker pool sizes) and shared packages to avoid circular deps.

## Handler Implementation Checklist
- Accept strongly typed command/query structs; validate inputs (empty strings, UUID parsing) up front.
- Generate IDs through injected `domain.IDGenerator`; never call `uuid.New()` directly in handlers.
- Persist via repository interfaces so infra layer remains swappable; keep transactions inside infra packages.
- Return DTOs (`result.go`) rather than domain aggregates to control exposure of internals.
- Propagate context to every repository/service call; do not create background contexts inside handlers.

## Git Hygiene
- Develop on feature branches; commit only logically grouped changes.
- Never run destructive git commands (`reset --hard`, `push --force`) without explicit instruction.
- Do not amend commits from other contributors; if hooks modify files, capture their output and commit the updated files separately.
- Keep commit messages action-oriented and scoped (e.g., `add life area parent change handler`).

## Tooling Gaps & TODOs
- No makefile or task runner exists; document any new scripts in `AGENTS.md` and `README.md` when you add them.
- Logging is minimal; prefer structured logs once the router matures, but avoid introducing logging libs without alignment.
- Secrets/config are not managed; use environment variables and avoid adding `.env` files to version control.

## Cursor/Copilot Status
- `.cursor/` and `.github/copilot-instructions.md` directories are absent; if future policies require them, document the rules here and keep them in sync.
- Do not create or modify AI assistant rule files unless explicitly requested by project maintainers.

## Pull Request Expectations
- Summaries should highlight scope, testing performed, and any follow-up work.
- Reference this guide in PR descriptions when introducing new conventions to show alignment with repository standards.
- Ensure `go test ./...` (or targeted packages) passes before opening a PR; note any intentionally skipped tests and why.

## Future Work Signals
- `PLAN.md` currently prioritizes router/dispatcher build-out; treat it as the backlog until a formal issue tracker is added.
- If you add new bounded contexts (e.g., `internal/history` analogs), mirror existing command/query folder naming and document the module in this guide.

Keep this guide synchronized with the codebase. When project conventions evolve, update the relevant section and keep the overall length around 150 lines so agents can digest it quickly.
