package server

import (
	"context"
	"crypto/sha1"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/idpzero/idpzero/pkg/configuration"
	_ "modernc.org/sqlite"
)

const flowConfigYAML = `
server:
  port: 4379
  keyphrase: flow-keyphrase
clients:
  - name: Confidential Web
    client_id: web
    auth_method: client_secret_basic
    access_token_type: bearer
    grant_types: [authorization_code]
    response_types: [code]
    redirect_uris: [http://rp.example/callback]
    post_logout_redirect_uris: [http://rp.example/signed-out]
    custom_scopes:
      myscope: ["Custom:Sample"]
users:
  - subject: alice
    login_display: Alice
    claims:
      email: alice@example.com
      custom:
        groups: [admin, editor]
        "Custom:Sample": Sample
`

const rpCallback = "http://rp.example/callback"

// newFlowServer runs a full idpzero (provider, storage, web routes) over
// httptest with flowConfigYAML.
func newFlowServer(t *testing.T) *httptest.Server {
	t.Helper()
	dir := t.TempDir()
	cm, err := configuration.NewConfigurationManager(dir)
	if err != nil {
		t.Fatalf("NewConfigurationManager: %v", err)
	}
	t.Cleanup(cm.Close)
	if err := os.WriteFile(cm.GetConfigurationFilePath(), []byte(flowConfigYAML), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}
	db, err := sql.Open("sqlite", cm.GetStateDatabasePath())
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	s, err := NewServer(slog.New(slog.NewTextHandler(io.Discard, nil)), cm, db)
	if err != nil {
		t.Fatalf("NewServer: %v", err)
	}
	if err := s.onStart(context.Background()); err != nil {
		t.Fatalf("onStart: %v", err)
	}
	ts := httptest.NewServer(s.server.Handler)
	t.Cleanup(ts.Close)
	return ts
}

// noRedirects is a client that returns each redirect instead of following it.
var noRedirects = &http.Client{CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}

func redirectTo(t *testing.T, res *http.Response) *url.URL {
	t.Helper()
	defer func() { _ = res.Body.Close() }()
	if res.StatusCode != http.StatusFound {
		body, _ := io.ReadAll(res.Body)
		t.Fatalf("%s %s: status %d, want 302: %s", res.Request.Method, res.Request.URL.Path, res.StatusCode, body)
	}
	loc, err := res.Request.URL.Parse(res.Header.Get("Location"))
	if err != nil {
		t.Fatalf("Location: %v", err)
	}
	return loc
}

type tokens struct {
	AccessToken string `json:"access_token"`
	IDToken     string `json:"id_token"`
}

// login signs alice in to client web with scope and returns the tokens.
func login(t *testing.T, ts *httptest.Server, scope string) tokens {
	t.Helper()
	q := url.Values{
		"client_id": {"web"}, "redirect_uri": {rpCallback}, "response_type": {"code"},
		"scope": {scope}, "state": {"st"}, "nonce": {"nc"},
	}
	res, err := noRedirects.Get(ts.URL + "/authorize?" + q.Encode())
	if err != nil {
		t.Fatal(err)
	}
	loginURL := redirectTo(t, res)
	req := loginURL.Query().Get("req")
	if loginURL.Path != "/login" || req == "" {
		t.Fatalf("authorize redirected to %s, want /login?req=", loginURL)
	}

	res, err = noRedirects.PostForm(ts.URL+"/login", url.Values{"username": {"alice"}, "req": {req}})
	if err != nil {
		t.Fatal(err)
	}
	res, err = noRedirects.Get(redirectTo(t, res).String())
	if err != nil {
		t.Fatal(err)
	}
	cb := redirectTo(t, res)
	code := cb.Query().Get("code")
	if !strings.HasPrefix(cb.String(), rpCallback) || code == "" {
		t.Fatalf("callback redirected to %s, want %s?code=", cb, rpCallback)
	}

	tr, _ := http.NewRequest(http.MethodPost, ts.URL+"/oauth/token", strings.NewReader(url.Values{
		"grant_type": {"authorization_code"}, "code": {code}, "redirect_uri": {rpCallback},
	}.Encode()))
	tr.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	secret := sha1.Sum([]byte("flow-keyphrase" + "web"))
	tr.SetBasicAuth("web", hex.EncodeToString(secret[:]))
	res, err = http.DefaultClient.Do(tr)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = res.Body.Close() }()
	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		t.Fatalf("token: status %d: %s", res.StatusCode, body)
	}
	var tok tokens
	if err := json.NewDecoder(res.Body).Decode(&tok); err != nil {
		t.Fatal(err)
	}
	return tok
}

