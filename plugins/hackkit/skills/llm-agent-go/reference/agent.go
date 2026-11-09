// Package agent runs the AI-buddy tool-use loop independent of the LLM provider.
package agent

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type Role string

const (
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
)

// ToolSpec is a tool as the model sees it.
type ToolSpec struct {
	Name        string
	Description string
	InputSchema map[string]any // JSON Schema object
}

type ToolCall struct {
	ID    string
	Name  string
	Input json.RawMessage
}

type ToolResult struct {
	CallID  string
	Content string
	IsError bool
}

// Message is a provider-neutral conversation turn.
type Message struct {
	Role        Role
	Text        string
	ToolCalls   []ToolCall
	ToolResults []ToolResult
	// ProviderData holds the provider's own form of an assistant turn so it can be
	// replayed verbatim (e.g. Claude thinking blocks must be sent back unchanged).
	ProviderData any
}

// Model is the only thing that talks to an LLM. Tests use a scripted fake.
type Model interface {
	Complete(ctx context.Context, system string, history []Message, tools []ToolSpec) (Message, error)
}

// Tools is where tools come from: an MCP session opened with the user's own token.
type Tools interface {
	List(ctx context.Context) ([]ToolSpec, error)
	Call(ctx context.Context, call ToolCall) ToolResult
}

// Step is one tool call shown to the user ("what I did").
type Step struct {
	Tool    string
	Input   json.RawMessage
	Output  string
	IsError bool
}

type Result struct {
	Text  string
	Steps []Step
}

type Agent struct {
	Model    Model
	Tools    Tools
	System   string
	MaxTurns int           // hard stop for runaway loops
	Timeout  time.Duration // whole run; the UI shows a spinner meanwhile
}

var ErrTooManyTurns = errors.New("agent: too many tool turns")

// Run answers one user message, calling tools as the model asks.
func (a *Agent) Run(ctx context.Context, history []Message, userText string) (Result, []Message, error) {
	if a.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, a.Timeout)
		defer cancel()
	}
	maxTurns := a.MaxTurns
	if maxTurns <= 0 {
		maxTurns = 8
	}
	tools, err := a.Tools.List(ctx)
	if err != nil {
		return Result{}, history, fmt.Errorf("list tools: %w", err)
	}
	history = append(history, Message{Role: RoleUser, Text: userText})
	var res Result
	for range maxTurns {
		msg, err := a.Model.Complete(ctx, a.System, history, tools)
		if err != nil {
			return res, history, fmt.Errorf("model: %w", err)
		}
		history = append(history, msg)
		if len(msg.ToolCalls) == 0 {
			res.Text = msg.Text
			return res, history, nil
		}
		results := make([]ToolResult, 0, len(msg.ToolCalls))
		for _, call := range msg.ToolCalls {
			r := a.Tools.Call(ctx, call)
			results = append(results, r)
			res.Steps = append(res.Steps, Step{Tool: call.Name, Input: call.Input, Output: r.Content, IsError: r.IsError})
		}
		// All results of one turn go back in a single user message.
		history = append(history, Message{Role: RoleUser, ToolResults: results})
	}
	return res, history, ErrTooManyTurns
}

// MCPTools adapts an MCP client session to Tools.
type MCPTools struct{ Session *mcp.ClientSession }

func (m MCPTools) List(ctx context.Context) ([]ToolSpec, error) {
	lt, err := m.Session.ListTools(ctx, nil)
	if err != nil {
		return nil, err
	}
	out := make([]ToolSpec, 0, len(lt.Tools))
	for _, t := range lt.Tools {
		schema := map[string]any{}
		raw, err := json.Marshal(t.InputSchema)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(raw, &schema); err != nil {
			return nil, err
		}
		out = append(out, ToolSpec{Name: t.Name, Description: t.Description, InputSchema: schema})
	}
	return out, nil
}

func (m MCPTools) Call(ctx context.Context, call ToolCall) ToolResult {
	res, err := m.Session.CallTool(ctx, &mcp.CallToolParams{Name: call.Name, Arguments: call.Input})
	if err != nil {
		// Report the failure to the model instead of aborting: it can explain or retry.
		return ToolResult{CallID: call.ID, Content: "tool failed: " + err.Error(), IsError: true}
	}
	var b strings.Builder
	for _, c := range res.Content {
		if tc, ok := c.(*mcp.TextContent); ok {
			b.WriteString(tc.Text)
		}
	}
	return ToolResult{CallID: call.ID, Content: b.String(), IsError: res.IsError}
}
