package configuration

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/fatih/color"
	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v2"
)

const (
	defaultDirectoryName string = ".idpzero"
	stateDirectoryName   string = "cache"
	// filenames
	configurationFilename string = "server.yaml"
	dbFilename            string = "state.sqlite"
)

type ConfigurationManager struct {
	// storage locations
	configurationDirectory string
	// file paths for the configuration and state database
	stateDbFilePath       string
	configurationFilePath string
	// watcher configuration for the configuration file
	w             *fsnotify.Watcher
	done          chan struct{}
	serverChanged []func(x *ServerConfig)
}

func NewConfigurationManager(dir string) (*ConfigurationManager, error) {
	wtch, err := fsnotify.NewWatcher()

	if err != nil {
		return nil, err
	}

	cm := ConfigurationManager{
		configurationDirectory: dir,
		// paths to use
		configurationFilePath: path.Join(dir, configurationFilename),
		stateDbFilePath:       path.Join(dir, stateDirectoryName, dbFilename),
		// watch configuration
		w:             wtch,
		done:          make(chan struct{}),
		serverChanged: make([]func(x *ServerConfig), 0),
	}

	if err := ensureDirectory(path.Dir(cm.configurationFilePath)); err != nil {
		return nil, err
	}
	if err := ensureDirectory(path.Dir(cm.stateDbFilePath)); err != nil {
		return nil, err
	}

	// Watch the configuration *directory* rather than the file directly. Many
	// editors (and our own atomic writes) save by writing a temp file and renaming
	// it over the target, which replaces the inode and silently orphans a
	// file-level watch. Watching the directory keeps working across those
	// replacements; the watcher filters events down to the config file.
	wtch.Add(cm.configurationDirectory)

	// start the watcher
	go watcher(&cm)

	return &cm, nil
}

func (r *ConfigurationManager) IsInitialized() (bool, error) {
	return fileExists(r.configurationFilePath)
}

func (r *ConfigurationManager) SaveConfiguration(config ServerConfig) error {
	return marshal(r.configurationFilePath, config)
}

func (r *ConfigurationManager) OnServerChanged(changed func(x *ServerConfig)) {
	r.serverChanged = append(r.serverChanged, changed)
}

func (r *ConfigurationManager) GetConfigurationFilePath() string {
	return r.configurationFilePath
}

func (r *ConfigurationManager) GetStateDatabasePath() string {
	return r.stateDbFilePath
}

func (r *ConfigurationManager) LoadConfiguration() (*ServerConfig, error) {
	file, err := os.Open(r.configurationFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	doc := &ServerConfig{}
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(file)

	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(buf.Bytes(), doc); err != nil {
		return nil, err
	}

	// process the parts that shouldnt be marshalled or processed separately.
	// Only confidential clients (client_secret_basic / client_secret_post) get a
	// derived secret. Public clients (auth_method: none) authenticate via PKCE and
	// intentionally have no secret so their config can be committed to source control.
	for _, c := range doc.Clients {
		if !c.RequiresClientSecret() {
			c.ClientSecret = ""
			continue
		}

		h := sha1.New()
		h.Write([]byte(doc.Server.KeyPhrase))
		h.Write([]byte(c.ClientID))
		c.ClientSecret = hex.EncodeToString(h.Sum(nil))
	}

	return doc, nil

}

func (w *ConfigurationManager) Close() {
	w.w.Close()
	<-w.done // wait for the go routine to finish
}

func marshal[T ServerConfig](path string, config T) error {
	data, err := yaml.Marshal(config)

	if err != nil {
		return err
	}

	// make sure the directory exists before writing the file
	if err := ensureDirectory(filepath.Dir(path)); err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// EnsureDirectory checks if the directory exists at the path provided and creates it if it doesn't.
func ensureDirectory(path string) error {

	if fi, err := os.Stat(path); os.IsNotExist(err) {

		if fi != nil && !fi.IsDir() {
			return errors.New("path exists but is not a directory")
		}

		// create the directory
		if err := os.Mkdir(path, 0755); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	return nil
}

func watcher(cm *ConfigurationManager) {

	// Start listening for events.
	defer func() {
		cm.done <- struct{}{}
	}()

	// A single save can emit several events (e.g. CREATE + WRITE, or the multiple
	// WRITEs of a non-atomic save), so debounce before reloading to avoid firing
	// change callbacks repeatedly for one logical change.
	var debounce *time.Timer
	defer func() {
		if debounce != nil {
			debounce.Stop()
		}
	}()

	reload := func() {
		color.Yellow("Server configuration changed.")

		t, err := cm.LoadConfiguration()
		if err != nil {
			color.Red("Error loading config file from watch: %v", err)
			return
		}
		for _, changed := range cm.serverChanged {
			go changed(t)
		}
	}

	for {
		select {
		case event, ok := <-cm.w.Events:
			if !ok {
				return
			}

			// We watch the directory, so filter down to the configuration file.
			if filepath.Clean(event.Name) != filepath.Clean(cm.configurationFilePath) {
				continue
			}

			// React to writes and to atomic replacements (create/rename over the
			// target); ignore pure removals.
			if !event.Has(fsnotify.Write) && !event.Has(fsnotify.Create) && !event.Has(fsnotify.Rename) {
				continue
			}

			if debounce != nil {
				debounce.Stop()
			}
			debounce = time.AfterFunc(100*time.Millisecond, reload)

		case err, ok := <-cm.w.Errors:
			if !ok {
				return
			}

			color.Red("Error occured during watch: %v", err)
		}
	}

}

func fileExists(file string) (bool, error) {
	_, err := os.Stat(file)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, err
		}
	} else {
		return true, nil
	}
}
