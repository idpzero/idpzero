package cli

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/idpzero/idpzero/internal/config"
	"github.com/spf13/cobra"
)

var removeKeyCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a key from the configuration if it exists",
	Long:  `Generate and append a new key to the configuration`,
	RunE: func(cmd *cobra.Command, args []string) error {

		// get the config dir to use from the path or discovery
		conf, err := config.Resolve(*location)

		if err != nil {
			return err
		}

		ensureInitialized(conf)

		cfg := &config.IDPConfiguration{}
		if config.LoadFromFile(cfg, conf.Config().Path()); err != nil {
			return err
		}

		removed := config.RemoveKey(cfg, *kid)

		if removed {
			fmt.Printf("Key '%s' removed from configuration\n", *kid)

			if config.Save(cfg, conf.Config().Path()); err != nil {
				color.Red("Failed to save configuration")
				return err
			}

			fmt.Println()
			color.Green("Configuration saved OK.")

		} else {
			fmt.Printf("Key '%s' not found. No action required.\n", *kid)
			fmt.Println()
		}

		return nil
	},
}
