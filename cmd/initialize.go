package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/google/uuid"
	"github.com/idpzero/idpzero/configuration"
	"github.com/spf13/cobra"
)

var initializeCmd = &cobra.Command{
	Use:   "init",
	Short: "Initalize configuration for idpzero",
	Long:  `Setup the configuration and data directory for idpzero`,
	RunE: func(cmd *cobra.Command, args []string) error {

		if *location == "" {
			defaultDir, err := configuration.DefaultDirectory()

			if err != nil {
				return err
			}

			*location = defaultDir
		}

		// get the config dir to use from the path or discovery
		conf, err := configuration.Resolve(*location)

		if err != nil {
			return err
		}

		if conf.Initialized() {
			color.Red("Configuration already initialized in '%s'", conf.Directory().Path())
			fmt.Println()
			os.Exit(1)
		}

		cfg := configuration.IDPConfiguration{}
		cfg.Server = configuration.ServerConfig{}
		cfg.Server.Port = 4379
		cfg.Server.KeyPhrase = uuid.New().String()

		signingKey, err := configuration.NewRSAKey("signing-key", "sig")
		if err != nil {
			return err
		}

		fmt.Printf("Initializing new configuration directory '%s'\n", conf.Directory().Path())
		fmt.Println()

		cfg.Server.Keys = append(cfg.Server.Keys, *signingKey)
		cfg.Clients = []configuration.ClientConfig{}

		if conf.Save(&cfg); err != nil {
			return err
		}

		conf, err = configuration.Resolve(*location)

		if err != nil {
			return err
		}

		conf.PrintStatus()

		color.Green("Configuration initialized OK.")
		fmt.Println()

		return nil
	},
}
