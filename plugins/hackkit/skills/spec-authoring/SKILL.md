---
name: spec-authoring
description: Write OpenSpec changes for AI-buddy - agent behavior and MCP tool contracts as requirements with WHEN/THEN scenarios, mandatory permission and failure scenarios, design.md in Avito's Tech Design Review format with rights and integration map. Use with /opsx:propose, when editing specs or design.md, or reviewing a change before apply.
---

# Spec authoring (OpenSpec, AI-buddy)

The jury receives our SDD specs "по которым велась разработка AI-функций". They must read as a product contract and drive the tests. OpenSpec mechanics (`/opsx:propose`, `/opsx:apply`, `/opsx:archive`, deltas ADDED/MODIFIED/REMOVED) come from the OpenSpec skills installed by `openspec init`; project context and rules are injected from `openspec/config.yaml`. This skill is what "good" means for this case.

## One change = one intent

Name it after the user outcome: `add-timeoff-request`, `add-onboarding-plan-view`, `add-knowledge-answers`. If the proposal needs "and also", split it. A scenario that is a whole demo story usually spans 2–3 changes (knowledge answers, action tool, progress view).

## Capabilities (spec folders)

Keep behavior and contracts apart:

- `buddy-chat` — how the assistant converses: tone, citing sources, confirmations, refusals, failure messages.
- `onboarding-plan` — stages, task statuses, next steps shown to the intern.
- `<service>-tools` — one per mocked corporate service: the MCP tool contract (`timeoff-tools`, `checklist-tools`, `access-tools`).
- `knowledge-answers` — retrieval over test documents, sources, "I don't know" behavior.

## Requirement patterns

Agent behavior:

```markdown
### Requirement: Подтверждение действий
Ассистент SHALL выполнять изменяющие действия только после явного подтверждения стажёра в интерфейсе.

#### Scenario: Стажёр подтверждает заявку на отгул
- GIVEN стажёр Анна с ролью intern
- WHEN Анна пишет «оформи отгул на 20 ноября, экзамен» и нажимает «Подтвердить» в карточке предпросмотра
- THEN создаётся заявка со статусом «на согласовании» у руководителя Марии
- AND в чате показан номер заявки и следующий шаг

#### Scenario: Стажёр не подтвердил
- GIVEN карточка предпросмотра показана
- WHEN Анна нажимает «Отмена»
- THEN заявка не создаётся и ассистент предлагает изменить дату или причину
```

Tool contract (one requirement per tool):

```markdown
### Requirement: Инструмент timeoff_create_request
Сервис отсутствий SHALL предоставлять MCP-инструмент `timeoff_create_request(date: YYYY-MM-DD, reason: string)`, который создаёт заявку от имени аутентифицированного пользователя и возвращает `{id, status, approver}`; требуется право `timeoff:write`.

#### Scenario: Нет права timeoff:write
- GIVEN пользователь с правами только `timeoff:read`
- WHEN агент запрашивает список инструментов
- THEN `timeoff_create_request` отсутствует в списке
- AND прямой вызов инструмента возвращает ошибку доступа без создания заявки

#### Scenario: Сервис недоступен
- GIVEN сервис отсутствий не отвечает 5 секунд
- WHEN агент вызывает инструмент
- THEN стажёр видит «Сервис отсутствий сейчас недоступен» и кнопку «Повторить», заявка не дублируется
```

## Every requirement gets three kinds of scenario

1. **Happy path** with concrete synthetic data (names, dates, ids).
2. **Permission**: a role without the right — the data is not revealed, the action not executed (case constraint 5).
3. **Failure**: LLM timeout / service down / invalid input — what the user sees, nothing half-done (case constraint 10).

Plus, where relevant: "answer unknown → route to mentor", "prompt injection inside a document is ignored", "UI action < 1 s".

Each scenario becomes at least one test named after it (`hackkit:writing-plans` maps them). Avito's QA lead recommends exactly this: every acceptance criterion links to its test.

Ground scenarios in the official intern checklist (periods, who does what): `hackkit:avito-context` → `references/intern-journey.md`; seed: `data/seed/onboarding-checklist.yaml`.

## proposal.md must contain

- The intern's pain and who else benefits (mentor hours saved — the case cites 40–48 h per 6-month program).
- Success metric and how to check it on interns (e.g. share of questions answered without the mentor, time to first completed onboarding step, task completion via buddy, CSAT after a dialog); which analytics events feed it.
- Which corporate system is touched and that it is a mock; out of scope.

## design.md must contain

Use the section order of Avito's Tech Design Review (Playbook, `tech_design_review.md`) — Avito engineers on the jury know it: **Проблема → Описание → Схема зависимостей → Компромиссы (и техдолг) → Альтернативы → FAQ: нагрузка и ресурсы → Структура БД → Информационная безопасность → Метрики**. Inside those sections cover:

- Mermaid architecture diagram (C4 level 2): frontend → backend API → agent → MCP services; Postgres/Redis; LLM provider.
- Rights flow: where identity comes from, how scopes filter tools, where checks run.
- Integration map to Avito's platform: auth & role model, AI-agent platform with prompts and MCP («Виталик»: the buddy is a domain agent — prompt, tools, knowledge base, golden questions, owner; our tools register in its MCP Hub), AI cron (reminders about next onboarding steps), feature flags (roll out per team, A/B), analytics events (Clickstream/Kafka) — what we mock now and what changes in production. Details: `hackkit:avito-context` → `references/ai-platform.md`.
- Why an agent at all: which steps are open-ended (LLM) and which are deterministic code or cron (plan, deadlines, reminders).
- How knowledge reaches the model (whole context with cached prefix, RAG, or a compiled prompt) and why.
- LLM: provider, model, env var, timeouts, fallback (`LLM_PROVIDER=fake`), cost/latency notes.
- Risks and how the demo avoids them; "что система гарантирует, а что нет".

## Before `/opsx:apply`

Run `openspec validate <change-id>` and check: one intent; every requirement has a SHALL/MUST and the three scenario kinds; no implementation details in requirements; tasks are vertical slices ending in tests; design covers rights and integrations.
