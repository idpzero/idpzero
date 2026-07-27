package scim

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// A target behind an auth gateway may redirect an unauthenticated API call to an
// HTML login page. The client must not follow it and must report the redirect
// clearly rather than failing to parse HTML as JSON.
func TestFindByExternalID_RedirectToLoginIsReported(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/Users" {
			http.Redirect(w, r, "/auth/login", http.StatusFound)
			return
		}
		// The login page an auto-follow would have landed on.
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("<html><body>login</body></html>"))
	}))
	defer srv.Close()

	_, err := New(srv.URL, "bad-token").FindByExternalID(context.Background(), "someone")
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if !strings.Contains(err.Error(), "redirect") || !strings.Contains(err.Error(), "/auth/login") {
		t.Errorf("error should mention the redirect and location, got: %v", err)
	}
	if strings.Contains(err.Error(), "invalid character") {
		t.Errorf("error should not be a raw JSON-decode failure, got: %v", err)
	}
}

// A target that ignores the externalId filter and returns an unrelated user must
// not be treated as a match (which would replace the wrong record); the client
// must report that no genuine match was found.
func TestFindByExternalID_IgnoredFilterIsNotAMatch(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Always return the same unrelated user regardless of the filter.
		writeJSON(w, http.StatusOK, ListResponse{
			Schemas:      []string{ListSchema},
			TotalResults: 1,
			Resources:    []User{{ID: "some-other-id", ExternalID: "someone-else", UserName: "admin"}},
		})
	}))
	defer srv.Close()

	_, err := New(srv.URL, "").FindByExternalID(context.Background(), "john_verified_email")
	if err == nil {
		t.Fatal("expected an error when the returned resource does not match the requested externalId")
	}
	if !strings.Contains(err.Error(), "does not appear to support filtering") {
		t.Errorf("error should explain the filter is not honoured, got: %v", err)
	}
}

// A 2xx response whose body is HTML (not JSON) must be reported as a non-JSON
// response, not as a cryptic JSON parse error.
func TestFindByExternalID_HTMLBodyWith200IsReported(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("<!doctype html><html>nope</html>"))
	}))
	defer srv.Close()

	_, err := New(srv.URL, "").FindByExternalID(context.Background(), "someone")
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if !strings.Contains(err.Error(), "expected a JSON response") {
		t.Errorf("error should call out the non-JSON response, got: %v", err)
	}
}
