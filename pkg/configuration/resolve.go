package configuration

import (
	"errors"
	"os"
	"path/filepath"
)

var (
	ErrDiscoveryFailed = errors.New("no configuration directory found")
)

func Resolve(path string) (*ConfigurationManager, error) {

	configDir, err := resolveDirectory(path)

	if err != nil {
		return nil, err
	}

	return NewConfigurationManager(configDir)
}

func DefaultDirectory() (string, error) {
	cwd, err := os.Getwd()

	if err != nil {
		return "", err
	}

	return filepath.Join(cwd, defaultDirectoryName), nil
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

func discoverConfigDir(cwd string) (string, error) {

	currentPath := cwd
	for {
		if info, err := os.Stat(filepath.Join(currentPath, defaultDirectoryName)); !os.IsNotExist(err) {
			if info.IsDir() {
				return filepath.Join(currentPath, defaultDirectoryName), nil
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
