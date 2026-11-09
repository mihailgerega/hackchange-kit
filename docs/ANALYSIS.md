# Разбор источников: что взяли и почему

Все репозитории склонированы и прочитаны локально 9.11.2026. Звёзды взяты из GitHub API в тот же день. Ниже отмечено, что вошло в `hackkit` и что осталось за бортом. Критерий отбора: польза для **этого** кейса (Go + React, обязательный SDD, MCP, покрытие > 80%, видео, презентация, 2–3 дня работы) при минимальной нагрузке на контекст.

## Главное ограничение, которое определило отбор

У Claude Code есть бюджет на список описаний скиллов: около 1% контекстного окна, то есть ~30k символов при окне 1M и ~6k при 200k. Описания, которые не влезли, выпадают из контекста, и модель видит только имя скилла. В опубликованном эксперименте с 160 скиллами сохранились 20 описаний, а из трёх скиллов, оставшихся без описания, модель не нашла ни одного.

| Набор | Скиллов | Символов в описаниях |
|---|---|---|
| ECC целиком | 314 | ~104 000 (+15 000 у команд) |
| Superpowers | 15 | ~2 700 |
| mattpocock/skills | 39 | ~6 200 |
| **hackkit** | **35** | **~7 500**; `claude plugin details`: ~2.5k токенов постоянно + ~1.1k от хука |

Поэтому кит собран из отобранных частей, а не из целых репозиториев.

---

## 1. obra/superpowers: 315k★, v7.0.2, MIT

**Что это.** Методология разработки, оформленная как скиллы: brainstorming → writing-plans → subagent-driven-development (свежий субагент на каждую задачу, ревью после каждой) → TDD → verification → review. Хук SessionStart подмешивает `using-superpowers`.

**Что смотрели.** Все 15 скиллов, промпты исполнителя и ревьюеров, скрипты `sdd-workspace`, `task-brief`, `review-package`, хук.

**Взяли (с адаптацией):**
- `test-driven-development`: «железный закон» TDD. Добавили соответствие «сценарий спеки → тест», команды Go/Vitest и порог покрытия.
- `writing-plans`: на вход теперь подаётся OpenSpec-change, план сохраняется в `openspec/changes/<id>/plan.md`, задачи связаны с `tasks.md`.
- `subagent-driven-development` и `executing-plans` со скриптами. Изменён только путь рабочей папки: `.superpowers/sdd` → `.hackkit/sdd`.
- `verification-before-completion`: добавлена таблица «что значит готово» для кейса: `make check`, e2e, чистый клон без ключа.
- `requesting-code-review`: добавлена маршрутизация на специализированных ревьюеров.
- `systematic-debugging`, `dispatching-parallel-agents`, `using-git-worktrees`, `finishing-a-development-branch` — без изменений.
- Механизм хука SessionStart: переписан под `using-hackkit`, только на builtins bash.

**Не взяли.** `brainstorming`: он отличный, но длинный и разговорный, и для хакатона с обязательным OpenSpec его заменяет `grilling`. Также `writing-skills`, `diagnosing-superpowers` и `receiving-code-review` (последний закрывается fix-loop в SDD).

## 2. affaan-m/ECC (бывший everything-claude-code): 300k★, v2.2.3, MIT

**Что это.** Сетап победителя хакатона Anthropic x Forum Ventures (сентябрь 2025, проект zenith.chat), выросший до продукта для 15 агентов-ассистентов: 71 агент, 314 скиллов, 98 команд, 27 хуков.

**Что смотрели.** Гайды автора, манифесты профилей установки, хуки, агентов и скиллы под наш стек.

**Взяли:**
- Агентов `go-reviewer`, `go-build-resolver`, `typescript-reviewer`, `react-reviewer`, `security-reviewer`, `database-reviewer`. Убран шаблонный блок «Prompt Defense Baseline», ссылки переведены на скиллы hackkit. В `security-reviewer` добавлен чек-лист для LLM-агента и MCP: авторизация в тулзе, фильтр тулзов по роли, prompt injection, подтверждение записи.
- Скиллы `react-patterns` и `react-testing` (RTL + Vitest + MSW, TanStack Query, Zustand), `api-design`, `postgres-patterns`, `redis-patterns`.
- Из гайда автора: не больше 10 MCP и 80 тулзов, worktrees для параллельной работы, LSP-плагины.

**Не взяли:**
- Полную установку: она переполняет бюджет описаний.
- `--with lang:go`: эта опция ставит 72 скилла на все языки.
- Хуки: среди них есть хук, который блокирует `npm run dev` вне tmux, 8 хуков на Stop и хуки на каждый вызов инструмента.
- `tdd-workflow`: 590 строк, отдельный коммит на каждую стадию RED/GREEN, зависимость от скриптов ECC.
- `orch-build-mvp`: зависит от GAN-цикла до 15 итераций.
- `mcp-server-patterns`: только Node/TypeScript.
- Go-скиллы ECC: монолитные, а у samber лучше.

