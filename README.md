# hackchange-kit

Сетап Claude Code для хакатона **Hack&Change 2026** (кейс Авито «AI-buddy: ИИ-наставник для стажёров»). Это маркетплейс плагинов с одним плагином `hackkit`. В нём собраны лучшие части Superpowers, ECC (сетап победителя хакатона Anthropic), mattpocock/skills, samber/cc-skills-golang, официальных скиллов и плагинов Anthropic, интеграция с OpenSpec и свои скиллы под требования кейса.

Почему именно эти части, а не репозитории целиком, разобрано в [docs/ANALYSIS.md](docs/ANALYSIS.md).

## Что получаете

- **Процесс.** `grilling` → OpenSpec (`/opsx:propose`) → план → TDD-сборка субагентами с ревью → проверка → архив. SDD-спеки, которые кейс требует сдать, получаются сами собой.
- **Правила кейса в каждой сессии.** Хук SessionStart подмешивает `using-hackkit` (~1.2k токенов): синтетические данные, права на сервере, покрытие > 80%, отзывчивый UI, приоритеты по весам критериев.
- **Знания об Авито.** Скилл `avito-context` — выжимка всех восьми материалов кейса: путь стажёра по официальному чек-листу (и сид на 48 пунктов), платформа ИИ-агентов «Виталик» и правила MCP от инженеров Авито, принципы Playbook и формат Tech Design Review, чек-лист из разборов стажёрских заданий, заметки о дизайн-системе AKITA.
- **Стек.** Go (8 скиллов samber), React/Vitest/MSW, API, Postgres, Redis, дизайн интерфейса, LSP для Go и TypeScript.
- **Ядро AI-buddy с проверенным кодом.** Цикл агента с интерфейсом модели и фейком, адаптер Anthropic Go SDK, MCP-имитация сервиса отсутствий с авторизацией по токену и фильтром тулзов по роли. Всё собрано и протестировано: 10 тестов, `-race`, покрытие 85%.
- **Ревью.** Агенты `go-reviewer`, `go-build-resolver`, `typescript-reviewer`, `react-reviewer`, `database-reviewer`, `security-reviewer` (с чек-листом для LLM-агента и MCP).
- **Сдача.** Шаблоны README, `CASE.md`, Makefile, скрипт проверки покрытия с воспроизводимым отчётом, CI, Docker Compose; e2e и запись видео через Playwright MCP; презентация в Marp → PDF; проверка сдачи по требованиям кейса; агент `jury`, который оценивает по официальным критериям и весам.
- **MCP.** Playwright (тесты, видео, PDF) и Context7 (актуальная документация библиотек).

Постоянная нагрузка на контекст: около 2.6k токенов на описания скиллов и агентов плюс около 1.2k от хука. Для сравнения, полный ECC занимает около 33k.

## Установка

Нужны Claude Code (CLI или вкладка Code в десктопном приложении), `git`, Node 22+ (для `npx`-серверов MCP и OpenSpec), Go 1.25+, Docker.

### Локально (только вы)

```bash
claude plugin marketplace add /home/misha/hackchange-kit
```
```bash
claude plugin install hackkit@hackchange
```

В десктопном приложении: вкладка Code → **+** рядом с полем ввода → **Add plugins** → добавить маркетплейс по пути → установить `hackkit`.

Плагин из локального маркетплейса читается прямо с диска: правки в `~/hackchange-kit` применяются при следующем старте сессии или после `/reload-plugins`.

### Для всей команды (после публикации на GitHub)

```bash
claude plugin marketplace add <owner>/hackchange-kit --scope project
```
```bash
claude plugin install hackkit@hackchange --scope project
```

Команды выполняются в репозитории проекта. Они записывают `.claude/settings.json`; его нужно закоммитить, и тогда каждый участник получит кит, когда откроет проект. Поле `version` в плагине не задано специально: команда получает новые коммиты кита без ручного повышения версии (`/plugin marketplace update hackchange`).

