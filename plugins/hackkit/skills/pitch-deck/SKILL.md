---
name: pitch-deck
description: Build the product presentation PDF and prepare the final defense - slide structure mapped to the jury criteria, Marp markdown source exported to PDF, demo and Q&A rehearsal with likely jury questions. Use when drafting slides, exporting the deck, or preparing the live demo.
---

# Product deck and defense

No deck = 0 points at admission. The deck itself is only 5% online, but the final (60% of the total) is judged on the live demo (25%), implementability & development (25%), value (20%) and technical decisions (20%) — the deck is the spine of that talk.

## Source and export

Keep the deck as Markdown in the repo: `presentation/slides.md` (Marp), export:

```bash
npx @marp-team/marp-cli@latest presentation/slides.md --pdf --allow-local-files -o presentation/ai-buddy.pdf
```

Screenshots go in `presentation/img/` (take them with the Playwright MCP `browser_take_screenshot` on the final build). Alternative if the team prefers PowerPoint: Anthropic's `pptx` skill (`/plugin install document-skills@anthropic-agent-skills`), then export to PDF.

Front matter: `marp: true`, `paginate: true`, a clean light theme; ≤ 40 words per slide; one idea per slide; Russian. Look in the spirit of Avito's design system (black accent, generous radii, semantic colors for statuses — `hackkit:avito-context` → `references/design-system.md`) but without Avito's logo, fonts or icons.

## Slides (12, mapped to criteria)

1. **Title** — AI-buddy, team, **video link and repo link** (the case asks for the video link on the first slide).
2. **Problem** — intern's first weeks: scattered instructions, accesses, "whom to ask"; mentor spends 40–48 h per 6-month program (case data); what we take off the mentor.
3. **Who** — intern (primary), mentor, lead, HR BP: one line each on what they get.
4. **Journey & scenarios** — the official Avito intern checklist on a timeline (week 1 → equator at 3 months → graduation) with our 3–5 scenarios (or one deep one) placed on it; `hackkit:avito-context` → `references/intern-journey.md`.
5. **Demo: request → result** — 3 screenshots of the main scenario (question with source → action preview → saved result + next step).
6. **How it works** — architecture diagram: React → Go API → agent → MCP services; Postgres, Redis; LLM behind an interface. Say where the LLM is used and where deterministic code is enough (plan, deadlines, reminders), and the main trade-offs and alternatives — Avito engineers think in Tech Design Review terms.
7. **Rights and data** — identity from the token, tools filtered by role, checks in every tool, confirmation of writes, synthetic data, no PII/finance.
8. **Built into Avito's platform** (integrations, 30% online) — the buddy as a domain agent of «Виталик» (prompt, tools, knowledge, golden set, owner); our MCP servers register in its MCP Hub; auth & role model; AI cron for reminders; feature flags (Trisigma) for rollout and A/B by team; Clickstream/Kafka events; Go/PostgreSQL/Redis/React/TanStack/Zustand/Zod stack.
9. **Quality** — SDD with OpenSpec (n specs, n scenarios, each linked to a test), tests, **coverage % from `coverage/REPORT.md`**, golden-set pass rate of the agent, CI with linters, key-less demo mode.
10. **Metrics and how we'd prove it** — e.g. share of questions resolved without the mentor, time to first completed onboarding step, mentor hours saved, CSAT; events that measure them; A/B by team via feature flag.
11. **Rollout & development** (25% final) — pilot with one team → more services (business trips, benefits) → mentor digest → onboarding analytics; cost and latency of the LLM; risks.
12. **Team** — who built what (team ownership is scored), contacts.

## Defense rehearsal

- **Live demo plan**: the same path as `docs/DEMO.md`, run on the final build, with `LLM_PROVIDER=fake` as a fallback if the network or key fails; the recorded video as a second fallback. Rehearse it twice with a timer.
- **Who answers what**: assign each criterion to a teammate (value, tech decisions, implementability, security). Every member must be able to explain the architecture slide (team ownership 5%).
- **Prepare answers to likely questions:**
  - How do you make sure the bot never does something the intern isn't allowed to? (show the tool filter + handler check + test)
  - What happens when the LLM is wrong or down? (confirmation step, timeouts, fake mode, "ask your mentor")
  - How would this connect to our real absence system / auth? (MCP adapter registered in the MCP Hub of «Виталик»; token passthrough; owner and tracing per tool)
  - Why build an agent at all, why not a workflow or a form? (LLM only for open-ended questions and preparing actions; plan and reminders are code and cron)
  - How do you measure answer quality? (golden set from the checklist, pass rate, likes/dislikes, escalation rate)
  - What does the system guarantee and what not? (writes only after confirmation, no cross-user data; no guarantee of complete knowledge — it says "не знаю" and calls the mentor)
  - Why this LLM, what does it cost per dialog, latency? (model choice, effort, caching)
  - How do you know it saves mentor time? (metrics, events, A/B plan)
  - Why OpenSpec, how did specs drive the code? (show a spec → test → code trail)
- Dispatch the `hackkit:jury` agent with the deck and repo for a mock score and the top questions it would ask.