## 3. mattpocock/skills: 319k★, MIT

**Что смотрели.** Все 39 скиллов: `grilling`, `tdd`, `prototype`, `to-spec`, `code-review` и другие.

**Взяли:**
- `grilling`: интервью раундами, на каждый вопрос сразу предлагается рекомендуемый ответ, факты агент ищет сам. Добавлен хакатон-режим (2–3 раунда, обязательные темы, передача в `/opsx:propose`).
- `prototype`: быстрые варианты UI или логики. UX весит 25%.
- Из `tdd`: вертикальные срезы и запрет тавтологических проверок (вошли в наш TDD).

**Не взяли.** `to-spec`, `to-tickets`, `wayfinder`, `triage`: они завязаны на трекер задач, а у нас эту роль играет OpenSpec.

## 4. Fission-AI/OpenSpec: 74k★, v1.16.1, MIT

**Почему он, а не spec-kit.** Кейс разрешает оба. OpenSpec легче: цикл explore → propose → apply → archive, изменения хранятся как дельты к постоянным спекам. В бенчмарке Uvik он вышел дешевле и точнее (84% успеха за $2.71 против 80% за $3.33 у spec-kit). Ещё у него есть `config.yaml` с полями `context` и `rules`, которые подмешиваются в каждый артефакт. Через них мы встроили требования кейса.

**Как используем.** Сам OpenSpec не вендорим: его ставит `npm` и `openspec init --tools claude` (CLI сам кладёт свои скиллы и `/opsx:*` в проект). Кит даёт:
- `openspec/config.yaml` с контекстом Авито: стек, ограничения кейса, русский язык, правила для proposal, specs, design и tasks;
- скилл `spec-authoring`: шаблоны требований для поведения агента и контрактов тулзов, три обязательных вида сценариев (успех, нет прав, сбой), состав `design.md`.

## 5. github/spec-kit: 149k★, v1.3.3, MIT

Смотрели 9 команд (около 2.3k строк промптов) и шаблон constitution. Идея constitution (неизменяемые принципы проекта) у нас реализована через `CLAUDE.md` и `context` в `config.yaml`. Сам spec-kit не взяли: процесс с жёсткими фазами даёт больше артефактов на каждое изменение и медленнее на хакатоне, а кейсу достаточно одного SDD-инструмента.

## 6. anthropics/skills: 190k★, Apache-2.0

**Взяли.** `frontend-design` без изменений, вместе с LICENSE.txt. Из `mcp-builder` взяли `reference/mcp_best_practices.md`, тоже с лицензией: это принципы именования, ошибок, пагинации и аннотаций тулзов.

**Не взяли:**
- Код `mcp-builder`: он для TypeScript и Python, поэтому для Go написан свой скилл `go-mcp-server`.
- `webapp-testing`: он на Python Playwright, а у нас `@playwright/test` на TS.
- `pptx`: не вендорим, это отдельная установка; в `pitch-deck` дан Marp как основной путь.

## 7. anthropics/claude-plugins-official: 47k★, Apache-2.0

Смотрели `champion-kit` (его `/hackathon` для **организаторов** хакатонов, не для участников), `feature-dev`, `code-review`, `security-guidance` (Python-хуки), `ralph-loop` (автономный цикл, рискованный на дедлайне), `mcp-server-dev` (TS, MCPB), LSP-плагины.

**Взяли.** Формат `lspServers` для `gopls` и `typescript-language-server` (объявлены прямо в нашем `plugin.json`) и конфиги MCP для Context7 и Playwright.

## 8. samber/cc-skills-golang: 3.8k★, MIT

40+ Go-скиллов с собственными evals (прирост +32…+47 п.п.). Они ограничены `.go`-файлами через `paths: **/*.go`, а детали вынесены в `references/` и подгружаются по требованию.

**Взяли 8:** `golang-testing`, `golang-error-handling`, `golang-database`, `golang-security`, `golang-concurrency`, `golang-context`, `golang-design-patterns`, `golang-project-layout`. Описания сокращены примерно с 500–800 до ~230 символов с сохранением триггеров. Каталоги evals убраны, перекрёстные ссылки переведены на `hackkit:`.

