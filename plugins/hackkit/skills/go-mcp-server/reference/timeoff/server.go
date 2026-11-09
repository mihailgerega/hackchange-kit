// Package timeoff is a mock of the corporate absence service, exposed as MCP tools.
package timeoff

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"sync"
	"time"

	"github.com/modelcontextprotocol/go-sdk/auth"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Scopes granted by the (already existing) corporate auth.
const (
	ScopeRead  = "timeoff:read"
	ScopeWrite = "timeoff:write"
)

// Request is a stored time-off request.
type Request struct {
	ID       string `json:"id"`
	UserID   string `json:"user_id"`
	Date     string `json:"date" jsonschema:"day off, YYYY-MM-DD"`
	Reason   string `json:"reason"`
	Status   string `json:"status"`
	Approver string `json:"approver"`
}

// Store is an in-memory fake of the absence system.
type Store struct {
	mu       sync.Mutex
	requests []Request
	managers map[string]string // user id -> manager id
}

func NewStore(managers map[string]string) *Store {
	return &Store{managers: managers}
}

type CreateInput struct {
	Date   string `json:"date" jsonschema:"day off in YYYY-MM-DD format"`
	Reason string `json:"reason" jsonschema:"short reason shown to the manager"`
}

var errForbidden = errors.New("forbidden: your role cannot create time-off requests")

// caller returns the identity the auth middleware attached to the request.
// Never take the user id from tool arguments: the model could put anything there.
func caller(req *mcp.CallToolRequest, scope string) (string, error) {
	ti := req.Extra.TokenInfo
	if ti == nil || ti.UserID == "" {
		return "", errors.New("unauthenticated")
	}
	if !slices.Contains(ti.Scopes, scope) {
		return "", errForbidden
	}
	return ti.UserID, nil
}

func (s *Store) create(ctx context.Context, req *mcp.CallToolRequest, in CreateInput) (*mcp.CallToolResult, Request, error) {
	userID, err := caller(req, ScopeWrite)
	if err != nil {
		return nil, Request{}, err
	}
	if _, err := time.Parse(time.DateOnly, in.Date); err != nil {
		return nil, Request{}, fmt.Errorf("date must be YYYY-MM-DD, got %q", in.Date)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	r := Request{
		ID:       fmt.Sprintf("TO-%d", len(s.requests)+1),
		UserID:   userID,
		Date:     in.Date,
		Reason:   in.Reason,
		Status:   "pending_approval",
		Approver: s.managers[userID],
	}
	s.requests = append(s.requests, r)
	return nil, r, nil
}

type ListInput struct{}

type ListOutput struct {
	Requests []Request `json:"requests"`
}

func (s *Store) listMine(ctx context.Context, req *mcp.CallToolRequest, _ ListInput) (*mcp.CallToolResult, ListOutput, error) {
	userID, err := caller(req, ScopeRead)
	if err != nil {
		return nil, ListOutput{}, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	out := ListOutput{Requests: []Request{}}
	for _, r := range s.requests {
		if r.UserID == userID { // scope every query to the caller
			out.Requests = append(out.Requests, r)
		}
	}
	return nil, out, nil
}

// NewServer builds the MCP server. Tools a caller may not use are not registered
// for that caller's session, so the model never even sees them.
func NewServer(s *Store, scopes []string) *mcp.Server {
	srv := mcp.NewServer(&mcp.Implementation{Name: "timeoff", Version: "v0.1.0"}, nil)
	if slices.Contains(scopes, ScopeRead) {
		mcp.AddTool(srv, &mcp.Tool{
			Name:        "timeoff_list_my_requests",
			Description: "List the current user's time-off requests with their status.",
			Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
		}, s.listMine)
	}
	if slices.Contains(scopes, ScopeWrite) {
		mcp.AddTool(srv, &mcp.Tool{
			Name:        "timeoff_create_request",
			Description: "Create a time-off request for the current user. It goes to the user's manager for approval. Ask the user to confirm date and reason first.",
		}, s.create)
	}
	return srv
}

// Handler serves MCP over streamable HTTP behind bearer-token auth.
func Handler(s *Store, verify auth.TokenVerifier) http.Handler {
	h := mcp.NewStreamableHTTPHandler(func(r *http.Request) *mcp.Server {
		ti := auth.TokenInfoFromContext(r.Context())
		if ti == nil {
			return NewServer(s, nil)
		}
		return NewServer(s, ti.Scopes)
	}, nil)
	return auth.RequireBearerToken(verify, nil)(h)
}
