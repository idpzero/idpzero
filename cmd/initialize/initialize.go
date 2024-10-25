package initialize

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/google/uuid"
	"github.com/idpzero/idpzero/cmd/shared"
	"github.com/idpzero/idpzero/pkg/configuration"
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

		if conf.Initialized() {
			color.Red("Configuration already initialized in '%s'", conf.DirectoryPath())
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

		fmt.Printf("Initializing new configuration directory '%s'\n", conf.DirectoryPath())
		fmt.Println()

		cfg.Server.Keys = append(cfg.Server.Keys, *signingKey)
		cfg.Clients = []configuration.ClientConfig{}

		if conf.Save(&cfg); err != nil {
			return err
		}

		conf, err = configuration.Resolve(*shared.Location)

		if err != nil {
			return err
		}

		conf.PrintStatus()

		color.Green("Configuration initialized OK.")
		fmt.Println()

		return nil
	},
}
