# Источники из кейса (страница «Источники информации»)

Все прочитаны 9.11.2026. Лицензии у репозиториев avito-tech/playbook и avito-tech/tech-internship не указаны, поэтому в ките только пересказ и ссылки, без копирования текстов. В проекте на эти материалы ссылаемся, а свои тестовые инструкции явно помечаем как тестовые.

| # | Источник | Что внутри | Где пересказ |
|---|---|---|---|
| 1 | Чек-лист стажёра Авито — https://docs.google.com/spreadsheets/d/1qKFqOzfSZ2jEz57jBSmq7rukDXi_tRfWohgemViceV0/edit?usp=sharing | 2 листа: «Материалы для всех» (16 пунктов для стажёра и бадди, от первой недели до выпускного) и «Чек-лист Backend» (32 пункта по срокам) | `intern-journey.md`; сид — `kickoff/templates/seed/onboarding-checklist.yaml` |
| 2 | AvitoTech Playbook — https://github.com/avito-tech/playbook | 28 документов: культура, процессы, принципы разработки, TDR, PaaS, SPT, OKR, профили E1–E8 | `engineering.md` |
| 3 | Блог AvitoTech на Хабре — https://habr.com/ru/companies/avito/articles/ | статьи о внутренних инструментах; для кейса важнее всего «Виталик», MCP, SDD, Go для агентов, модерация на LLM | `ai-platform.md` |
| 4 | AvitoTech — https://avito.tech/ | цифры, курсы (JavaScript, «Паттерны и практики написания кода»), стажировки | `engineering.md` |
| 5 | AvitoTech Internship — https://github.com/avito-tech/tech-internship | тестовые задания 2023–2026 и разборы частых ошибок кандидатов | `review-checklist.md` |
| 6 | Дизайн-система Авито — https://habr.com/ru/companies/avito/articles/543450/ | устройство библиотеки в Figma: стили, компоненты, спецификации, гайдлайны | `design-system.md` |
| 7 | Стажировки Авито — https://start.avito.ru/internships | программы, условия, FAQ | `intern-journey.md` |
| 8 | Доклад о дизайн-системе — https://rutube.ru/video/6a203a00a9ef7cdd25d811a01cacd61b/ | 35 минут, Ф. Соломин; субтитров нет, разобраны слайды | `design-system.md` |

## Статьи Хабра, на которые опирается `ai-platform.md`

| ID | Название | Зачем |
|---|---|---|
| 1071126 | От LLM-портала до фабрики внутренних агентов: «Виталик» | платформа агентов, MCP Hub, RAG, бенчмарки — куда встраивается бадди |
| 1060190 | Агентская разработка: как обеспечить качество | SDD (SpecKit, Superpowers, OpenSpec), практики навыков и MCP, метрики |
| 1056028, 1056030 | От REST к MCP, части 1 и 2 | крупные инструменты, AI Gateway, prompt injection, когда агент не нужен |
| 1070660 | AI Hardcore Set: доклады про MCP и агентов в проде | тестирование MCP, tolerant reader, OWASP для агентов |
| 1069734 | Чем Go полезен при разработке AI-агентов | обвязка без фреймворка на anthropic-sdk-go, context, параллельные инструменты |
| 1095230 | Как мы внедряли LLM в модерацию Авито | промпт как параметр, GuardRails, кэш префикса, человек при низкой уверенности |
| 1085718 | ИИ-конвейер для разметки продуктового фидбэка | эталонная выборка, micro F1, «что гарантирует, а что нет» |
| 1072244 | Стажировка в AvitoTech (август 2026) | кого берут и что предлагают |
