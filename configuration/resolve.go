package configuration

import (
	"errors"
	"fmt"
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

type ConfigurationCheck struct {
	exists   bool
	location string
}

func (r ConfigurationCheck) Exists() bool {
	return r.exists
}

func (r ConfigurationCheck) Path() string {
	return r.location
}

type DirectoryCheck struct {
	exists   bool
	location string
}

func (r DirectoryCheck) Exists() bool {
	return r.exists
}

func (r DirectoryCheck) Path() string {
	return r.location
}

type ConfigInformation struct {
	dir    *DirectoryCheck
	config *ConfigurationCheck
}

func (c *ConfigInformation) Directory() *DirectoryCheck {
	return c.dir
}

func (c *ConfigInformation) Config() *ConfigurationCheck {
	return c.config
}

func (c *ConfigInformation) Initialized() bool {
	return c.dir.exists && c.config.exists
}

// PrintChecks prints the existance of each part of the configuration to the console
func (cfg *ConfigInformation) PrintStatus() {

	fmt.Println("Verifying IDP configuration...")
	printCheck(cfg.Directory().Exists(), fmt.Sprintf("Configuration Directory Exists (%s)", cfg.Directory().Path()))
	printCheck(cfg.Config().Exists(), fmt.Sprintf("Configuration File Exists (%s)", configFilename))

	fmt.Println()

}

func (r *ConfigInformation) Save(config *IDPConfiguration) error {

	data, err := yaml.Marshal(*config)

	if err != nil {
		return err
	}

	if !r.dir.exists {
		if err := os.Mkdir(r.dir.location, 0755); err != nil {
			return err
		}
		r.dir.exists = true
	}

	return os.WriteFile(r.config.location, data, 0644)
}

func (r *ConfigInformation) Load() (*IDPConfiguration, error) {
	file, err := os.Open(r.config.location)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return Parse(file)
}

func Resolve(path string) (*ConfigInformation, error) {

	configDir, err := resolveDirectory(path)

	if err != nil {
		return nil, err
	}

	dirCheck := &DirectoryCheck{true, configDir}

	_, err = os.Stat(filepath.Join(path, configFilename))

	if err != nil {
		if os.IsNotExist(err) {
			dirCheck.exists = false
		} else {
			return nil, err
		}
	}

	configCheck := &ConfigurationCheck{true, filepath.Join(configDir, configFilename)}
	_, err = os.Stat(filepath.Join(path, configFilename))

	if err != nil {
		if os.IsNotExist(err) {
			configCheck.exists = false
		} else {
			return nil, err
		}
	}

	return &ConfigInformation{dirCheck, configCheck}, nil
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
