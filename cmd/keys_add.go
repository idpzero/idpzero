package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/idpzero/idpzero/configuration"
	"github.com/spf13/cobra"
)

var addKeyCmd = &cobra.Command{
	Use:   "add",
	Short: "Add new key to configuration file",
	Long:  `Generate and append a new key to the configuration`,
	RunE: func(cmd *cobra.Command, args []string) error {

		// get the config dir to use from the path or discovery
		conf, err := configuration.Resolve(*location)

		if err != nil {
			return err
		}

		ensureInitialized(conf)

		cfg, err := conf.Load()

		if err != nil {
			color.Red("Failed to load configuration from '%s'", conf.Config().Path())
			return err
		}

		if configuration.KeyExists(cfg, *kid) && !*replace {
			color.Red("Key with ID '%s' already exists. Use --replace to force replacement.", *kid)
			fmt.Println()
			os.Exit(1)
		}

		// generate new RSA key aligned to IDP needs
		key, err := configuration.NewRSAKey(*kid, *use)

		if err != nil {
			return err
		}

		replaced := configuration.SetKey(cfg, *key, *replace)

		if replaced {
			fmt.Printf("Replaced existing key '%s' in configuration\n", *kid)
		} else {
			fmt.Printf("Added new key '%s' to configuration\n", *kid)
		}

		if conf.Save(cfg); err != nil {
			color.Red("Failed to save configuration")
			return err
		}

		fmt.Println()
		color.Green("Configuration saved OK.")

		return nil
	},
}
