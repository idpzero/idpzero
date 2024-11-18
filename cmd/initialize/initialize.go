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

		if ok, err := conf.IsServerInitialized(); err != nil {
			return err
		} else if ok {
			console.PrintCheck(console.IconCheck, "Server configuration already initialized. Skipping.")
		} else {
			cfg := configuration.ServerConfig{}
			cfg.Server = configuration.HostConfig{}
			cfg.Server.Port = 4379
			cfg.Server.KeyPhrase = uuid.New().String()

			fmt.Printf("Initializing new configuration directory.")
			fmt.Println()

			cfg.Clients = []configuration.ClientConfig{}

			if err := conf.SaveServer(cfg); err != nil {
				return err
			}

			console.PrintCheck(console.IconCheck, "Server configuration initialized successfully.")
		}

		// if ok, err := conf.IsKeysInitialized(); err != nil {
		// 	return err
		// } else if ok {
		// 	console.PrintCheck(console.IconCheck, "Keys configuration already initialized. Skipping.")
		// } else {
		// 	fmt.Printf("Initializing new keys directory.")
		// 	fmt.Println()

		// 	keys := configuration.KeysConfiguration{}
		// 	keys.Keys = []configuration.Key{}

		// 	nk, err := configuration.NewRSAKey("default", "sig")

		// 	if err != nil {
		// 		return err
		// 	}

		// 	keys.Keys = append(keys.Keys, *nk)

		// 	if err := conf.SaveKeys(keys); err != nil {
		// 		return err
		// 	}

		// 	console.PrintCheck(console.IconCheck, "Keys configuration initialized successfully.")
		// }

		conf, err = configuration.Resolve(*shared.Location)

		if err != nil {
			return err
		}

		fmt.Println()
		color.Green("Initialized OK!")
		fmt.Println()

		return nil
	},
}
