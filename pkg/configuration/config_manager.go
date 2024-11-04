package configuration

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/fsnotify/fsnotify"
	"github.com/idpzero/idpzero/pkg/dbg"
	"github.com/idpzero/idpzero/pkg/validation"
	"gopkg.in/yaml.v2"
)

const (
	dirName        string = ".idpzero"
	serverFilename string = "server.yaml"
	keysFilename   string = "keys.yaml"
)

type ConfigurationManager struct {
	keysPath   string
	dirPath    string
	configPath string

	w    *fsnotify.Watcher
	done chan struct{}

	changed []func(x *ServerConfig)
}

func NewConfigurationManager(serverDirectory string, keysDirectory string) (*ConfigurationManager, error) {
	wtch, err := fsnotify.NewWatcher()

	if err != nil {
		return nil, err
	}

	cm := ConfigurationManager{
		dirPath:    serverDirectory,
		configPath: path.Join(serverDirectory, serverFilename),
		keysPath:   path.Join(keysDirectory, keysFilename),
		w:          wtch,
		done:       make(chan struct{}),
		changed:    make([]func(x *ServerConfig), 0),
	}

	// add the watcher.
	wtch.Add(cm.configPath)

	// start the watcher
	go watcher(&cm)

	return &cm, nil
}

func dirExists(directory string) (bool, error) {
	_, err := os.Stat(directory)
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

func (c *ConfigurationManager) Initialized() error {
	ok, err := dirExists(c.dirPath)

	if err != nil {
		return err
	}

	if !ok {
		return errors.New("configuration directory does not exist")
	}

	ok, err = fileExists(c.configPath)

	if err != nil {
		return err
	}

	if !ok {
		return errors.New("server configuration does not exist")
	}

	return nil
}

// PrintChecks prints the existance of each part of the configuration to the console
func (cfg *ConfigurationManager) PrintStatus() {

	dirMsg := "Configuration Directory Exists"
	confMsg := "Configuration File Exists"

	if *dbg.Debug {
		dirMsg = fmt.Sprintf("Configuration Directory Exists (%s)", cfg.dirPath)
		confMsg = fmt.Sprintf("Configuration File Exists (%s)", cfg.configPath)
	}

	consoleVal := validation.NewValidationSet()
	checklist := validation.NewChecklist("Configuration checks:")
	checklist.Add(validation.NewChecklistItem(dbg.MustOrFalse(dirExists(cfg.dirPath)), dirMsg))
	checklist.Add(validation.NewChecklistItem(dbg.MustOrFalse(fileExists(cfg.configPath)), confMsg))
	consoleVal.AddChild(checklist)
	consoleVal.Render()
}

func (r *ConfigurationManager) SaveServer(config ServerConfig) error {
	return marshal(r.configPath, config)
}

func (r *ConfigurationManager) OnServerChanged(changed func(x *ServerConfig)) {
	r.changed = append(r.changed, changed)
}

func (r *ConfigurationManager) SaveKeys(config KeysConfiguration) error {
	return marshal(r.keysPath, config)
}

func (r *ConfigurationManager) LoadServer() (*ServerConfig, error) {
	file, err := os.Open(r.configPath)
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

func (r *ConfigurationManager) LoadKeys() (*KeysConfiguration, error) {
	file, err := os.Open(r.configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	doc := &KeysConfiguration{}
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

func marshal[T ServerConfig | KeysConfiguration](path string, config T) error {
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

		if !fi.IsDir() {
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
				if event.Name == cm.configPath {
					t, err := cm.LoadServer()
					if err != nil {
						color.Red("Error loading config file from watch")
					}
					if t != nil {
						for _, changed := range cm.changed {
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
