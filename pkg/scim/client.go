package scim

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client is a minimal SCIM 2.0 client for the User resource. It targets a single
// service provider identified by its base URL (e.g. https://host/scim/v2).
type Client struct {
	baseURL string
	token   string
	http    *http.Client
}

// New builds a SCIM client for the given endpoint. token is optional; when empty
// no Authorization header is sent.
func New(endpoint, token string) *Client {
	return &Client{
		baseURL: strings.TrimRight(endpoint, "/"),
		token:   token,
		http: &http.Client{
			Timeout: 30 * time.Second,
			// SCIM is a JSON API: a redirect means the request was not accepted
			// (typically an auth gateway bouncing an unauthenticated call to a
			// login page). Surface it rather than silently following it to an
			// HTML page and then failing to parse that as JSON.
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
	}
}

// FindByExternalID looks up a user by its externalId, returning nil when no match
// exists. idpzero sets externalId to the user's Subject.
func (c *Client) FindByExternalID(ctx context.Context, externalID string) (*User, error) {
	filter := fmt.Sprintf(`externalId eq %q`, externalID)
	path := "/Users?filter=" + url.QueryEscape(filter)

	var list ListResponse
	if err := c.do(ctx, http.MethodGet, path, nil, &list); err != nil {
		return nil, err
	}

	if len(list.Resources) == 0 {
		return nil, nil
	}

	// Only trust a resource whose externalId actually matches what we asked for.
	// Some targets ignore the filter and return unrelated users; blindly using
	// Resources[0] would then replace the wrong record (or the same record for
	// every user).
	for i := range list.Resources {
		if list.Resources[i].ExternalID == externalID {
			return &list.Resources[i], nil
		}
	}

	// Results came back but none match — the endpoint is not honouring the
	// externalId filter. Fail loudly rather than act on the wrong record.
	return nil, fmt.Errorf("target returned %d user(s) for externalId %q but none actually match (id of first result: %q) — the SCIM endpoint does not appear to support filtering by externalId", len(list.Resources), externalID, list.Resources[0].ID)
}

// Create provisions a new user on the target.
func (c *Client) Create(ctx context.Context, user User) (*User, error) {
	var created User
	if err := c.do(ctx, http.MethodPost, "/Users", user, &created); err != nil {
		return nil, err
	}
	return &created, nil
}

// Replace overwrites an existing user (identified by its service-provider id).
func (c *Client) Replace(ctx context.Context, id string, user User) (*User, error) {
	var updated User
	if err := c.do(ctx, http.MethodPut, "/Users/"+url.PathEscape(id), user, &updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

// do performs a SCIM request, encoding body (when non-nil) and decoding a
// successful response into out (when non-nil). Non-2xx responses are parsed into
// a descriptive error, preferring the SCIM error envelope when present.
func (c *Client) do(ctx context.Context, method, path string, body any, out any) error {
	var reader io.Reader
	if body != nil {
		encoded, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reader = bytes.NewReader(encoded)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, reader)
	if err != nil {
		return err
	}

	req.Header.Set("Accept", ContentType)
	if body != nil {
		req.Header.Set("Content-Type", ContentType)
	}
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// A redirect (returned to us because we disabled auto-following) almost always
	// means the target rejected the request before it reached the SCIM API — most
	// commonly an auth gateway sending an unauthenticated call to a login page.
	if resp.StatusCode >= 300 && resp.StatusCode < 400 {
		loc := resp.Header.Get("Location")
		return fmt.Errorf("scim %s %s: unexpected redirect (%d) to %q — the target rejected the request (usually an invalid/expired bearer token or an endpoint URL that is not a SCIM API)", method, path, resp.StatusCode, loc)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("scim %s %s: %s", method, path, describeError(resp.StatusCode, data))
	}

	if out != nil && len(data) > 0 {
		// Guard against a 2xx HTML body (e.g. a login page served with 200):
		// decoding it as JSON would produce a confusing "invalid character '<'".
		if ct := resp.Header.Get("Content-Type"); ct != "" && !strings.Contains(ct, "json") {
			return fmt.Errorf("scim %s %s: expected a JSON response but got Content-Type %q — is %s a SCIM API endpoint? body: %s", method, path, ct, c.baseURL, snippet(data))
		}
		if err := json.Unmarshal(data, out); err != nil {
			return fmt.Errorf("scim %s %s: decoding response: %w (body: %s)", method, path, err, snippet(data))
		}
	}
	return nil
}

// snippet returns a short, single-line preview of a response body for error messages.
func snippet(data []byte) string {
	const max = 120
	s := strings.TrimSpace(string(data))
	s = strings.Join(strings.Fields(s), " ")
	if len(s) > max {
		return s[:max] + "…"
	}
	return s
}

// describeError renders a human-readable message from a SCIM error response,
// falling back to the raw body when it is not a SCIM error envelope.
func describeError(status int, body []byte) string {
	var e ErrorResponse
	if err := json.Unmarshal(body, &e); err == nil && e.Detail != "" {
		return fmt.Sprintf("%d %s", status, e.Detail)
	}
	trimmed := strings.TrimSpace(string(body))
	if trimmed == "" {
		return fmt.Sprintf("%d (no response body)", status)
	}
	return fmt.Sprintf("%d %s", status, trimmed)
}
