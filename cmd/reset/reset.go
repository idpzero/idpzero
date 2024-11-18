package reset

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/idpzero/idpzero/cmd/shared"
	"github.com/idpzero/idpzero/pkg/configuration"
	"github.com/idpzero/idpzero/pkg/console"
	"github.com/idpzero/idpzero/pkg/store/query"
	"github.com/spf13/cobra"
)

var (
	hard *bool = new(bool)
)

func New() *cobra.Command {
	return initializeCmd
}

func init() {
	initializeCmd.Flags().BoolVar(hard, "hard", false, "hard reset will delelte the state file")
}

var initializeCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset the state of the IDP server",
	Long:  `Reset the state of the IDP server to a clean state`,
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

		path := conf.GetStateDatabasePath()

		exists, err := fileExists(path)

		if err != nil {
			return err
		}

		if !exists {
			console.PrintCheck(console.IconDash, "State file not found, skipping reset.")
		} else {
			if hard != nil && *hard {
				if err := os.Remove(path); err != nil {
					return err
				}
				console.PrintCheck(console.IconCheck, "Hard reset completed, state file deleted OK")
			} else {
				db, err := sql.Open("sqlite", path)
				if err != nil {
					return err
				}
				defer db.Close()

				queries := query.New(db)
				if err := queries.DeleteAllAuthRequests(cmd.Context()); err != nil {
					return err
				}

				if err := queries.DeleteAllTokens(cmd.Context()); err != nil {
					return err
				}

				if err := queries.DeleteAllKeys(cmd.Context()); err != nil {
					return err
				}

				console.PrintCheck(console.IconCheck, "State reset completed OK")
			}

		}

		fmt.Println()

		return nil
	},
}

func fileExists(file string) (bool, error) {
	_, err := os.Stat(file)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, err
		}
	} else {
		return true, nil
	}
}
