package scim

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/idpzero/idpzero/pkg/configuration"
)

// mockSCIM is a tiny in-memory SCIM service provider for tests. Users already
// present (keyed by externalId) are returned on filter queries; anything else is
// reported absent so the provisioner will create it.
type mockSCIM struct {
	existing map[string]User // externalId -> resource (must have an ID)
	failOn   string          // subject whose create/replace returns 500

	posts    int
	puts     int
	authSeen string // last Authorization header observed
}

func (m *mockSCIM) handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m.authSeen = r.Header.Get("Authorization")

		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/Users":
			filter := r.URL.Query().Get("filter")
			ext := extractExternalID(filter)
			list := ListResponse{Schemas: []string{ListSchema}}
			if u, ok := m.existing[ext]; ok {
				list.Resources = []User{u}
				list.TotalResults = 1
			}
			writeJSON(w, http.StatusOK, list)

		case r.Method == http.MethodPost && r.URL.Path == "/Users":
			m.posts++
			var in User
			_ = json.NewDecoder(r.Body).Decode(&in)
			if in.ExternalID == m.failOn {
				writeJSON(w, http.StatusInternalServerError, ErrorResponse{Detail: "boom", Status: "500"})
				return
			}
			in.ID = "id-" + in.ExternalID
			writeJSON(w, http.StatusCreated, in)

		case r.Method == http.MethodPut && strings.HasPrefix(r.URL.Path, "/Users/"):
			m.puts++
			var in User
			_ = json.NewDecoder(r.Body).Decode(&in)
			if in.ExternalID == m.failOn {
				writeJSON(w, http.StatusInternalServerError, ErrorResponse{Detail: "boom", Status: "500"})
				return
			}
			writeJSON(w, http.StatusOK, in)

		default:
			http.Error(w, "not found", http.StatusNotFound)
		}
	}
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", ContentType)
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

// extractExternalID pulls the quoted value out of `externalId eq "<value>"`.
func extractExternalID(filter string) string {
	_, after, found := strings.Cut(filter, "eq ")
	if !found {
		return ""
	}
	return strings.Trim(after, `"`)
}

func users(subjects ...string) []*configuration.User {
	out := make([]*configuration.User, 0, len(subjects))
	for _, s := range subjects {
		out = append(out, &configuration.User{Subject: s})
	}
	return out
}

func actionFor(r Report, subject string) string {
	for _, res := range r.Results {
		if res.Subject == subject {
			return res.Action
		}
	}
	return ""
}

func TestSync_CreatesAndUpdates(t *testing.T) {
	mock := &mockSCIM{existing: map[string]User{
		"existing": {ID: "id-existing", ExternalID: "existing", UserName: "existing"},
	}}
	srv := httptest.NewServer(mock.handler())
	defer srv.Close()

	client := New(srv.URL, "")
	report := Sync(context.Background(), client, srv.URL, users("newcomer", "existing"))

	if report.Failed != 0 {
		t.Fatalf("Failed = %d, want 0; results=%+v", report.Failed, report.Results)
	}
	if got := actionFor(report, "newcomer"); got != ActionCreated {
		t.Errorf("newcomer action = %q, want created", got)
	}
	if got := actionFor(report, "existing"); got != ActionUpdated {
		t.Errorf("existing action = %q, want updated", got)
	}
	if mock.posts != 1 {
		t.Errorf("posts = %d, want 1", mock.posts)
	}
	if mock.puts != 1 {
		t.Errorf("puts = %d, want 1", mock.puts)
	}
}

func TestSync_PartialFailureDoesNotAbort(t *testing.T) {
	mock := &mockSCIM{existing: map[string]User{}, failOn: "bad"}
	srv := httptest.NewServer(mock.handler())
	defer srv.Close()

	client := New(srv.URL, "")
	report := Sync(context.Background(), client, srv.URL, users("good1", "bad", "good2"))

	if report.Failed != 1 {
		t.Fatalf("Failed = %d, want 1; results=%+v", report.Failed, report.Results)
	}
	if got := actionFor(report, "bad"); got != ActionFailed {
		t.Errorf("bad action = %q, want failed", got)
	}
	if got := actionFor(report, "good1"); got != ActionCreated {
		t.Errorf("good1 action = %q, want created", got)
	}
	if got := actionFor(report, "good2"); got != ActionCreated {
		t.Errorf("good2 action = %q, want created", got)
	}
}

func TestSync_AuthHeaderOnlyWhenTokenSet(t *testing.T) {
	// With a token, the Authorization header is present.
	mock := &mockSCIM{existing: map[string]User{}}
	srv := httptest.NewServer(mock.handler())
	defer srv.Close()

	client := New(srv.URL, "s3cret")
	_ = Sync(context.Background(), client, srv.URL, users("someone"))
	if mock.authSeen != "Bearer s3cret" {
		t.Errorf("auth header = %q, want %q", mock.authSeen, "Bearer s3cret")
	}

	// Without a token, no Authorization header is sent.
	mock2 := &mockSCIM{existing: map[string]User{}}
	srv2 := httptest.NewServer(mock2.handler())
	defer srv2.Close()

	client2 := New(srv2.URL, "")
	_ = Sync(context.Background(), client2, srv2.URL, users("someone"))
	if mock2.authSeen != "" {
		t.Errorf("auth header = %q, want empty", mock2.authSeen)
	}
}
