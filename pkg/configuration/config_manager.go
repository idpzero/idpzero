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
	dirName        string = ".idpzero"
	serverFilename string = "server.yaml"
	keysFilename   string = "keys.yaml"
)

type ConfigurationManager struct {
	keysDirectory string
	keysPath      string

	dirPath    string
	configPath string

	w    *fsnotify.Watcher
	done chan struct{}

	serverChanged []func(x *ServerConfig)
	keysChanged   []func(x *KeysConfiguration)
}

func NewConfigurationManager(serverDirectory string, keysDirectory string) (*ConfigurationManager, error) {
	wtch, err := fsnotify.NewWatcher()

	if err != nil {
		return nil, err
	}

	cm := ConfigurationManager{
		keysDirectory: keysDirectory,
		dirPath:       serverDirectory,
		configPath:    path.Join(serverDirectory, serverFilename),
		keysPath:      path.Join(keysDirectory, keysFilename),
		w:             wtch,
		done:          make(chan struct{}),
		serverChanged: make([]func(x *ServerConfig), 0),
		keysChanged:   make([]func(x *KeysConfiguration), 0),
	}

	// add the watcher.
	wtch.Add(cm.configPath)

	// start the watcher
	go watcher(&cm)

	return &cm, nil
}

// func dirExists(directory string) (bool, error) {
// 	_, err := os.Stat(directory)
// 	if err != nil {
// 		if os.IsNotExist(err) {
// 			return false, nil
// 		} else {
// 			return false, err
// 		}
// 	} else {
// 		return true, nil
// 	}
// }

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

func (c *ConfigurationManager) IsInitialized() (bool, error) {
	si, err := c.IsServerInitialized()

	if err != nil {
		return false, err
	}

	if si {
		ki, err := c.IsKeysInitialized()

		if err != nil {
			return false, err
		}

		return ki, nil
	}

	return false, nil
}

func (c *ConfigurationManager) IsServerInitialized() (bool, error) {
	return fileExists(c.configPath)
}

func (c *ConfigurationManager) IsKeysInitialized() (bool, error) {
	return fileExists(c.keysPath)
}

func (r *ConfigurationManager) SaveServer(config ServerConfig) error {
	return marshal(r.configPath, config)
}

func (r *ConfigurationManager) OnServerChanged(changed func(x *ServerConfig)) {
	r.serverChanged = append(r.serverChanged, changed)
}

func (r *ConfigurationManager) OnKeysChanged(changed func(x *KeysConfiguration)) {
	r.keysChanged = append(r.keysChanged, changed)
}

func (r *ConfigurationManager) SaveKeys(config KeysConfiguration) error {
	return marshal(r.keysPath, config)
}

func (r *ConfigurationManager) GetServerPath() string {
	return r.configPath
}

func (r *ConfigurationManager) GetKeysPath() string {
	return r.keysPath
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
				if event.Name == cm.configPath {
					t, err := cm.LoadServer()
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
