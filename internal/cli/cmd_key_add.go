package cli

import (
	"fmt"
	"os"

	"github.com/idpzero/idpzero/internal/config"
	"github.com/idpzero/idpzero/internal/style"
	"github.com/spf13/cobra"
)

var addKeyCmd = &cobra.Command{
	Use:   "add",
	Short: "Add new key to configuration file",
	Long:  `Generate and append a new key to the configuration`,
	RunE: func(cmd *cobra.Command, args []string) error {

		// get the config dir to use from the path or discovery
		conf, err := config.Resolve(*location)

		if err != nil {
			return err
		}

		configDebug(conf)

		if !conf.Initialized() {
			fmt.Println(style.WarningTextStyle.Render("Configuration not valid. Run 'idpzero init' to initialize configuration"))
			fmt.Println()
			os.Exit(1)
		}

		return nil
	},
}
