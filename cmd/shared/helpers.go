package shared

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/idpzero/idpzero/pkg/configuration"
)

func EnsureInitialized(conf *configuration.ConfigurationManager) error {

	if initialized, err := conf.IsInitialized(); err != nil {
		return err
	} else if !initialized {
		fmt.Println("Initialization status:")
		configuration.PrintStatus(conf)

		color.Yellow("Configuration not valid. Run 'idpzero init' to initialize")
		fmt.Println()
		return ErrNotInitialized
	}

	return nil
}
