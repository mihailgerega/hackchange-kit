# Stack commands (hackkit)

Prefer the project's `make` targets when they exist (`make test`, `make cover`, `make cover-check`); they wrap the commands below.

## Go (backend)

```bash
go test ./... -run 'TestTimeOff' -count=1          # focused RED/GREEN loop
go test ./... -race -count=1                        # full suite before commit
go test ./... -coverpkg=./... -coverprofile=coverage/go.out   # coverage profile
go tool cover -func=coverage/go.out | tail -1       # total
go tool cover -html=coverage/go.out -o coverage/go.html
```

- Table-driven tests with named subtests; integration tests behind `//go:build integration`.
- HTTP handlers: `httptest.NewRecorder()` + `httptest.NewRequest()`; servers: `httptest.NewServer()`.
- MCP tools: `mcp.NewInMemoryTransports()` (see `hackkit:go-mcp-server`).
- DB: hide `pgx`/`database/sql` behind a small repository interface so unit tests use an in-memory fake; keep one integration test per repository against real Postgres from `docker compose`, in a separate test database (Avito reviewers flag tests that share the app DB; for e2e use a separate `docker-compose.e2e.yaml`). Test business rules (status transitions, rights, edge cases), not just "repository returned no error".

## Frontend (Vite + React + TS)

```bash
npx vitest run src/features/timeoff              # focused loop
npx vitest run --coverage                        # suite + coverage (v8)
npx playwright test                              # e2e (see hackkit:e2e-demo)
```

- React Testing Library + `@testing-library/user-event`; network via MSW handlers, never real fetch.
- Zod schemas get their own unit tests (valid, invalid, edge).

## Coverage gate (case requirement: > 80% of team-written lines)

`scripts/coverage-check.sh` (installed by `hackkit:kickoff`) merges nothing and hides nothing: it reports Go and frontend separately, excludes generated code (`// Code generated ... DO NOT EDIT.`), mocks and test files, and fails below `COVERAGE_MIN` (default 85 to keep a margin over the required 80).
