package configuration

import (
	"os"
	"path"
	"path/filepath"
	"testing"
	"time"
)

const watchConfigA = `
server:
  port: 4379
  keyphrase: test
`

const watchConfigB = `
server:
  port: 5555
  keyphrase: test
`

func writeConfig(t *testing.T, file, content string) {
	t.Helper()
	if err := os.WriteFile(file, []byte(content), 0644); err != nil {
		t.Fatalf("write config: %v", err)
	}
}

// atomicReplace writes to a temp file in the same directory then renames it over
// the target — the pattern used by editors and our own SaveConfiguration that
// orphaned the previous file-level watch.
func atomicReplace(t *testing.T, file, content string) {
	t.Helper()
	tmp := file + ".tmp"
	if err := os.WriteFile(tmp, []byte(content), 0644); err != nil {
		t.Fatalf("write temp: %v", err)
	}
	if err := os.Rename(tmp, file); err != nil {
		t.Fatalf("rename: %v", err)
	}
}

func waitForReload(t *testing.T, ch <-chan *ServerConfig) *ServerConfig {
	t.Helper()
	select {
	case cfg := <-ch:
		return cfg
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for configuration reload")
		return nil
	}
}

func TestWatchReloadsOnInPlaceWrite(t *testing.T) {
	dir := t.TempDir()
	file := path.Join(dir, configurationFilename)
	writeConfig(t, file, watchConfigA)

	cm, err := NewConfigurationManager(dir)
	if err != nil {
		t.Fatalf("NewConfigurationManager: %v", err)
	}
	t.Cleanup(cm.Close)

	changes := make(chan *ServerConfig, 1)
	cm.OnServerChanged(func(c *ServerConfig) { changes <- c })

	writeConfig(t, file, watchConfigB)

	cfg := waitForReload(t, changes)
	if cfg.Server.Port != 5555 {
		t.Fatalf("expected reloaded port 5555, got %d", cfg.Server.Port)
	}
}

func TestWatchReloadsOnAtomicReplace(t *testing.T) {
	dir := t.TempDir()
	file := path.Join(dir, configurationFilename)
	writeConfig(t, file, watchConfigA)

	cm, err := NewConfigurationManager(dir)
	if err != nil {
		t.Fatalf("NewConfigurationManager: %v", err)
	}
	t.Cleanup(cm.Close)

	changes := make(chan *ServerConfig, 1)
	cm.OnServerChanged(func(c *ServerConfig) { changes <- c })

	// This is the case the previous file-level watch missed.
	atomicReplace(t, file, watchConfigB)

	cfg := waitForReload(t, changes)
	if cfg.Server.Port != 5555 {
		t.Fatalf("expected reloaded port 5555 after atomic replace, got %d", cfg.Server.Port)
	}
}

func TestWatchIgnoresOtherFiles(t *testing.T) {
	dir := t.TempDir()
	file := path.Join(dir, configurationFilename)
	writeConfig(t, file, watchConfigA)

	cm, err := NewConfigurationManager(dir)
	if err != nil {
		t.Fatalf("NewConfigurationManager: %v", err)
	}
	t.Cleanup(cm.Close)

	changes := make(chan *ServerConfig, 1)
	cm.OnServerChanged(func(c *ServerConfig) { changes <- c })

	// Writing an unrelated file in the watched directory must not trigger a reload.
	writeConfig(t, filepath.Join(dir, "unrelated.txt"), "noise")

	select {
	case <-changes:
		t.Fatal("reload triggered by unrelated file change")
	case <-time.After(500 * time.Millisecond):
		// expected: no reload
	}
}
