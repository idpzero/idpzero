package shared

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/idpzero/idpzero/pkg/configuration"
)

func EnsureInitialized(conf *configuration.ConfigurationManager) error {

	conf.PrintStatus()

	if err := conf.Initialized(); err != nil {
		color.Yellow("Configuration not valid. Run 'idpzero init' to initialize configuration")
		fmt.Println()
		return ErrNotInitialized
	}

	return nil
}
