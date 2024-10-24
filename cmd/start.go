package cmd

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/fatih/color"
	"github.com/idpzero/idpzero/configuration"
	"github.com/idpzero/idpzero/idp"
	"github.com/idpzero/idpzero/server"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the IDP server",
	// Long:  `Start the IDP server based on the configuration path`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt)
		defer stop()

		// get the config dir to use from the path or discovery
		conf, err := configuration.Resolve(*location)

		if err != nil {
			return err
		}

		conf.PrintStatus()

		if !conf.Initialized() {
			color.Yellow("Configuration not valid. Run 'idpzero init' to initialize configuration")
			fmt.Println()
			os.Exit(1)
		}

		cfg := configuration.IDPConfiguration{}
		cfg.Server = configuration.ServerConfig{}
		cfg.Server.Port = 4379
		cfg.Server.KeyPhrase = "secret"
		key1, err := configuration.NewRSAKey("sample", "sig")
		if err != nil {
			return err
		}

		cfg.Server.Keys = append(cfg.Server.Keys, *key1)
		cfg.Clients = []configuration.ClientConfig{}

		if conf.Save(&cfg); err != nil {
			return err
		}

		idpStore, err := idp.NewStorage(logger)
		idpStore.SetConfig(&cfg)

		if err != nil {
			return err
		}

		s, err := server.NewServer(logger, cfg, idpStore)

		if err != nil {
			return err
		}

		return s.Run(ctx)
	},
}
