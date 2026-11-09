---
name: e2e-demo
description: Playwright end-to-end test for the main scenario plus the required demo video - a demo script, a smoke test that guards admission, and a recorded walkthrough with chapter cards via the Playwright MCP server. Use when adding e2e tests, before merging risky changes, or when recording the video.
---

# E2E smoke test and demo video

The case gives 0 points if the expert can't pass the main scenario, and asks for a video "from request to result". One script serves both: the e2e test proves it works, the recording shows it.

## 1. Demo script first (`docs/DEMO.md`)

Write the path as steps a stranger could follow, with the synthetic user and the exact inputs:

```
Пользователь: стажёр Анна (токен DEMO_INTERN_TOKEN, роль intern)
1. Открыть http://localhost:5173 → виден план адаптации, текущий этап «Неделя 1», 3 следующих шага
2. Спросить: «Как оформить отгул?» → ответ со ссылкой на тестовую инструкцию
3. «Оформи отгул на 20 ноября, у меня экзамен» → карточка предпросмотра
4. «Подтвердить» → заявка TO-1, статус «на согласовании», согласующий — Мария
5. Шаг «Оформить первый отгул» в плане отмечен выполненным
6. Под ролью гостя то же действие недоступно (права)
```

Keep it ≤ 3 minutes on screen. Show the value (time saved, next step), the action through the service, and one permission refusal.

## 2. E2E smoke test (`frontend/e2e/main-scenario.spec.ts`)

- `@playwright/test`, `baseURL` `http://localhost:5173`, runs against `make up` with `LLM_PROVIDER=fake`, so it is deterministic and key-less.
- Select by role/label/text (`getByRole('button', { name: 'Подтвердить' })`), not CSS. Assert the saved result (request id and status), not only that a toast appeared.
- One test per demo step group, plus the permission refusal. Also assert the waiting state appears while the agent runs and disappears.
- `make e2e` runs it; run it after any change to routing, auth, the agent loop or Docker files. It does not count toward unit coverage, it guards admission.

## 3. Record the video (Playwright MCP, bundled with hackkit)

The plugin's `playwright` MCP server starts with `--caps=devtools,pdf`, which adds video tools. With the stack up (`make up`, `LLM_PROVIDER=fake` for a smooth take, or `anthropic` for a live take):

1. `browser_navigate` to the app; set a 1440×900 window.
2. `browser_start_video` with `filename: "demo/ai-buddy-demo.webm"`, `cursor: true`.
3. `browser_video_show_actions` (callouts on each click/typing step).
4. For each part of `docs/DEMO.md`: `browser_video_chapter` with a short Russian title and one-line description (e.g. «Отгул за 30 секунд» / «Бадди оформляет заявку через сервис отсутствий»), then perform the steps.
5. `browser_stop_video`. Check the file plays and shows every step; re-record rather than editing.

Then: add a voice-over or captions if time allows (any editor), convert if the platform needs MP4 (`ffmpeg -i demo/ai-buddy-demo.webm -c:v libx264 -crf 23 demo/ai-buddy-demo.mp4`), put the file in the repo or upload it, and link it on the first slide and at the top of the README. Large files: link instead of committing if they exceed the repo host's limits.

## 4. Before submission

`hackkit:submission-check` re-runs the clean-clone check and `make e2e`. A recording made before the last merge must be re-checked against the final build.
