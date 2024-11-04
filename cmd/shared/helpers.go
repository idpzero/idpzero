package shared

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/idpzero/idpzero/pkg/configuration"
)

func EnsureInitialized(conf *configuration.ConfigurationManager) error {

	conf.PrintStatus()

	if initialized, err := conf.Initialized(); err != nil {
		return err
	} else if !initialized {
		color.Yellow("Configuration not valid. Run 'idpzero init' to initialize")
		fmt.Println()
		return ErrNotInitialized
	}

	return nil
}
