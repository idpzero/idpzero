package server

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
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
	"path/filepath"
	"strings"
	"testing"
	"time"

	jose "github.com/go-jose/go-jose/v4"
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
    grant_types: [authorization_code, refresh_token]
    response_types: [code]
    redirect_uris: [http://rp.example/callback]
    post_logout_redirect_uris: [http://rp.example/signed-out]
    custom_scopes:
      myscope: ["Custom:Sample"]
  - name: Key Client
    client_id: pkjwt
    auth_method: private_key_jwt
    grant_types: [authorization_code]
    response_types: [code]
    redirect_uris: [http://rp.example/callback]
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
	// Written before the manager starts watching, so no reload races the
	// server starting up.
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "server.yaml"), []byte(flowConfigYAML), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}
	cm, err := configuration.NewConfigurationManager(dir)
	if err != nil {
		t.Fatalf("NewConfigurationManager: %v", err)
	}
	t.Cleanup(cm.Close)
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
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	IDToken      string `json:"id_token"`
}

// webSecret is client web's derived secret, sha1(keyphrase + client_id).
func webSecret() string {
	sum := sha1.Sum([]byte("flow-keyphrase" + "web"))
	return hex.EncodeToString(sum[:])
}

// postForm posts form to path as client web (HTTP Basic) and returns the
// status and body.
func postForm(t *testing.T, ts *httptest.Server, path string, form url.Values) (int, []byte) {
	t.Helper()
	req, _ := http.NewRequest(http.MethodPost, ts.URL+path, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth("web", webSecret())
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = res.Body.Close() }()
	body, _ := io.ReadAll(res.Body)
	return res.StatusCode, body
}

func decodeTokens(t *testing.T, status int, body []byte) tokens {
	t.Helper()
	if status != http.StatusOK {
		t.Fatalf("token: status %d: %s", status, body)
	}
	var tok tokens
	if err := json.Unmarshal(body, &tok); err != nil {
		t.Fatal(err)
	}
	return tok
}

// login signs alice in to client web with scope and returns the tokens.
func login(t *testing.T, ts *httptest.Server, scope string) tokens {
	t.Helper()
	code := authorize(t, ts, "web", scope)
	status, body := postForm(t, ts, "/oauth/token", url.Values{
		"grant_type": {"authorization_code"}, "code": {code}, "redirect_uri": {rpCallback},
	})
	return decodeTokens(t, status, body)
}

