package configuration

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

var (
	ErrDiscoveryFailed = errors.New("no configuration directory found")
)

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

	if err := yaml.Unmarshal(buf.Bytes(), doc); err != nil {
		return nil, err
	}

	return doc, nil
}

// EnsureDirectory checks if the directory exists at the path provided and creates it if it doesn't.
func ensureDirectory(path string) error {
	if !filepath.IsAbs(path) {
		return errors.New("path must be absolute")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		// create the directory
		if err := os.Mkdir(path, 0755); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	return nil
}
