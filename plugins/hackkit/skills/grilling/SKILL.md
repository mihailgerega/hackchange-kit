---
name: grilling
description: Interview the user in rounds (numbered questions, each with a recommended answer) until a scenario, plan or decision is fully understood. Use before proposing a new OpenSpec change, to pick hackathon scenarios, or on any 'grill' phrase.
---

Interview the user relentlessly until you reach a shared understanding. Map this as a **design tree**: every decision branches into the decisions that hang off it.

Work the tree in **rounds**. The **frontier** is every decision whose prerequisites are already settled: the questions you can ask _now_ without guessing at answers you haven't heard yet. Ask the whole frontier in one round: number each question and give your recommended answer. Then wait for the user's answers before the next round.

Format a round like so:

```
❓ **Q1** - **<question title>**: <question body, might be multiple paragraphs, including multiple choices>

➡️ <your recommended answer>

---

❓ **Q2** - **<question title>**: <question body, might be multiple paragraphs, including multiple choices>

➡️ <your recommended answer>
```

Word each question so "yes" accepts your recommended answer.

Each round the user answers reshapes the tree: settled decisions push the frontier outward and unblock questions that depended on them. Recompute the frontier and ask the next round. A question whose answer depends on another question still open in this round belongs to a _later_ round, not this one.

Finding _facts_ is your job, never the user's. When a frontier question needs a fact from the environment (filesystem, tools, etc.), dispatch a sub-agent with a fresh context window to find it; don't ask the user for anything you could look up yourself. Don't block on it: a running exploration is an unsettled prerequisite, so only the questions downstream of it wait for the sub-agent to report; ask the rest of the frontier now. The _decisions_ are the user's: put each to them and wait.

The session is done when the frontier is empty: every branch of the design tree visited, nothing left silently assumed. Do not act on it until the user confirms you have reached a shared understanding.

## hackkit: hackathon mode

- **Time-box it.** Two or three rounds is normal at a hackathon. When the remaining frontier would not change what gets built in the next few hours, say so and stop.
- **Ground every question in the case.** Read `CLAUDE.md` and `openspec/config.yaml` first: the case brief, the scoring weights and the constraints there are facts, not questions. For scenario choice also read `hackkit:avito-context` → `references/intern-journey.md` (the official Avito intern checklist by week, candidate scenarios, metrics): propose scenarios from it with a recommendation instead of asking open questions. Ask only what these leave open (which scenarios, which roles, what the intern sees, what counts as success).
- **Always settle these before a change is proposed:** the user and role; the trigger; the happy path end to end; which corporate service is touched (and that it is a mock); what the user is *not* allowed to do; what happens when the LLM or a service fails; how we will show it in the demo.
- **Hand off.** When the user confirms shared understanding, summarize the decisions as a short brief and continue with `/opsx:propose <change-id>` using that brief, following `hackkit:spec-authoring`. Do not write code from a grilling session.
