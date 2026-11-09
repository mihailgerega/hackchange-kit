---
name: kickoff
description: Set up the Hack&Change AI-buddy project repo - CLAUDE.md, case brief, OpenSpec with Avito context, Makefile, coverage gate, CI, Docker Compose, Go and React skeletons with the verified agent and MCP examples. Use once at the start of the hackathon or when the repo lacks this setup.
argument-hint: "[go module path, e.g. github.com/team/ai-buddy]"
---

# Kickoff: project scaffold

Say in one line what you will do: scaffold the hackathon repo so that `make check` is green within the first hour, then hand off to scenario discovery.

Templates live in `${CLAUDE_PLUGIN_ROOT}/skills/kickoff/templates/`. Verified example code lives in `${CLAUDE_PLUGIN_ROOT}/skills/llm-agent-go/reference/` and `${CLAUDE_PLUGIN_ROOT}/skills/go-mcp-server/reference/`.

**Never overwrite an existing file silently.** If a target exists, show a short diff summary and ask whether to merge, replace or skip.

## 1. Preflight (report, don't fix silently)

Check and list what is missing, with the install command for each:

| Tool | Check | Install |
|---|---|---|
| git repo | `git rev-parse --show-toplevel` | ask before `git init` |
| Go ≥ 1.25 | `go version` | https://go.dev/dl |
| Node ≥ 22 | `node --version` | https://nodejs.org |
| Docker + Compose v2 | `docker compose version` | Docker Desktop / docker-ce |
| OpenSpec CLI | `openspec --version` | `npm install -g @fission-ai/openspec@latest` |
| golangci-lint v2 (Avito reviewers expect a linter config) | `golangci-lint version` shows 2.x | `go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest` |
| gopls (Go LSP for the plugin) | `gopls version` | `go install golang.org/x/tools/gopls@latest` |
| TypeScript LSP | `typescript-language-server --version` | `npm install -g typescript typescript-language-server` |

## 2. Ask only what the repo can't tell you

Use AskUserQuestion, one small batch: Go module path (default from the argument or `github.com/<team>/ai-buddy`), team members and who owns backend/agent/frontend/specs, and whether to scaffold the backend and frontend skeletons now (recommended: yes).

## 3. Project files from templates

| Template | Target | Notes |
|---|---|---|
| `CLAUDE.md` | `CLAUDE.md` | fill `{{TEAM}}`, `{{GO_VERSION}}`; leave `{{SCENARIOS}}` until grilling |
| `README.md` | `README.md` | submission README skeleton: every section the case requires; fill as you go |
| `docs/CASE.md` | `docs/CASE.md` | the case brief; never edit requirements in it |
| `docs/STATUS.md` | `docs/STATUS.md` | fill owners |
| `Makefile` | `Makefile` | |
| `scripts/coverage-check.sh` | `scripts/coverage-check.sh` | `chmod +x` |
| `docker-compose.yml` | `docker-compose.yml` | every variable has a default: `docker compose up --build` must work without `.env`; never add `env_file: .env` |
| `env.example` | `.env.example` | documents overrides (real LLM key); `.env` is optional and goes to `.gitignore` |
| `golangci.yml` | `backend/.golangci.yml` | golangci-lint v2 config, verified on the reference code |
| `seed/onboarding-checklist.yaml` | `data/seed/onboarding-checklist.yaml` | both sheets of the official Avito intern checklist (48 items, period, role, kind); progress and people stay synthetic |
| `github-workflows/ci.yml` | `.github/workflows/ci.yml` | |
| `claude/settings.json` | `.claude/settings.json` | merge `permissions.allow` if the file exists |

Append to `.gitignore`: `.env`, `coverage/`, `.hackkit/`, `node_modules/`, `frontend/dist/`, `test-results/`, `playwright-report/`.

## 4. OpenSpec

```bash
openspec init --tools claude
```

Then replace the generated `openspec/config.yaml` with the `openspec-config.yaml` template, keeping any keys `init` wrote that the template lacks. Check with `openspec list --json` (expects a `root`).

## 5. Backend skeleton (`backend/`)

1. `go mod init <module>`; layout per `hackkit:golang-project-layout`: `cmd/buddy/main.go`, `internal/httpapi` (handlers split by domain, not one `handlers.go`), `internal/agent`, `internal/onboarding` (loads the checklist seed), `internal/services/<service>`, `internal/auth`, `migrations/` (golang-migrate, one file per change: `0001_init.up.sql`, never a single `init.sql`).
2. Copy the verified examples and fix the import path from `example.com/buddydemo` to the module:
   - `llm-agent-go/reference/*.go` → `backend/internal/agent/`
   - `go-mcp-server/reference/timeoff/*.go` → `backend/internal/services/timeoff/` (the case's own example scenario; keep it even if the team picks others — it doubles as the pattern)
   then `go get github.com/modelcontextprotocol/go-sdk@latest github.com/anthropics/anthropic-sdk-go@latest && go mod tidy`.
3. With TDD: `GET /healthz` returning `200 {"status":"ok"}`, and `main.go` wiring config from env (`HTTP_ADDR`, `LLM_PROVIDER`, `LLM_MODEL`, `LLM_TIMEOUT`), a JSON `slog` logger, and graceful shutdown (`signal.NotifyContext` + `srv.Shutdown` with a timeout) — Avito reviewers flag its absence.
4. `backend/Dockerfile`: multi-stage `golang:<go.mod version>-alpine` → `alpine:3` (keeps `wget` for the compose health check), non-root user, `EXPOSE 8080`.

## 6. Frontend skeleton (`frontend/`)

```bash
npm create vite@latest frontend -- --template react-ts
cd frontend
npm install @tanstack/react-query @tanstack/react-router zustand zod
npm install -D vitest @vitest/coverage-v8 jsdom @testing-library/react @testing-library/user-event @testing-library/jest-dom msw @playwright/test
npx playwright install chromium
```

- Merge `templates/vitest.coverage.ts` into `vite.config.ts` (`coverage.include` is required, otherwise untested files are invisible and coverage is overstated).
- `src/test/setup.ts`: import `@testing-library/jest-dom/vitest`, start/stop the MSW server.
- `src/api/client.ts`: the explicit API client, base URL from `import.meta.env.VITE_API_URL`; requests validated with Zod. Components never build URLs or headers.
- Keep the ESLint config from the Vite template, add Prettier and `eslint-plugin-import` (`import/order`, `import/no-cycle`), and an `npm run lint` script; `make lint` runs it.
- One real component test (e.g. the app shell renders the chat input), written test-first.
- `frontend/Dockerfile`: `node:24-alpine` build with `ARG VITE_API_URL` → `nginx:alpine` serving `dist/`.

## 7. Prove it

Run `make check` and `docker compose config -q` (with no `.env` present). Both must pass before you say the scaffold is done (`hackkit:verification-before-completion`). If coverage is below the gate on the skeleton, add the missing tests now — it only gets harder later.

## 8. Share the plugin with the team (ask first)

If the kit is on GitHub, register it for everyone who opens this repo:

```bash
claude plugin marketplace add <owner>/hackchange-kit --scope project
claude plugin install hackkit@hackchange --scope project
```

This writes `.claude/settings.json`; commit it.

## 9. Hand off

Offer to commit the scaffold (ask; don't commit unasked). Then continue with `hackkit:grilling` to choose the scenarios: read `docs/CASE.md` and `hackkit:avito-context` (`references/intern-journey.md` maps the official checklist to candidate scenarios and metrics) first, and propose 3–5 candidate scenarios with your recommendation (one must perform an action through a mocked corporate service).
