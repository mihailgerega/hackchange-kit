package agent

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/modelcontextprotocol/go-sdk/auth"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"example.com/buddydemo/internal/services/timeoff"
)

// ScriptedModel replays canned assistant turns and records what it was shown.
type ScriptedModel struct {
	Turns []Message
	Seen  [][]Message
	Err   error
}

func (s *ScriptedModel) Complete(ctx context.Context, _ string, history []Message, _ []ToolSpec) (Message, error) {
	s.Seen = append(s.Seen, append([]Message(nil), history...))
	if s.Err != nil {
		return Message{}, s.Err
	}
	if len(s.Turns) == 0 {
		<-ctx.Done() // behave like a slow LLM
		return Message{}, ctx.Err()
	}
	m := s.Turns[0]
	s.Turns = s.Turns[1:]
	m.Role = RoleAssistant
	return m, nil
}

type bearer struct{ token string }

func (b bearer) RoundTrip(r *http.Request) (*http.Response, error) {
	r = r.Clone(r.Context())
	r.Header.Set("Authorization", "Bearer "+b.token)
	return http.DefaultTransport.RoundTrip(r)
}

func mcpToolsFor(t *testing.T, scopes ...string) Tools {
	t.Helper()
	verify := func(_ context.Context, token string, _ *http.Request) (*auth.TokenInfo, error) {
		return &auth.TokenInfo{UserID: token, Scopes: scopes, Expiration: time.Now().Add(time.Hour)}, nil
	}
	ts := httptest.NewServer(timeoff.Handler(timeoff.NewStore(map[string]string{"anna": "maria"}), verify))
	t.Cleanup(ts.Close)
	client := mcp.NewClient(&mcp.Implementation{Name: "buddy", Version: "v0"}, nil)
	cs, err := client.Connect(t.Context(), &mcp.StreamableClientTransport{
		Endpoint: ts.URL, HTTPClient: &http.Client{Transport: bearer{"anna"}},
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { cs.Close() })
	return MCPTools{Session: cs}
}

func call(name, input string) ToolCall {
	return ToolCall{ID: "call-1", Name: name, Input: json.RawMessage(input)}
}

func TestRun_CreatesTimeOffThroughTool(t *testing.T) {
	model := &ScriptedModel{Turns: []Message{
		{ToolCalls: []ToolCall{call("timeoff_create_request", `{"date":"2026-11-20","reason":"exam"}`)}},
		{Text: "Заявка TO-1 отправлена руководителю."},
	}}
	a := &Agent{Model: model, Tools: mcpToolsFor(t, timeoff.ScopeRead, timeoff.ScopeWrite)}

	res, _, err := a.Run(t.Context(), nil, "оформи отгул на 20 ноября, у меня экзамен")
	if err != nil {
		t.Fatal(err)
	}
	if res.Text != "Заявка TO-1 отправлена руководителю." {
		t.Errorf("text = %q", res.Text)
	}
	if len(res.Steps) != 1 || res.Steps[0].IsError || !strings.Contains(res.Steps[0].Output, `"approver":"maria"`) {
		t.Errorf("steps = %+v", res.Steps)
	}
	last := model.Seen[1][len(model.Seen[1])-1]
	if len(last.ToolResults) != 1 || last.ToolResults[0].CallID != "call-1" {
		t.Errorf("tool result not fed back to the model: %+v", last)
	}
}

func TestRun_ForbiddenToolErrorGoesBackToModel(t *testing.T) {
	model := &ScriptedModel{Turns: []Message{
		{ToolCalls: []ToolCall{call("timeoff_create_request", `{"date":"2026-11-20","reason":"x"}`)}},
		{Text: "У вас нет прав оформить отгул, обратитесь к наставнику."},
	}}
	a := &Agent{Model: model, Tools: mcpToolsFor(t, timeoff.ScopeRead)}

	res, _, err := a.Run(t.Context(), nil, "оформи отгул")
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Steps) != 1 || !res.Steps[0].IsError {
		t.Fatalf("want one failed step, got %+v", res.Steps)
	}
}

func TestRun_StopsAfterMaxTurns(t *testing.T) {
	loop := Message{ToolCalls: []ToolCall{call("timeoff_list_my_requests", `{}`)}}
	model := &ScriptedModel{Turns: []Message{loop, loop, loop}}
	a := &Agent{Model: model, Tools: mcpToolsFor(t, timeoff.ScopeRead), MaxTurns: 2}

	if _, _, err := a.Run(t.Context(), nil, "?"); !errors.Is(err, ErrTooManyTurns) {
		t.Errorf("err = %v, want ErrTooManyTurns", err)
	}
}

func TestRun_TimesOutInsteadOfHanging(t *testing.T) {
	a := &Agent{Model: &ScriptedModel{}, Tools: mcpToolsFor(t, timeoff.ScopeRead), Timeout: 50 * time.Millisecond}

	_, _, err := a.Run(t.Context(), nil, "привет")
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("err = %v, want deadline exceeded", err)
	}
}

func TestAnthropicModel_ParsesToolUse(t *testing.T) {
	var body map[string]any
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, &body)
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"id":"msg_1","type":"message","role":"assistant","model":"claude-opus-5",
		 "content":[{"type":"text","text":"Оформляю."},
		            {"type":"tool_use","id":"toolu_1","name":"timeoff_create_request","input":{"date":"2026-11-20","reason":"exam"}}],
		 "stop_reason":"tool_use","usage":{"input_tokens":10,"output_tokens":5}}`)
	}))
	defer api.Close()
	m := NewAnthropicModel("")
	m.Client = anthropic.NewClient(option.WithBaseURL(api.URL), option.WithAPIKey("test"), option.WithMaxRetries(0))

	msg, err := m.Complete(t.Context(), "Ты AI-buddy.", []Message{{Role: RoleUser, Text: "отгул"}},
		[]ToolSpec{{Name: "timeoff_create_request", Description: "d", InputSchema: map[string]any{
			"type": "object", "properties": map[string]any{"date": map[string]any{"type": "string"}}, "required": []any{"date"},
		}}})
	if err != nil {
		t.Fatal(err)
	}
	if msg.Text != "Оформляю." || len(msg.ToolCalls) != 1 || msg.ToolCalls[0].Name != "timeoff_create_request" {
		t.Fatalf("msg = %+v", msg)
	}
	if body["model"] != "claude-opus-5" {
		t.Errorf("model sent = %v", body["model"])
	}
	tools := body["tools"].([]any)
	if req := tools[0].(map[string]any)["input_schema"].(map[string]any)["required"]; req == nil {
		t.Error("required fields lost in schema conversion")
	}
	if _, ok := msg.ProviderData.(anthropic.MessageParam); !ok {
		t.Error("assistant turn must keep provider data for verbatim replay")
	}
}
