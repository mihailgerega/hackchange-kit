---
name: parallel-sessions
description: Coordinate several people and Claude sessions on one hackathon repo - file ownership table, frozen API contract with mocks, status log as message bus, one git worktree per session, merge rhythm. Use when splitting work across teammates or running more than one session at once.
---

# Parallel sessions without collisions

The rule that matters: **two sessions never edit the same file.** Everything below exists to make that true and cheap. (Distilled from hackathon write-ups: a frozen-contract + ownership-table split took ~7 minutes to set up and produced zero merge conflicts; a team running 140 PRs in a week used the same "independent chunks, no shared state" rule with worktrees.)

## 1. Ownership table first (before any code)

`docs/STATUS.md` → "Владение". Every path has exactly one owner (a person, and their Claude session). Typical split for 3–5 people:

| Owner | Paths |
|---|---|
| Agent & prompts | `backend/internal/agent/**`, `openspec/specs/buddy-chat/**` |
| Services (MCP mocks) | `backend/internal/services/**`, `backend/internal/auth/**` |
| API & wiring | `backend/internal/httpapi/**`, `backend/cmd/**`, `docker-compose.yml` |
| Frontend | `frontend/**` |
| Specs, docs, deck, video | `openspec/changes/**` (drafts), `docs/**`, `README.md` |

Shared files (`api/openapi.yaml`, `go.mod`, `package.json`, `Makefile`) have one named owner who applies others' requests.

## 2. Freeze the contract, build against mocks

- `api/openapi.yaml` (HTTP between frontend and backend) and the MCP tool contracts in `openspec/specs/*-tools/` are frozen once agreed. Changing them needs a line in the status log and an OK from both sides.
- Frontend builds against MSW handlers that satisfy the contract; flip to the real backend with `VITE_API_URL`. Backend builds against `LLM_PROVIDER=fake`.

## 3. Status log is the message bus

Each session appends one line to its own section of `docs/STATUS.md` (`- [14:05] A3 DONE — timeoff tool + tests, coverage 88%`), reads the others' sections before starting a task, and never edits another section. Requests to another owner go in your own section: `→ C: need GET /api/onboarding/plan, shape in openapi PR #4`.

## 4. One worktree per session

```bash
git worktree add ../buddy-agent -b agent/timeoff-confirmation
cd ../buddy-agent && claude
```

Use `hackkit:using-git-worktrees` for setup and `hackkit:finishing-a-development-branch` to merge back. `main` must stay green and runnable (`make check`, `make up`); merge small, often, after the specialist review (`hackkit:requesting-code-review`).

## 5. Inside one session: parallel subagents

For independent investigations or fixes (different packages, different failing tests) use `hackkit:dispatching-parallel-agents`. Never fan out tasks that touch the same files.

## Failure modes seen at hackathons

- **Cleanup collisions:** a session's cleanup globbed and deleted another's artifacts. Clean up only exact paths you created.
- **Lost outward-facing work:** a draft form or message left in a browser tab got closed by another session. Save drafts to files.
- **Irreversible steps by agents:** submission, pushes to shared branches, publishing the video — agents prepare, a human confirms.
- **Trusting reports:** verify a subagent's "done" with the command output (`hackkit:verification-before-completion`).
