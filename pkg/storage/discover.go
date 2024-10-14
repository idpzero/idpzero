package storage

import (
	"os"
	"path/filepath"
)

const (
	folder   string = ".idpzero"
	filename string = "config.yaml"
)

func DiscoverConfigFile() (string, error) {

	// set the path to the current working directory if not provided
	pwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	currentPath := pwd
	for {
		if info, err := os.Stat(filepath.Join(currentPath, folder)); !os.IsNotExist(err) {
			if info.IsDir() {
				return filepath.Join(currentPath, folder, filename), nil
			}

		}
		parentPath := filepath.Dir(currentPath)
		if parentPath == currentPath {
			break
		}
		currentPath = parentPath
	}
	return "", ErrConfigNotFound
}
