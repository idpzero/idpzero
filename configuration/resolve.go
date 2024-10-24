package configuration

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

var (
	ErrDiscoveryFailed = errors.New("no configuration directory found")
)

const (
	directory      string = ".idpzero"
	configFilename string = "config.yaml"
)

type ConfigInformation struct {
	dirPath   string
	dirExists bool

	configPath   string
	configExists bool
}

func (c *ConfigInformation) DirectoryPath() string {
	return c.dirPath
}

func (c *ConfigInformation) ConfigFilePath() string {
	return c.configPath
}

func (c *ConfigInformation) Initialized() bool {
	return c.dirExists && c.configExists
}

// PrintChecks prints the existance of each part of the configuration to the console
func (cfg *ConfigInformation) PrintStatus() {

	fmt.Println("IDP configuration checks:")
	printCheck(cfg.dirExists, fmt.Sprintf("Configuration Directory Exists (%s)", cfg.dirPath))
	printCheck(cfg.configExists, fmt.Sprintf("Configuration File Exists (%s)", cfg.configPath))

	fmt.Println()

}

func (r *ConfigInformation) Save(config *IDPConfiguration) error {

	data, err := yaml.Marshal(*config)

	if err != nil {
		return err
	}

	if !r.dirExists {
		if err := os.Mkdir(r.dirPath, 0755); err != nil {
			return err
		}
		r.dirExists = true
	}

	return os.WriteFile(r.configPath, data, 0644)
}

func (r *ConfigInformation) Load() (*IDPConfiguration, error) {
	file, err := os.Open(r.configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return parse(file)
}

func Resolve(path string) (*ConfigInformation, error) {

	ci := &ConfigInformation{}

	configDir, err := resolveDirectory(path)

	if err != nil {
		return nil, err
	}

	ci.dirPath = configDir

	_, err = os.Stat(ci.dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			ci.dirExists = false
		} else {
			return nil, err
		}
	} else {
		ci.dirExists = true
	}

	ci.configPath = filepath.Join(configDir, configFilename)
	_, err = os.Stat(ci.configPath)

	if err != nil {
		if os.IsNotExist(err) {
			ci.configExists = false
		} else {
			return nil, err
		}
	} else {
		ci.configExists = true
	}

	return ci, nil
}

func DefaultDirectory() (string, error) {
	cwd, err := os.Getwd()

	if err != nil {
		return "", err
	}

	return filepath.Join(cwd, directory), nil
}

func resolveDirectory(path string) (string, error) {

	cwd, err := os.Getwd()

	if err != nil {
		return "", err
	}

	if path == "" {
		return discoverConfigDir(cwd)
	} else {
		if filepath.IsLocal(path) {
			return filepath.Join(cwd, path), nil
		} else {
			// assume absolute
			return path, nil
		}
	}
}

// func isInitialized(path string) (bool, error) {

// 	// check if config file exists
// 	if _, err := os.Stat(filepath.Join(path, configFilename)); os.IsNotExist(err) {
// 		return false, nil
// 	} else if err != nil {
// 		return false, err
// 	}

// 	return true, nil
// }

func discoverConfigDir(cwd string) (string, error) {

	currentPath := cwd
	for {
		if info, err := os.Stat(filepath.Join(currentPath, directory)); !os.IsNotExist(err) {
			if info.IsDir() {
				return filepath.Join(currentPath, directory), nil
			}
		}
		parentPath := filepath.Dir(currentPath)
		if parentPath == currentPath {
			break
		}
		currentPath = parentPath
	}

	return "", ErrDiscoveryFailed
}

func parse(reader io.Reader) (*IDPConfiguration, error) {
	doc := &IDPConfiguration{}
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(reader)

	if err != nil {
		return nil, err
	}

	if yaml.Unmarshal(buf.Bytes(), doc); err != nil {
		return nil, err
	}

	return doc, nil
}

// // EnsureDirectory checks if the directory exists at the path provided and creates it if it doesn't.
// func ensureDirectory(path string) error {
// 	if !filepath.IsAbs(path) {
// 		return errors.New("path must be absolute")
// 	}

// 	if _, err := os.Stat(path); os.IsNotExist(err) {
// 		// create the directory
// 		if err := os.Mkdir(path, 0755); err != nil {
// 			return err
// 		}
// 	} else if err != nil {
// 		return err
// 	}

// 	return nil
// }