### Языковые серверы для LSP

```bash
go install golang.org/x/tools/gopls@latest
```
```bash
npm install -g typescript typescript-language-server @fission-ai/openspec@latest
```

## Как работать

1. **Старт проекта.** `/hackkit:kickoff github.com/<team>/ai-buddy` создаёт каркас, OpenSpec с контекстом Авито, шаблоны и проверку `make check`.
2. **Выбор сценариев.** `/hackkit:grilling`: Claude прочитает `docs/CASE.md` и предложит 3–5 сценариев с рекомендацией.
3. **Спека.** `/opsx:propose add-timeoff-request`, по правилам `hackkit:spec-authoring`.
4. **План и сборка.** `hackkit:writing-plans` → `hackkit:subagent-driven-development` (TDD, ревью после каждой задачи).
5. **Готово.** `make check` → `/opsx:archive <id>`.
6. **Сдача.** `/hackkit:submission-check`, `/hackkit:e2e-demo` (видео), `/hackkit:pitch-deck`, затем агент `hackkit:jury`.

Работа нескольких людей одновременно описана в `/hackkit:parallel-sessions`: таблица владения файлами, замороженный контракт, журнал статуса, worktree на каждую сессию.

## Состав

| Группа | Скиллы и агенты | Источник |
|---|---|---|
| Процесс | `using-hackkit`, `kickoff`, `spec-authoring`, `parallel-sessions`, `e2e-demo`, `submission-check`, `pitch-deck`, агент `jury` | свои |
| Контекст Авито | `avito-context` (+ сид чек-листа стажёра в `kickoff`) | материалы кейса, пересказ со ссылками |
| Дисциплина | `test-driven-development`, `writing-plans`, `subagent-driven-development`, `executing-plans`, `verification-before-completion`, `requesting-code-review`, `systematic-debugging`, `dispatching-parallel-agents`, `using-git-worktrees`, `finishing-a-development-branch` | obra/superpowers (адаптировано) |
| Discovery и прототипы | `grilling`, `prototype` | mattpocock/skills (адаптировано) |
| AI-ядро | `llm-agent-go`, `go-mcp-server` (+ проверенный код в `reference/`) | свои; принципы MCP из anthropics/skills |
| Go | `golang-testing`, `golang-error-handling`, `golang-database`, `golang-security`, `golang-concurrency`, `golang-context`, `golang-design-patterns`, `golang-project-layout` | samber/cc-skills-golang |
| Веб и данные | `react-patterns`, `react-testing`, `api-design`, `postgres-patterns`, `redis-patterns` | affaan-m/ECC |
| Дизайн | `frontend-design` | anthropics/skills |
| Ревьюеры | `go-reviewer`, `go-build-resolver`, `typescript-reviewer`, `react-reviewer`, `database-reviewer`, `security-reviewer` | affaan-m/ECC (доработано) |
| MCP и LSP | `playwright` (`--caps=devtools,pdf`), `context7`; `gopls`, `typescript-language-server` | microsoft, upstash; формат из claude-plugins-official |

Скиллы написаны по-английски: модель лучше следует таким инструкциям. Шаблоны проекта (`CLAUDE.md`, `CASE.md`, README, `STATUS.md`) на русском, как и OpenSpec-спеки (это задаётся в `config.yaml`).

## Управление

- Отключить подмешивание правил: `HACKKIT_DISABLE=1 claude`.
- Посмотреть состав и стоимость: `claude plugin details hackkit`.
- Удалить: `claude plugin marketplace remove hackchange`.
- Проверить после правок: `claude plugin validate ./plugins/hackkit && claude plugin validate .`.

## Лицензии

Кит распространяется по MIT. Вендоренные части сохраняют свои лицензии (MIT и Apache-2.0), подробности в [THIRD_PARTY_NOTICES.md](THIRD_PARTY_NOTICES.md).
