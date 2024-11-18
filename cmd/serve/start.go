package serve

import (
	"database/sql"
	"fmt"
	"os"
	"os/signal"

	"github.com/fatih/color"
	"github.com/idpzero/idpzero/cmd/shared"
	"github.com/idpzero/idpzero/pkg/configuration"
	"github.com/idpzero/idpzero/pkg/dbg"
	"github.com/idpzero/idpzero/pkg/server"
	"github.com/idpzero/idpzero/pkg/store"
	"github.com/idpzero/idpzero/pkg/store/query"
	"github.com/spf13/cobra"

	_ "modernc.org/sqlite"
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

		if initialized, err := conf.IsInitialized(); err != nil {
			return err
		} else if !initialized {
			color.Yellow("Configuration not valid. Run 'idpzero init' to initialize")
			fmt.Println()
			os.Exit(1)
		}

		db, err := sql.Open("sqlite", conf.GetStateDatabasePath())
		if err != nil {
			return err
		}

		defer db.Close()

		// migrate the database to the latest version in memory
		if err = store.Migrate(db); err != nil {
			return err
		}

		// access to the data layer.
		qry := query.New(db)

		s, err := server.NewServer(dbg.Logger, conf, qry)

		if err != nil {
			return err
		}

		return s.Run(ctx)
	},
}
