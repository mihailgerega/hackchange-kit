---
name: using-hackkit
description: Start here in any hackkit session - the hackathon workflow (grill, OpenSpec, plan, TDD build, verify, ship), which hackkit skill or agent to use when, and the case rules that always apply.
---

<SUBAGENT-STOP>
If you were dispatched as a subagent to execute a specific task, ignore this skill and do your task.
</SUBAGENT-STOP>

# Using hackkit

hackkit is this team's Claude Code setup for the Hack&Change 2026 hackathon (Avito case: AI-buddy, an AI mentor for interns). It combines OpenSpec (spec-driven development, required by the jury), a disciplined build loop (adapted from obra/superpowers), interview-first discovery (mattpocock/skills), stack skills for Go and React, and specialist reviewers.

**Rule:** when a hackkit skill below matches what you are about to do, invoke it with the Skill tool before acting, then follow it. Project instructions (`CLAUDE.md`) and the user's direct requests win over skills. If the project has no `CLAUDE.md`/`openspec/` yet and the user wants to start the hackathon project, use `hackkit:kickoff`.

## The loop

| Phase | Do | Skill / command |
|---|---|---|
| 0. Setup (once) | scaffold repo, OpenSpec, Makefile, coverage gate, CI | `hackkit:kickoff` |
| 1. Discover | settle the scenario in rounds of questions | `hackkit:grilling` |
| 2. Specify | write the change: proposal, delta specs, design, tasks | `/opsx:propose <id>` + `hackkit:spec-authoring` |
| 3. Plan | detailed TDD plan next to the change | `hackkit:writing-plans` |
| 4. Build | fresh subagent per task + review, or inline | `hackkit:subagent-driven-development` / `hackkit:executing-plans`, always `hackkit:test-driven-development` |
| 5. Verify | evidence before "done"; then archive the change | `hackkit:verification-before-completion`, `/opsx:archive <id>` |
| 6. Ship | checklist, demo video, deck, mock jury | `hackkit:submission-check`, `hackkit:e2e-demo`, `hackkit:pitch-deck`, agent `hackkit:jury` |

Small fixes skip phases 1–3 but never skip TDD or verification. Bugs: `hackkit:systematic-debugging` first.

## Which skill for which code

- Avito facts (intern checklist, agent platform «Виталик» and MCP rules, Playbook principles and TDR, how Avito reviews code, design system): `hackkit:avito-context`
- AI agent loop, LLM calls, prompts, fake model: `hackkit:llm-agent-go`
- Mock corporate services as MCP tools, permission checks: `hackkit:go-mcp-server`
- Go: `hackkit:golang-testing`, `golang-error-handling`, `golang-database`, `golang-security`, `golang-concurrency`, `golang-context`, `golang-design-patterns`, `golang-project-layout`
- HTTP API contract: `hackkit:api-design`; Postgres: `hackkit:postgres-patterns`; Redis: `hackkit:redis-patterns`
- React: `hackkit:react-patterns`, `hackkit:react-testing`; look and feel: `hackkit:frontend-design` with Avito's design-system notes (`hackkit:avito-context`); try UI or logic variants fast: `hackkit:prototype`
- Parallel work (several people or sessions): `hackkit:parallel-sessions`, `hackkit:dispatching-parallel-agents`, `hackkit:using-git-worktrees`, `hackkit:finishing-a-development-branch`
- Review before merge: `hackkit:requesting-code-review` → agents `hackkit:go-reviewer`, `go-build-resolver`, `typescript-reviewer`, `react-reviewer`, `database-reviewer`, `security-reviewer`

## Always true in this project

1. The main scenario must run from a clean clone by the README (`docker compose up --build`, no `.env`), with `LLM_PROVIDER=fake` and no API key. If it breaks, fixing it beats every other task.
2. Synthetic data only; test instructions are labeled as test data and never presented as Avito policy.
3. Rights are enforced server-side in the tool/handler from the authenticated identity; the model never decides who may do what.
4. Tests cover > 80% of team-written lines (aim for 85%): `make cover-check`. No real LLM or network in unit tests.
5. UI: local actions < 1 s; LLM/service calls show a waiting state, time out, and fail with a clear message and retry.
6. Scoring priorities (online 40%: integrations 30, UX 25, working scenarios 20, tech 15, security 5, deck 5; final 60%: live demo 25, implementability 25, value 20, tech decisions 20): finish one end-to-end scenario before widening, and show how it plugs into Avito's platform — the buddy is a domain agent of Avito's agent platform «Виталик» (MCP Hub), plus auth/roles, AI cron, feature flags, analytics.
8. The jury are Avito engineers: ground scenarios in the official intern checklist and check code against Avito's review checklist before merging (`hackkit:avito-context`).
7. Context hygiene: keep ≤ 10 MCP servers enabled; use Context7 (`use context7`) for current docs of TanStack, Zustand, Zod, Vite, the MCP Go SDK and the Anthropic Go SDK instead of guessing APIs.
