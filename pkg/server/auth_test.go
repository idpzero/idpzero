package server

import (
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
	"time"
)

// exchange redeems code at client web for tokens.
func exchange(t *testing.T, ts *httptest.Server, code string) tokens {
	t.Helper()
	status, body := postForm(t, ts, "/oauth/token", url.Values{
		"grant_type": {"authorization_code"}, "code": {code}, "redirect_uri": {rpCallback},
	})
	return decodeTokens(t, status, body)
}

// checkAuthClaims checks the acr, amr and auth_time claims of an ID token.
func checkAuthClaims(t *testing.T, idToken, wantACR string, wantAMR []string, notBefore, notAfter time.Time) {
	t.Helper()
	claims := jwtClaims(t, idToken)
	acr, hasACR := claims["acr"]
	switch {
	case wantACR == "" && hasACR:
		t.Errorf("acr = %v, want none", acr)
	case wantACR != "" && acr != wantACR:
		t.Errorf("acr = %v, want %q", acr, wantACR)
	}
	var amr []string
	raw, _ := claims["amr"].([]any)
	for _, v := range raw {
		amr = append(amr, v.(string))
	}
	if !reflect.DeepEqual(amr, wantAMR) {
		t.Errorf("amr = %v, want %v", claims["amr"], wantAMR)
	}
	at, ok := claims["auth_time"].(float64)
	if !ok {
		t.Fatalf("auth_time missing: %v", claims)
	}
	if authTime := time.Unix(int64(at), 0); authTime.Before(notBefore.Truncate(time.Second)) || authTime.After(notAfter) {
		t.Errorf("auth_time = %s, want between %s and %s", authTime, notBefore, notAfter)
	}
}

func TestUserAuthClaims(t *testing.T) {
	ts := newFlowServer(t)

	cases := []struct {
		user string
		acr  string
		amr  []string
	}{
		{user: "alice", acr: "", amr: []string{"pwd"}},         // no auth configured
		{user: "bob", acr: "phr", amr: []string{"hwk", "pwd"}}, // auth: {acr: phr, amr: [hwk, pwd]}
	}
	for _, tc := range cases {
		t.Run(tc.user, func(t *testing.T) {
			before := time.Now()
			tok := exchange(t, ts, authorizeUser(t, ts, "web", "openid offline_access", tc.user, nil))
			checkAuthClaims(t, tok.IDToken, tc.acr, tc.amr, before, time.Now())

			// Refreshed ID tokens, across rotations, keep the acr, amr and
			// auth_time of the login.
			login := jwtClaims(t, tok.IDToken)
			refreshToken := tok.RefreshToken
			for i := 1; i <= 2; i++ {
				status, body := postForm(t, ts, "/oauth/token", url.Values{
					"grant_type": {"refresh_token"}, "refresh_token": {refreshToken},
				})
				refreshed := decodeTokens(t, status, body)
				if refreshed.IDToken == "" {
					t.Fatalf("refresh %d: no id_token", i)
				}
				claims := jwtClaims(t, refreshed.IDToken)
				for _, c := range []string{"acr", "amr", "auth_time"} {
					if got, want := claims[c], login[c]; !reflect.DeepEqual(got, want) {
						t.Errorf("refresh %d: %s = %v, want %v", i, c, got, want)
					}
				}
				refreshToken = refreshed.RefreshToken
			}
		})
	}
}

// TestPromptLoginAndMaxAge: every authorization request shows the login
// page, so prompt=login and max_age always get a fresh login.
func TestPromptLoginAndMaxAge(t *testing.T) {
	ts := newFlowServer(t)

	for name, extra := range map[string]url.Values{
		"prompt=login": {"prompt": {"login"}},
		"max_age=0":    {"max_age": {"0"}},
	} {
		t.Run(name, func(t *testing.T) {
			// A first login, then one that must not reuse it.
			_ = exchange(t, ts, authorizeUser(t, ts, "web", "openid", "bob", nil))
			time.Sleep(1100 * time.Millisecond) // auth_time has second precision
			before := time.Now()
			tok := exchange(t, ts, authorizeUser(t, ts, "web", "openid", "bob", extra))
			checkAuthClaims(t, tok.IDToken, "phr", []string{"hwk", "pwd"}, before, time.Now())
		})
	}
}
