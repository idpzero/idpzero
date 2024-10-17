package discovery

import (
	"errors"
	"os"
	"path/filepath"
)

var (
	ErrDiscoveryFailed = errors.New("no configuration directory found")
)

const (
	folder         string = ".idpzero"
	configFileName string = "idpzero.yaml"
)

type ConfigurationInfo struct {
	Directory string
}

func (ci *ConfigurationInfo) ConfigPath() string {
	return filepath.Join(ci.Directory, configFileName)
}

// Ensure checks if the directory exists at the path provided and creates it if it doesn't.
func Ensure(path string) (*ConfigurationInfo, error) {
	target := filepath.Join(path)
	if _, err := os.Stat(target); os.IsNotExist(err) {
		// create the directory
		if err := os.Mkdir(target, 0755); err != nil {
			return nil, err
		}
		return &ConfigurationInfo{
			Directory: target,
		}, nil
	} else if err != nil {
		return nil, err
	} else {
		return &ConfigurationInfo{
			Directory: target,
		}, nil
	}
}

func Discover() (*ConfigurationInfo, error) {

	// set the path to the current working directory if not provided
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	currentPath := pwd
	for {
		if info, err := os.Stat(filepath.Join(currentPath, folder)); !os.IsNotExist(err) {
			if info.IsDir() {
				return &ConfigurationInfo{
					Directory: filepath.Join(currentPath, folder),
				}, nil
			}

		}
		parentPath := filepath.Dir(currentPath)
		if parentPath == currentPath {
			break
		}
		currentPath = parentPath
	}

	return nil, ErrDiscoveryFailed
}
