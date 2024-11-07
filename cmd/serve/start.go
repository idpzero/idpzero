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

		defer conf.Close()

		configuration.PrintStatus(conf)

		if initialized, err := conf.IsInitialized(); err != nil {
			return err
		} else if !initialized {
			color.Yellow("Configuration not valid. Run 'idpzero init' to initialize")
			fmt.Println()
			os.Exit(1)
		}

		cfg, err := conf.LoadServer()

		if err != nil {
			return err
		}

		idpStore, err := idp.NewStorage(dbg.Logger)

		if err != nil {
			return err
		}

		idpStore.SetConfig(cfg)

		s, err := server.NewServer(dbg.Logger, cfg, idpStore)

		if err != nil {
			return err
		}

		// watch for changes and set it again.
		conf.OnServerChanged(func(x *configuration.ServerConfig) {
			color.Yellow("Configuration changed. Reloading...")
			idpStore.SetConfig(x)
			s.SetConfig(x)
		})

		return s.Run(ctx)
	},
}
