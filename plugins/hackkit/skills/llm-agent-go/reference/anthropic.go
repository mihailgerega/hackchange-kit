package agent

import (
	"context"
	"encoding/json"

	"github.com/anthropics/anthropic-sdk-go"
)

// AnthropicModel implements Model with the official Go SDK.
// Configure with env: ANTHROPIC_API_KEY (read by the SDK), LLM_MODEL (default claude-opus-5).
type AnthropicModel struct {
	Client    anthropic.Client
	Model     string
	MaxTokens int64
}

func NewAnthropicModel(model string) *AnthropicModel {
	if model == "" {
		model = "claude-opus-5"
	}
	return &AnthropicModel{Client: anthropic.NewClient(), Model: model, MaxTokens: 16000}
}

func (m *AnthropicModel) Complete(ctx context.Context, system string, history []Message, tools []ToolSpec) (Message, error) {
	params := anthropic.MessageNewParams{
		Model:     m.Model,
		MaxTokens: m.MaxTokens,
		System:    []anthropic.TextBlockParam{{Text: system}},
		Messages:  toAnthropicMessages(history),
		Tools:     toAnthropicTools(tools),
	}
	resp, err := m.Client.Messages.New(ctx, params)
	if err != nil {
		return Message{}, err
	}
	out := Message{Role: RoleAssistant, ProviderData: resp.ToParam()}
	for _, block := range resp.Content {
		switch v := block.AsAny().(type) {
		case anthropic.TextBlock:
			out.Text += v.Text
		case anthropic.ToolUseBlock:
			out.ToolCalls = append(out.ToolCalls, ToolCall{
				ID: v.ID, Name: v.Name, Input: json.RawMessage(v.JSON.Input.Raw()),
			})
		}
	}
	return out, nil
}

func toAnthropicTools(tools []ToolSpec) []anthropic.ToolUnionParam {
	out := make([]anthropic.ToolUnionParam, 0, len(tools))
	for _, t := range tools {
		schema := anthropic.ToolInputSchemaParam{Properties: t.InputSchema["properties"]}
		if req, ok := t.InputSchema["required"].([]any); ok {
			for _, r := range req {
				if s, ok := r.(string); ok {
					schema.Required = append(schema.Required, s)
				}
			}
		}
		tp := anthropic.ToolParam{Name: t.Name, Description: anthropic.String(t.Description), InputSchema: schema}
		out = append(out, anthropic.ToolUnionParam{OfTool: &tp})
	}
	return out
}

func toAnthropicMessages(history []Message) []anthropic.MessageParam {
	out := make([]anthropic.MessageParam, 0, len(history))
	for _, m := range history {
		if p, ok := m.ProviderData.(anthropic.MessageParam); ok {
			out = append(out, p) // replay the assistant turn verbatim (keeps thinking blocks)
			continue
		}
		switch {
		case len(m.ToolResults) > 0:
			blocks := make([]anthropic.ContentBlockParamUnion, 0, len(m.ToolResults))
			for _, r := range m.ToolResults {
				blocks = append(blocks, anthropic.NewToolResultBlock(r.CallID, r.Content, r.IsError))
			}
			out = append(out, anthropic.NewUserMessage(blocks...))
		case m.Role == RoleAssistant:
			out = append(out, anthropic.NewAssistantMessage(anthropic.NewTextBlock(m.Text)))
		default:
			out = append(out, anthropic.NewUserMessage(anthropic.NewTextBlock(m.Text)))
		}
	}
	return out
}
