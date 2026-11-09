---
name: avito-context
description: Avito knowledge for the AI-buddy case - the intern journey and official checklist, Avito's AI-agent platform (Vitalik, MCP Hub) and its MCP rules, engineering principles and TDR, how Avito engineers review code, the AKITA design system. Use when choosing scenarios, writing specs or design, building UI or MCP tools, or preparing README and defense.
---

# Avito context

Digest of the eight sources the case tells participants to study (checklist, Playbook, AvitoTech blog, avito.tech, tech-internship, design-system article and talk, internship site), read on 2026-11-09. The jury are Avito engineers: they judge with these principles and checklists in mind. Facts here are paraphrased with links; never present test instructions as Avito policy.

## Read the reference that matches the task

| Task | Read |
|---|---|
| Choosing scenarios, seed data, onboarding plan, metrics of value | `references/intern-journey.md` |
| Agent architecture, integration story, MCP tool design, prompt/eval strategy, agent security | `references/ai-platform.md` |
| design.md, architecture slide, defense answers, observability, data rules | `references/engineering.md` |
| Code review, README, Docker/compose, tests — "what Avito reviewers flag" | `references/review-checklist.md` |
| UI look and feel, tokens, components, UX details | `references/design-system.md` |
| Links and what each source contains | `references/sources.md` |

Seed data for the onboarding plan (both checklist sheets, with period and role): `${CLAUDE_PLUGIN_ROOT}/skills/kickoff/templates/seed/onboarding-checklist.yaml`.

## The ten things that matter most

1. **Position the buddy as a domain agent in Avito's existing agent platform ("Vitalik")**: system prompt + MCP tools + knowledge base + golden question set + owner. Our Go agent loop is a replaceable harness; mock services are MCP servers ready for Avito's MCP Hub; auth/roles, AI cron, feature flags (Trisigma), event analytics (Clickstream) are mocked, not rebuilt. This carries the 30% integrations criterion.
2. **Ground scenarios in the real checklist**: week-1 accesses and channels, the weekly plan, tool questions (avito-cli, PaaS, Sentry, Grafana, SPT, Zero Bug Policy), absence for exams (interns study), 1:1 prep, the equator meeting at 3 months.
3. **Answer like Vitalik and the DevProblems bot**: route, links, next step; cite sources; when unsure, "call the mentor" with a summary instead of guessing.
4. **MCP tools: one tool = one human task**, few tools, common names, tolerant reader (accept dates in several formats), compact outputs, actionable errors; test on a mid-tier model.
5. **Agent only where steps are unknown in advance.** Plan, deadlines and reminders are deterministic code and cron.
6. **Security**: data from documents and services is untrusted (indirect prompt injection); writes need human confirmation; rights checked in the tool; render model output without raw HTML; persist every run with `pending → success | failed`; the system prompt never reaches the client.
7. **design.md in TDR format**: problem, solution, dependency diagram, trade-offs, alternatives, load, DB structure, security, metrics; plus "what the system guarantees and what it does not".
8. **Avito reviewer checklist** (`review-checklist.md`) before every merge and before submission: one-command start with defaults, layered Go, consumer-side interfaces, business rules in code not in DB, separate migrations, graceful shutdown, `errors.Is`, linter configs in repo, loading/error/empty on every page, API layer, role-guarded routes.
9. **README the way Avito asks interns**: C4 level 2–3 and CJM as code (Mermaid), DB schema and indexes, decisions and trade-offs, LLM request example and failure handling, provider switch, and the plan/prompts the project was built from (Avito calls itself AI-first).
10. **UI in the spirit of AKITA**: semantic tokens, black accent, L/M/S controls, popup on desktop and bottom sheet on mobile, light and dark themes — without Avito's logo, fonts or icons.
