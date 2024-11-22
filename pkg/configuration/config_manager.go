package configuration

import (
	"bytes"
	"errors"
	"os"
	"path"
	"path/filepath"

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

	// add the watcher.
	wtch.Add(cm.configurationFilePath)

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

	for {
		select {
		case event, ok := <-cm.w.Events:
			if !ok {
				return
			}
			if event.Has(fsnotify.Write) {
				if event.Name == cm.configurationFilePath {

					color.Yellow("Server configuration changed.")

					t, err := cm.LoadConfiguration()
					if err != nil {
						color.Red("Error loading config file from watch")
					}
					if t != nil {
						for _, changed := range cm.serverChanged {
							go changed(t)
						}
					}
				}
			}
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
