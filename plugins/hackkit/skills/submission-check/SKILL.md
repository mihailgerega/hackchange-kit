---
name: submission-check
description: Check the Hack&Change submission against the case requirements - admission blockers, README sections, LLM config, secrets, SDD specs, coverage report, demo video, deck PDF, jury access - and run the clean-clone test. Use before every submission or demo freeze.
---

# Submission check

Run in this order and report a table of ✅/❌/⚠️ with the fix for each ❌. Admission blockers first: if the expert cannot start the prototype by the README, or there is no code, or no presentation — the team gets 0.

## 1. Static checklist

```bash
bash "${CLAUDE_PLUGIN_ROOT}/skills/submission-check/check.sh"
```

It checks README sections (run, tests, test users, check order, versions, LLM service/model/env var), `.env.example` and committed secrets, compose validity and that it starts without `.env`, what Avito reviewers look for (linter configs, separate migrations, decisions, C4/CJM diagrams, DB schema, LLM request example, run statuses — as warnings), OpenSpec specs and `openspec validate --all`, permission scenarios, `coverage/REPORT.md` (> 80% per part), e2e presence, test-data labeling, video link, presentation PDF. The README skeleton with every required section is in `${CLAUDE_PLUGIN_ROOT}/skills/kickoff/templates/README.md`.

## 2. Clean-clone run (the expert's experience)

```bash
rm -rf /tmp/buddy-check && git clone <repo> /tmp/buddy-check && cd /tmp/buddy-check
make up && make e2e && make down    # no .env: defaults must be enough
```

Use the URL the jury will get (the pushed branch), not the working tree. Time it: if it takes more than ~10 minutes or needs anything not in the README, fix the README or the setup.

## 3. Fresh coverage

`make cover-check` on the final commit, commit `coverage/REPORT.md` and the HTML reports (or attach them), so the numbers in the deck match the repo.

## 4. Things a human confirms

- Private repo: GitHub account `hack-and-change-2026` has access (or the ZIP contains code + deck).
- Video link on the first slide and at the top of the README; the video matches the final build.
- Deck exported to PDF (`hackkit:pitch-deck`), with the repo link.
- Upload done before the deadline — the version is frozen after it.

Do not push, upload or change repository access yourself; prepare and let the user do the outward-facing step.

## 5. Mock jury

Dispatch the `hackkit:jury` agent on the repo and deck for a criterion-by-criterion score with the top fixes; spend the remaining time on the highest-weight gaps.
