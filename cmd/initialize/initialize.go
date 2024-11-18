package initialize

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/google/uuid"
	"github.com/idpzero/idpzero/cmd/shared"
	"github.com/idpzero/idpzero/pkg/configuration"
	"github.com/idpzero/idpzero/pkg/console"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	return initializeCmd
}

var initializeCmd = &cobra.Command{
	Use:   "init",
	Short: "Initalize configuration for idpzero",
	Long:  `Setup the configuration and data directory for idpzero`,
	RunE: func(cmd *cobra.Command, args []string) error {

		if *shared.Location == "" {
			defaultDir, err := configuration.DefaultDirectory()

			if err != nil {
				return err
			}

			shared.Location = &defaultDir
		}

		// get the config dir to use from the path or discovery
		conf, err := configuration.Resolve(*shared.Location)

		if err != nil {
			return err
		}

		if ok, err := conf.IsInitialized(); err != nil {
			return err
		} else if ok {
			console.PrintCheck(console.IconDash, "Server configuration already initialized. Skipping.")
		} else {
			cfg := configuration.ServerConfig{}
			cfg.Server = configuration.HostConfig{}
			cfg.Server.Port = 4379
			cfg.Server.KeyPhrase = uuid.New().String()

			fmt.Printf("Initializing new configuration directory.")
			fmt.Println()

			cfg.Clients = []configuration.ClientConfig{}

			if err := conf.SaveConfiguration(cfg); err != nil {
				return err
			}

			console.PrintCheck(console.IconCheck, "Server configuration initialized successfully.")
		}

		fmt.Println()
		color.Green("Initialized OK!")
		fmt.Println()

		return nil
	},
}
