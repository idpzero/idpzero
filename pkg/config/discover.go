package config

import (
	"errors"
	"os"
	"path/filepath"
)

var (
	ErrConfigNotFound = errors.New("configuration not found")
)

const (
	folder         string = ".idpzero"
	configFilename string = "config.yaml"
	stateFilename  string = "idpzero.db"
)

type ConfigurationInfo struct {
	ConfigurationFile string
	StateFile         string
}

func Discover() (ConfigurationInfo, error) {

	// set the path to the current working directory if not provided
	pwd, err := os.Getwd()
	if err != nil {
		return ConfigurationInfo{}, err
	}

	currentPath := pwd
	for {
		if info, err := os.Stat(filepath.Join(currentPath, folder)); !os.IsNotExist(err) {
			if info.IsDir() {
				return ConfigurationInfo{
					ConfigurationFile: filepath.Join(currentPath, folder, configFilename),
					StateFile:         filepath.Join(currentPath, folder, stateFilename),
				}, nil
			}

		}
		parentPath := filepath.Dir(currentPath)
		if parentPath == currentPath {
			break
		}
		currentPath = parentPath
	}
	return ConfigurationInfo{}, ErrConfigNotFound
}
