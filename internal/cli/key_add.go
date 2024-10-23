package cli

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/idpzero/idpzero/internal/config"
	"github.com/spf13/cobra"
)

func ensureInitialized(conf *config.ConfigInformation) {

	config.PrintChecks(conf)

	if !conf.Initialized() {
		color.Yellow("Configuration not valid. Run 'idpzero init' to initialize configuration")
		fmt.Println()
		os.Exit(1)
	}
}

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

		ensureInitialized(conf)

		cfg := &config.IDPConfiguration{}
		if config.LoadFromFile(cfg, conf.Config().Path()); err != nil {
			return err
		}

		if config.KeyExists(cfg, *kid) && !*replace {
			color.Red("Key with ID '%s' already exists. Use --replace to force replacement.", *kid)
			fmt.Println()
			os.Exit(1)
		}

		// generate new RSA key aligned to IDP needs
		key, err := config.NewRSAKey(*kid, *use)

		if err != nil {
			return err
		}

		replaced := config.SetKey(cfg, *key, *replace)

		if replaced {
			fmt.Printf("Replaced existing key '%s' in configuration\n", *kid)
		} else {
			fmt.Printf("Added new key '%s' to configuration\n", *kid)
		}

		if config.Save(cfg, conf.Config().Path()); err != nil {
			color.Red("Failed to save configuration")
			return err
		}

		fmt.Println()
		color.Green("Configuration saved OK.")

		return nil
	},
}
