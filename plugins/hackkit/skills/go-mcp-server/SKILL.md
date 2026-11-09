---
name: go-mcp-server
description: Build mock corporate services (absences, access, onboarding checklist, business trips) as MCP servers in Go with the official go-sdk - tool design, bearer-token auth, per-role tool lists, server-side permission checks, and fast tests. Use when adding or changing an MCP tool or service mock.
---

# Go MCP servers for corporate-service mocks

Avito already has these systems (absences, business trips, consents, benefits, newcomer checklist) plus auth, roles, and an AI-agent platform with MCP. We don't rebuild them: each is a small **mock service exposed as MCP tools** that the AI-buddy agent calls with the user's own identity. That is the integration story the jury scores (implementability & integrations: 30% online, 25% final).

Library: `github.com/modelcontextprotocol/go-sdk` (official; `mcp`, `auth` packages). For current API details ask Context7 (`use context7`, library `modelcontextprotocol/go-sdk`). Tool naming, descriptions, pagination and error style: [mcp_best_practices.md](mcp_best_practices.md) (from Anthropic's mcp-builder; its TypeScript/Python code does not apply here).

**Working reference (compiled and tested with go-sdk v1.8.0):** [reference/timeoff/server.go](reference/timeoff/server.go) and [reference/timeoff/server_test.go](reference/timeoff/server_test.go). Start new services by copying that shape.

## Shape of a service

```
internal/services/<service>/
  store.go     in-memory (or Postgres) state + seed data, synthetic only
  server.go    NewServer(store, scopes) *mcp.Server; Handler(store, verify) http.Handler
  server_test.go
```

1. **Typed tools.** `mcp.AddTool(srv, &mcp.Tool{Name, Description, Annotations}, handler)` with `func(ctx, *mcp.CallToolRequest, In) (*mcp.CallToolResult, Out, error)`. Input/output structs with `json` and `jsonschema:"..."` tags generate the schemas; returning `Out` fills `StructuredContent`.
2. **Names** are `<service>_<verb>_<object>` (`timeoff_create_request`, `checklist_list_my_steps`, `access_request_role`). Descriptions say when to use the tool and, for writes, that the user must confirm first.
3. **Annotations**: `ReadOnlyHint: true` on reads; writes leave it false. The agent uses this to decide which calls need a confirmation step in the UI.
4. **Errors**: return a Go `error` from the handler for anything the model should see and explain (bad date, forbidden, not found); the SDK turns it into `IsError: true`. Messages are short, actionable, in the user's language when shown to them.

## Avito's own MCP rules (AvitoTech blog and meetup talks, 2026)

Their engineers judge tools by these; details in `hackkit:avito-context` → `references/ai-platform.md`.

- **One tool = one human task.** `timeoff_create_request`, `checklist_show_my_week`, not CRUD wrappers over REST; fewer calls to reach the result means fewer places for the model to get lost.
- **Few tools**: add a parameter instead of a new tool, but don't merge unrelated actions. One service = one server.
- **Common names** for tools and parameters (`date`, `reason`, `status`), so the model guesses them right.
- **Tolerant reader**: accept imperfect agent input — parse `2026-11-20`, `20.11.2026` and `20 ноября` on the server instead of failing; reject only what is truly ambiguous, with a hint.
- **Compact output**: return only the fields the next step needs; it saves context and avoids leaking data.
- Test the server with a **mid-tier model** too (e.g. the cheapest Claude model): it makes more mistakes and exposes weak names, descriptions and error messages.
- MCP is a **protective facade** over the real services: arguments are validated, limits and audit live in the gateway, data returned by services is untrusted text for the model (indirect prompt injection).

## Identity and permissions (the part the jury will poke)

- Serve over streamable HTTP behind `auth.RequireBearerToken(verify, nil)`. The verifier maps the corporate token to `*auth.TokenInfo{UserID, Scopes, Expiration}` (expiration is required by the middleware).
- **Per-role tool list:** `mcp.NewStreamableHTTPHandler(func(r *http.Request) *mcp.Server { ti := auth.TokenInfoFromContext(r.Context()); return NewServer(store, ti.Scopes) }, nil)` registers only the tools the caller may use, so the model never sees forbidden actions.
- **Check again in the handler:** read `req.Extra.TokenInfo`, require the scope, and scope every query by `ti.UserID`. Never accept a user id or role from tool arguments.
- The backend opens one MCP client session per user with that user's token (`mcp.StreamableClientTransport{Endpoint, HTTPClient}` with a RoundTripper adding `Authorization: Bearer …`). Mentor/lead scenarios use their own tokens and scopes.

## Tests (fast, no network beyond loopback)

Copy the patterns from `reference/timeoff/server_test.go`:

- `httptest.NewServer(Handler(store, fakeVerify))` + an MCP client per test user.
- One test per spec scenario: happy path, **role without the scope does not see the tool and cannot call it**, **user A cannot read user B's data**, invalid input → `IsError`, unknown token → HTTP 401.
- For pure tool logic without auth, `mcp.NewInMemoryTransports()` connects client and server in-process.

## Wiring into the stack

- Each service listens on its own path (`/mcp/timeoff`, `/mcp/checklist`) in the backend binary for the hackathon — or as a separate compose service, which reads even better as "a real system behind MCP" (Avito's 2026 backend intern task asks for exactly such a separate partner service). Document in `design.md` that in production these are adapters in front of Avito's real systems, registered in the MCP Hub of the internal agent platform («Виталик») with an owner, rights and tracing.
- Emit an analytics event per successful write (`timeoff_request_created`) — the metrics in the proposal rely on them.
- State changes are visible in the UI immediately (the case counts an action as done when the service returns a displayed, saved result).