// jwtClaims decodes a JWT payload without verifying it.
func jwtClaims(t *testing.T, raw string) map[string]any {
	t.Helper()
	parts := strings.Split(raw, ".")
	if len(parts) != 3 {
		t.Fatalf("not a JWT: %d parts", len(parts))
	}
	b, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		t.Fatal(err)
	}
	var claims map[string]any
	if err := json.Unmarshal(b, &claims); err != nil {
		t.Fatal(err)
	}
	return claims
}

func userinfo(t *testing.T, ts *httptest.Server, accessToken string) map[string]any {
	t.Helper()
	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/userinfo", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = res.Body.Close() }()
	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		t.Fatalf("userinfo: status %d: %s", res.StatusCode, body)
	}
	var claims map[string]any
	if err := json.NewDecoder(res.Body).Decode(&claims); err != nil {
		t.Fatal(err)
	}
	return claims
}

func TestCustomClaims(t *testing.T) {
	ts := newFlowServer(t)

	tests := []struct {
		scope      string
		wantSample bool
	}{
		// groups is mapped by no custom scope, so it comes with openid;
		// Custom:Sample is mapped by myscope, so only with myscope.
		{scope: "openid email", wantSample: false},
		{scope: "openid email myscope", wantSample: true},
	}
	for _, tt := range tests {
		t.Run(tt.scope, func(t *testing.T) {
			tok := login(t, ts, tt.scope)
			for name, claims := range map[string]map[string]any{
				"id_token": jwtClaims(t, tok.IDToken),
				"userinfo": userinfo(t, ts, tok.AccessToken),
			} {
				groups, _ := json.Marshal(claims["groups"])
				if string(groups) != `["admin","editor"]` {
					t.Errorf("%s groups = %s, want [\"admin\",\"editor\"]", name, groups)
				}
				if _, ok := claims["Custom:Sample"]; ok != tt.wantSample {
					t.Errorf("%s has Custom:Sample = %v, want %v", name, ok, tt.wantSample)
				}
				if claims["email"] != "alice@example.com" || claims["sub"] != "alice" {
					t.Errorf("%s standard claims: sub=%v email=%v", name, claims["sub"], claims["email"])
				}
			}
		})
	}
}

func TestEndSession(t *testing.T) {
	ts := newFlowServer(t)
	tok := login(t, ts, "openid")

	t.Run("id_token_hint and post_logout_redirect_uri", func(t *testing.T) {
		q := url.Values{
			"id_token_hint":            {tok.IDToken},
			"post_logout_redirect_uri": {"http://rp.example/signed-out"},
			"state":                    {"bye"},
		}
		res, err := noRedirects.Get(ts.URL + "/end_session?" + q.Encode())
		if err != nil {
			t.Fatal(err)
		}
		if got := redirectTo(t, res).String(); got != "http://rp.example/signed-out?state=bye" {
			t.Errorf("redirect = %s", got)
		}
	})

	t.Run("unregistered post_logout_redirect_uri", func(t *testing.T) {
		q := url.Values{"client_id": {"web"}, "post_logout_redirect_uri": {"http://evil.example/"}}
		res, err := noRedirects.Get(ts.URL + "/end_session?" + q.Encode())
		if err != nil {
			t.Fatal(err)
		}
		defer func() { _ = res.Body.Close() }()
		if loc := res.Header.Get("Location"); strings.HasPrefix(loc, "http://evil.example") {
			t.Errorf("redirected to unregistered %s", loc)
		}
	})

	t.Run("no parameters lands on the logged-out page", func(t *testing.T) {
		res, err := noRedirects.Get(ts.URL + "/end_session")
		if err != nil {
			t.Fatal(err)
		}
		loc := redirectTo(t, res)
		if loc.Path != pathLoggedOut {
			t.Fatalf("redirect = %s, want %s", loc, pathLoggedOut)
		}
		res, err = http.Get(loc.String())
		if err != nil {
			t.Fatal(err)
		}
		defer func() { _ = res.Body.Close() }()
		body, _ := io.ReadAll(res.Body)
		if res.StatusCode != http.StatusOK || !strings.Contains(string(body), "Signed Out") {
			t.Errorf("logged-out page: status %d", res.StatusCode)
		}
	})
}
