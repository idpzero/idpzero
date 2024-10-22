package config

import (
	"errors"
	"os"
	"path/filepath"
)

var (
	ErrDiscoveryFailed = errors.New("no configuration directory found")
)

const (
	directory      string = ".idpzero"
	configFilename string = "config.yaml"
)

type ResolutionCheck struct {
	exists   bool
	location string
}

func (r ResolutionCheck) Exists() bool {
	return r.exists
}

func (r ResolutionCheck) Path() string {
	return r.location
}

type ConfigInformation struct {
	dir    ResolutionCheck
	config ResolutionCheck
}

func (c *ConfigInformation) Directory() ResolutionCheck {
	return c.dir
}

func (c *ConfigInformation) Config() ResolutionCheck {
	return c.config
}

func Resolve(path string) (*ConfigInformation, error) {

	configDir, err := resolveDirectory(path)

	if err != nil {
		return nil, err
	}

	dirCheck := ResolutionCheck{true, configDir}

	_, err = os.Stat(filepath.Join(path, configFilename))

	if err != nil {
		if os.IsNotExist(err) {
			dirCheck.exists = false
		} else {
			return nil, err
		}
	}

	configCheck := ResolutionCheck{true, filepath.Join(configDir, configFilename)}
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
