package configuration

import (
	"os"
	"path"
	"testing"
)

const testConfigYAML = `
server:
  port: 4379
  keyphrase: test-keyphrase
clients:
  - name: Public SPA
    client_id: spa
    auth_method: none
    grant_types:
      - authorization_code
    redirect_uris:
      - http://localhost:3000/callback
    response_types:
      - code
  - name: Default Public
    client_id: default
    grant_types:
      - authorization_code
    redirect_uris:
      - http://localhost:3000/callback
    response_types:
      - code
  - name: Confidential Web
    client_id: web
    auth_method: client_secret_basic
    grant_types:
      - authorization_code
    redirect_uris:
      - http://localhost:3000/callback
    response_types:
      - code
`

func newTestManager(t *testing.T) *ConfigurationManager {
	t.Helper()

	dir := t.TempDir()
	cm, err := NewConfigurationManager(dir)
	if err != nil {
		t.Fatalf("NewConfigurationManager: %v", err)
	}
	t.Cleanup(cm.Close)

	if err := os.WriteFile(path.Join(dir, configurationFilename), []byte(testConfigYAML), 0644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	return cm
}

func TestLoadConfigurationSecretDerivation(t *testing.T) {
	cm := newTestManager(t)

	cfg, err := cm.LoadConfiguration()
	if err != nil {
		t.Fatalf("LoadConfiguration: %v", err)
	}

	byID := map[string]*ClientConfig{}
	for _, c := range cfg.Clients {
		byID[c.ClientID] = c
	}

	// Public client (auth_method: none) must have no secret so its config can be
	// committed to source control safely.
	if got := byID["spa"].ClientSecret; got != "" {
		t.Errorf("public client 'spa' should have no secret, got %q", got)
	}

	// A client with no auth_method defaults to a public (PKCE) client and also
	// gets no secret.
	if got := byID["default"].ClientSecret; got != "" {
		t.Errorf("defaulted public client 'default' should have no secret, got %q", got)
	}

	// Confidential client must have a derived, non-empty secret.
	if got := byID["web"].ClientSecret; got == "" {
		t.Error("confidential client 'web' should have a derived secret, got empty")
	}
}

func TestLoadConfigurationSecretIsStable(t *testing.T) {
	// The derived secret must be deterministic for a given keyphrase + client id
	// so it stays consistent between server restarts and dashboard views.
	cm := newTestManager(t)

	first, err := cm.LoadConfiguration()
	if err != nil {
		t.Fatalf("LoadConfiguration (first): %v", err)
	}
	second, err := cm.LoadConfiguration()
	if err != nil {
		t.Fatalf("LoadConfiguration (second): %v", err)
	}

	if first.Clients[2].ClientSecret != second.Clients[2].ClientSecret {
		t.Errorf("derived secret not stable across loads: %q vs %q",
			first.Clients[2].ClientSecret, second.Clients[2].ClientSecret)
	}
}