// authorize signs alice in to clientID with scope and returns the code.
func authorize(t *testing.T, ts *httptest.Server, clientID, scope string) string {
	t.Helper()
	q := url.Values{
		"client_id": {clientID}, "redirect_uri": {rpCallback}, "response_type": {"code"},
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
	return code
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
	status, claims := userinfoStatus(t, ts, accessToken)
	if status != http.StatusOK {
		t.Fatalf("userinfo: status %d", status)
	}
	return claims
}

func userinfoStatus(t *testing.T, ts *httptest.Server, accessToken string) (int, map[string]any) {
	t.Helper()
	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/userinfo", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = res.Body.Close() }()
	var claims map[string]any
	if res.StatusCode == http.StatusOK {
		if err := json.NewDecoder(res.Body).Decode(&claims); err != nil {
			t.Fatal(err)
		}
	}
	return res.StatusCode, claims
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

func TestRefreshTokenRotation(t *testing.T) {
	ts := newFlowServer(t)
	first := login(t, ts, "openid offline_access")
	if first.RefreshToken == "" {
		t.Fatal("no refresh token for offline_access")
	}

	refresh := func(rt string) (int, []byte) {
		return postForm(t, ts, "/oauth/token", url.Values{"grant_type": {"refresh_token"}, "refresh_token": {rt}})
	}
	status, body := refresh(first.RefreshToken)
	second := decodeTokens(t, status, body)
	if second.RefreshToken == "" || second.RefreshToken == first.RefreshToken {
		t.Fatalf("refresh token not rotated: %q", second.RefreshToken)
	}
	if claims := jwtClaims(t, second.IDToken); claims["sub"] != "alice" {
		t.Errorf("refreshed id_token sub = %v", claims["sub"])
	}
	if got := userinfo(t, ts, second.AccessToken)["sub"]; got != "alice" {
		t.Errorf("refreshed access token userinfo sub = %v", got)
	}

	if status, _ := refresh(first.RefreshToken); status != http.StatusBadRequest {
		t.Errorf("reusing the rotated refresh token: status %d, want 400", status)
	}
	if status, _ := userinfoStatus(t, ts, first.AccessToken); status == http.StatusOK {
		t.Error("access token of the rotated refresh token still works")
	}
	if status, _ := refresh("not-a-token"); status != http.StatusBadRequest {
		t.Errorf("unknown refresh token: status %d, want 400", status)
	}
}

func TestRevocation(t *testing.T) {
	ts := newFlowServer(t)

	t.Run("access token", func(t *testing.T) {
		tok := login(t, ts, "openid")
		if status, body := postForm(t, ts, "/revoke", url.Values{"token": {tok.AccessToken}, "token_type_hint": {"access_token"}}); status != http.StatusOK {
			t.Fatalf("revoke: status %d: %s", status, body)
		}
		if status, _ := userinfoStatus(t, ts, tok.AccessToken); status == http.StatusOK {
			t.Error("revoked access token still works")
		}
	})

	t.Run("refresh token and its access token", func(t *testing.T) {
		tok := login(t, ts, "openid offline_access")
		if status, body := postForm(t, ts, "/revoke", url.Values{"token": {tok.RefreshToken}}); status != http.StatusOK {
			t.Fatalf("revoke: status %d: %s", status, body)
		}
		if status, _ := postForm(t, ts, "/oauth/token", url.Values{"grant_type": {"refresh_token"}, "refresh_token": {tok.RefreshToken}}); status != http.StatusBadRequest {
			t.Errorf("revoked refresh token: status %d, want 400", status)
		}
		if status, _ := userinfoStatus(t, ts, tok.AccessToken); status == http.StatusOK {
			t.Error("access token of the revoked refresh token still works")
		}
	})

	t.Run("unknown token", func(t *testing.T) {
		if status, body := postForm(t, ts, "/revoke", url.Values{"token": {"not-a-token"}}); status != http.StatusOK {
			t.Errorf("revoke unknown: status %d, want 200: %s", status, body)
		}
	})
}

func TestIntrospection(t *testing.T) {
	ts := newFlowServer(t)
	tok := login(t, ts, "openid email")

	introspect := func(token string) map[string]any {
		status, body := postForm(t, ts, "/oauth/introspect", url.Values{"token": {token}})
		if status != http.StatusOK {
			t.Fatalf("introspect: status %d: %s", status, body)
		}
		var res map[string]any
		if err := json.Unmarshal(body, &res); err != nil {
			t.Fatal(err)
		}
		return res
	}

	res := introspect(tok.AccessToken)
	if res["active"] != true || res["sub"] != "alice" || res["client_id"] != "web" || res["email"] != "alice@example.com" {
		t.Errorf("introspection = %v", res)
	}
	if res := introspect("not-a-token"); res["active"] != false {
		t.Errorf("unknown token active = %v", res["active"])
	}

	if status, body := postForm(t, ts, "/revoke", url.Values{"token": {tok.AccessToken}}); status != http.StatusOK {
		t.Fatalf("revoke: status %d: %s", status, body)
	}
	if res := introspect(tok.AccessToken); res["active"] != false {
		t.Errorf("revoked token active = %v", res["active"])
	}
}

func TestPrivateKeyJWTRejected(t *testing.T) {
	ts := newFlowServer(t)
	code := authorize(t, ts, "pkjwt", "openid")

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	signer, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.RS256, Key: key}, (&jose.SignerOptions{}).WithHeader("kid", "k1"))
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now()
	payload, _ := json.Marshal(map[string]any{
		"iss": "pkjwt", "sub": "pkjwt", "aud": ts.URL, "jti": "j1",
		"iat": now.Unix(), "exp": now.Add(time.Minute).Unix(),
	})
	jws, err := signer.Sign(payload)
	if err != nil {
		t.Fatal(err)
	}
	assertion, _ := jws.CompactSerialize()

	res, err := http.PostForm(ts.URL+"/oauth/token", url.Values{
		"grant_type": {"authorization_code"}, "code": {code}, "redirect_uri": {rpCallback}, "client_id": {"pkjwt"},
		"client_assertion_type": {"urn:ietf:params:oauth:client-assertion-type:jwt-bearer"}, "client_assertion": {assertion},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = res.Body.Close() }()
	body, _ := io.ReadAll(res.Body)
	// zitadel reports a key lookup failure as a 400 server_error.
	if res.StatusCode != http.StatusBadRequest || !strings.Contains(string(body), "private_key_jwt is not supported") {
		t.Errorf("private_key_jwt: status %d, want 400 naming the cause: %s", res.StatusCode, body)
	}
}