**Не взяли.** Скиллы про samber/*-библиотеки, gRPC, GraphQL, CLI, Wire, Fx: для кейса не нужны.

## 9. Инструменты, подключённые через `.mcp.json` плагина

- **microsoft/playwright-mcp** (39k★): e2e-проверки и **запись видеодемо** с флагом `--caps=devtools` (`browser_start_video`, главы с карточками, подсветка действий); `--caps=pdf` сохраняет PDF.
- **upstash/context7** (64k★): актуальная документация TanStack, Zustand, Zod, MCP Go SDK и Anthropic Go SDK вместо догадок модели.

## 10. SDK, на которых проверен референсный код

- **modelcontextprotocol/go-sdk** v1.8.0, официальный: `mcp.AddTool`, `NewStreamableHTTPHandler`, `auth.RequireBearerToken`, `req.Extra.TokenInfo`.
- **anthropics/anthropic-sdk-go** v1.84.0: `Messages.New`, tool use, `ToParam()` для повтора хода ассистента.

Референсный код (`skills/llm-agent-go/reference`, `skills/go-mcp-server/reference`) собран и прогнан: 10 тестов, `-race`, покрытие 85.4%. Тесты проверяют фильтрацию тулзов по роли, отказ без права, отсутствие утечки чужих заявок, таймаут агента и адаптер Anthropic против фейкового API-сервера.

## 11. Рассмотрели и отклонили

| Репозиторий | Почему нет |
|---|---|
| open-gsd/gsd-core (наследник GSD, 64k★ у архивного) | 72 команды, 64 агента, 72 скилла; процесс из многих фаз рассчитан на долгие проекты |
| bmad-code-org/BMAD-METHOD (55k★) | тяжёлый агентный процесс уровня «enterprise» |
| wshobson/agents (42k★) | Go почти не покрыт |
| nizos/tdd-guard (2.4k★) | блокирующие хуки с LLM-проверками: дорого и медленно на хакатоне; TDD закрыт скиллом |
| hesreallyhim/awesome-claude-code, davila7/claude-code-templates | каталоги: полезны для поиска, ставить из них нечего |
| gotalab/cc-sdd (3.8k★) | ещё один SDD-фреймворк; жюри называет OpenSpec и spec-kit |

## 12. Материалы кейса Авито (страница «Источники информации»)

Все восемь источников прочитаны 9.11.2026: чек-лист стажёра (оба листа, 48 пунктов), AvitoTech Playbook (28 документов), блог AvitoTech на Хабре (100 последних статей по RSS, 16 ключевых прочитаны), avito.tech, репозиторий tech-internship (задания 2023–2026 и разборы ошибок), статья о дизайн-системе, сайт стажировок, доклад о дизайн-системе (на Rutube нет субтитров, а поток защищён, поэтому разобраны слайды через встроенный браузер). Лицензий у репозиториев Авито нет, так что в ките пересказ и ссылки, без копий текстов.

Что это изменило в ките:

| Вывод из материалов | Изменение |
|---|---|
| В Авито уже есть платформа ИИ-агентов «Виталик» (PortalManager + ADK, RAG, MCP Hub, бенчмарки), кейс прямо просит не строить её заново | бадди позиционируется как доменный агент платформы: в `llm-agent-go`, `spec-authoring`, `openspec/config.yaml`, README-шаблоне, презентации и у `jury` |
| Инженеры Авито в разборах стажёрских работ раз за разом отмечают одни и те же ошибки (запуск, слои, миграции, линтеры, loading/error/empty) | `avito-context/references/review-checklist.md`; дополнения в `go-reviewer`, `react-reviewer`, `security-reviewer`, `jury`; предупреждения в `submission-check` |
| Задание Fullstack 2026 «каталог AI-ассистентов» почти совпадает с нашим кейсом: мок LLM, статусы запусков, системный промпт только на бэкенде, `docker-compose up` без ручной настройки | правила в `llm-agent-go`; найден и исправлен риск недопуска: шаблон compose требовал `.env`, теперь у всех переменных значения по умолчанию |
| Авито просит конфиг линтера в репозитории | шаблон `golangci.yml` (golangci-lint v2), проверен на референсном коде: нашёл две мелочи, они исправлены |
| Практики MCP от инженеров Авито: «один инструмент — одна задача», tolerant reader, компактные ответы | `go-mcp-server` |
| Tech Design Review — стандартный формат дизайн-документа в Авито | структура `design.md` в `spec-authoring` и `openspec/config.yaml` |
| Официальный чек-лист стажёра по неделям | сид `kickoff/templates/seed/onboarding-checklist.yaml` и кандидаты сценариев с метриками для `grilling` |
| Дизайн-система AKITA: семантические токены, чёрный акцент, размеры контролов, светлая и тёмная тема; бот DevProblems с «призвать дежурного» | заметки для UI и эскалация «Позвать наставника» в поведении агента |
| Руководитель тестирования Авито называет SpecKit, Superpowers и OpenSpec и связывает критерии приёмки с тестами | подтверждает выбор процесса кита; правило «каждому сценарию — тест» усилено в `spec-authoring` |

## Свои части (написаны для кейса)

`using-hackkit` (маршрутизация и правила), `kickoff` (каркас и шаблоны), `spec-authoring`, `go-mcp-server`, `llm-agent-go`, `parallel-sessions`, `e2e-demo`, `submission-check` (со скриптом), `pitch-deck`, агент `jury`. Шаблоны проекта: `CLAUDE.md`, `docs/CASE.md`, `README.md` для сдачи, `openspec/config.yaml`, `Makefile`, `scripts/coverage-check.sh` (проверен на Go и Vitest), `docker-compose.yml`, CI, `.env.example`, `docs/STATUS.md`.
