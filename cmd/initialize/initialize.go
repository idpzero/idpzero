package initialize

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/google/uuid"
	"github.com/idpzero/idpzero/cmd/shared"
	"github.com/idpzero/idpzero/pkg/configuration"
	"github.com/idpzero/idpzero/pkg/console"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var (
	includeSampleConfig *bool = new(bool)
	//go:embed "samples"
	samples embed.FS
)

func New() *cobra.Command {
	return initializeCmd
}

func init() {
	initializeCmd.Flags().BoolVar(includeSampleConfig, "with-sample-config", false, "include sample configuration for clients")
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

		fmt.Printf("Initializing configuration.")
		fmt.Println()

		if ok, err := conf.IsInitialized(); err != nil {
			return err
		} else if ok {
			console.PrintCheck(console.IconCheck, "Existing configuration found. Skipping.")
		} else {
			cfg := configuration.ServerConfig{}
			cfg.Server = configuration.HostConfig{}
			cfg.Server.Port = 4379
			cfg.Server.KeyPhrase = uuid.New().String()

			cfg.Users = []*configuration.User{}
			cfg.Clients = []*configuration.ClientConfig{}

			if *includeSampleConfig {
				sampleConfig, err := getSampleConfig()

				if err != nil {
					return err
				}

				// we just take the users
				cfg.Users = sampleConfig.Users
				cfg.Clients = sampleConfig.Clients
			}

			if err := conf.SaveConfiguration(cfg); err != nil {
				return err
			}

			console.PrintCheck(console.IconCheck, "Default configuration initialized.")

			if *includeSampleConfig {
				console.PrintCheck(console.IconCheck, "Sample configuration applied.")
			}

			fmt.Println()
			fmt.Println(
				color.YellowString("Update your"),
				color.MagentaString(".gitignore"),
				color.YellowString("to include"),
				color.MagentaString(".idpzero/cache"),
				color.YellowString("as this directory should not be added to source control."),
			)
		}

		fmt.Println()

		fmt.Println(
			"To start the server run",
			color.CyanString(fmt.Sprintf("%s serve", filepath.Base(os.Args[0]))),
		)
		fmt.Println()

		return nil
	},
}

func getSampleConfig() (configuration.ServerConfig, error) {
	config := configuration.ServerConfig{}

	// load the sample users
	userData, err := samples.ReadFile("samples/sample.yaml")

	if err != nil {
		return config, err
	}

	if err := yaml.Unmarshal(userData, &config); err != nil {
		return config, err
	}

	return config, nil
}
