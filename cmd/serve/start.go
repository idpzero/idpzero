package serve

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/fatih/color"
	"github.com/idpzero/idpzero/cmd/shared"
	"github.com/idpzero/idpzero/pkg/configuration"
	"github.com/idpzero/idpzero/pkg/dbg"
	"github.com/idpzero/idpzero/pkg/idp"
	"github.com/idpzero/idpzero/pkg/server"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	return startCmd
}

var startCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the IDP server and login experience",
	// Long:  `Start the IDP server based on the configuration path`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt)
		defer stop()

		// get the config dir to use from the path or discovery
		conf, err := configuration.Resolve(*shared.Location)

		if err != nil {
			return err
		}

		conf.PrintStatus()

		if !conf.Initialized() {
			color.Yellow("Configuration not valid. Run 'idpzero init' to initialize configuration")
			fmt.Println()
			os.Exit(1)
		}

		cfg, err := conf.Load()

		if err != nil {
			return err
		}

		idpStore, err := idp.NewStorage(dbg.Logger)

		if err != nil {
			return err
		}

		valid := idp.PrintValidation(cfg)
		if !valid {
			color.Red("Configuration not valid. Fix the configuration and try again.")
			fmt.Println()
			os.Exit(1)
		}

		idpStore.SetConfig(cfg)

		// watch for changes and set it again.
		w, err := configuration.NewWatcher(conf, func(x *configuration.IDPConfiguration) {
			color.Yellow("Configuration changed. Reloading...")

			valid := idp.PrintValidation(x)
			if valid {
				idpStore.SetConfig(x)
			}
		})

		if err != nil {
			return err
		}

		defer w.Close() // wait for the tiy up.

		s, err := server.NewServer(dbg.Logger, cfg, idpStore)

		if err != nil {
			return err
		}

		return s.Run(ctx)
	},
}
