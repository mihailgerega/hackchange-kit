package timeoff

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/auth"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Fake corporate auth: token -> identity. Real auth already exists at Avito.
var tokens = map[string]*auth.TokenInfo{
	"intern-anna":  {UserID: "anna", Scopes: []string{ScopeRead, ScopeWrite}},
	"intern-boris": {UserID: "boris", Scopes: []string{ScopeRead, ScopeWrite}},
	"guest":        {UserID: "guest", Scopes: []string{ScopeRead}},
}

func verify(_ context.Context, token string, _ *http.Request) (*auth.TokenInfo, error) {
	ti, ok := tokens[token]
	if !ok {
		return nil, auth.ErrInvalidToken
	}
	cp := *ti
	cp.Expiration = time.Now().Add(time.Hour) // the middleware rejects tokens without expiry
	return &cp, nil
}

type bearer struct{ token string }

func (b bearer) RoundTrip(r *http.Request) (*http.Response, error) {
	r = r.Clone(r.Context())
	r.Header.Set("Authorization", "Bearer "+b.token)
	return http.DefaultTransport.RoundTrip(r)
}

func connect(t *testing.T, url, token string) *mcp.ClientSession {
	t.Helper()
	client := mcp.NewClient(&mcp.Implementation{Name: "test", Version: "v0"}, nil)
	cs, err := client.Connect(t.Context(), &mcp.StreamableClientTransport{
		Endpoint:   url,
		HTTPClient: &http.Client{Transport: bearer{token}},
	}, nil)
	if err != nil {
		t.Fatalf("connect as %s: %v", token, err)
	}
	t.Cleanup(func() { cs.Close() })
	return cs
}

func newTestServer(t *testing.T) string {
	store := NewStore(map[string]string{"anna": "maria", "boris": "maria"})
	ts := httptest.NewServer(Handler(store, verify))
	t.Cleanup(ts.Close)
	return ts.URL
}

func TestCreateRequest_InternCreatesOwnRequest(t *testing.T) {
	cs := connect(t, newTestServer(t), "intern-anna")
	res, err := cs.CallTool(t.Context(), &mcp.CallToolParams{
		Name:      "timeoff_create_request",
		Arguments: map[string]any{"date": "2026-11-20", "reason": "exam"},
	})
	if err != nil || res.IsError {
		t.Fatalf("CallTool: err=%v isError=%v content=%v", err, res.IsError, res.Content)
	}
	var got Request
	raw, _ := json.Marshal(res.StructuredContent)
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatal(err)
	}
	if got.UserID != "anna" || got.Approver != "maria" || got.Status != "pending_approval" {
		t.Errorf("got %+v, want request owned by anna, approved by maria, pending", got)
	}
}

func TestCreateRequest_ReadOnlyRoleDoesNotSeeTool(t *testing.T) {
	cs := connect(t, newTestServer(t), "guest")
	tools, err := cs.ListTools(t.Context(), nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, tool := range tools.Tools {
		if tool.Name == "timeoff_create_request" {
			t.Fatal("read-only caller must not be offered timeoff_create_request")
		}
	}
	res, err := cs.CallTool(t.Context(), &mcp.CallToolParams{
		Name:      "timeoff_create_request",
		Arguments: map[string]any{"date": "2026-11-20", "reason": "x"},
	})
	if err == nil && !res.IsError {
		t.Fatal("read-only caller must not be able to call timeoff_create_request")
	}
}

func TestListMyRequests_DoesNotLeakOtherUsers(t *testing.T) {
	url := newTestServer(t)
	anna := connect(t, url, "intern-anna")
	if _, err := anna.CallTool(t.Context(), &mcp.CallToolParams{
		Name: "timeoff_create_request", Arguments: map[string]any{"date": "2026-11-20", "reason": "exam"},
	}); err != nil {
		t.Fatal(err)
	}
	boris := connect(t, url, "intern-boris")
	res, err := boris.CallTool(t.Context(), &mcp.CallToolParams{Name: "timeoff_list_my_requests", Arguments: map[string]any{}})
	if err != nil || res.IsError {
		t.Fatalf("list: %v %v", err, res.Content)
	}
	var out ListOutput
	raw, _ := json.Marshal(res.StructuredContent)
	_ = json.Unmarshal(raw, &out)
	if len(out.Requests) != 0 {
		t.Errorf("boris sees %d requests of other users", len(out.Requests))
	}
}

func TestCreateRequest_RejectsBadDate(t *testing.T) {
	cs := connect(t, newTestServer(t), "intern-anna")
	res, err := cs.CallTool(t.Context(), &mcp.CallToolParams{
		Name: "timeoff_create_request", Arguments: map[string]any{"date": "20.11.2026", "reason": "x"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !res.IsError {
		t.Error("want tool error for non-ISO date")
	}
}

func TestHandler_RejectsUnknownToken(t *testing.T) {
	url := newTestServer(t)
	resp, err := http.Post(url, "application/json", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("status %d, want 401", resp.StatusCode)
	}
}
