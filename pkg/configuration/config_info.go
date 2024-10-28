package configuration

import (
	"fmt"
	"os"

	"github.com/idpzero/idpzero/pkg/dbg"
	"gopkg.in/yaml.v2"
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

	dirMsg := "Configuration Directory Exists"
	confMsg := "Configuration File Exists"

	if *dbg.Debug {
		dirMsg = fmt.Sprintf("Configuration Directory Exists (%s)", cfg.dirPath)
		confMsg = fmt.Sprintf("Configuration File Exists (%s)", cfg.configPath)
	}

	fmt.Println("Configuration checks:")
	PrintCheck(cfg.dirExists, dirMsg)
	PrintCheck(cfg.configExists, confMsg)

	fmt.Println()

}

func (r *ConfigInformation) Save(config *IDPConfiguration) error {

	data, err := yaml.Marshal(*config)

	if err != nil {
		return err
	}

	// make sure the directory exists before writing the file
	if err := ensureDirectory(r.dirPath); err != nil {
		return err
	}

	r.dirExists = true

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
