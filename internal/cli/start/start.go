package start

import (
	"os"
	"os/signal"

	"github.com/spf13/cobra"
)

var (
	flagPath string
)

func Register(parent *cobra.Command) {

	var cmd = &cobra.Command{
		Use:   "start",
		Short: "Start the IDP server",
		// Long:  `Start the IDP server based on the configuration path`,
		RunE: func(cmd *cobra.Command, args []string) error {
			_, stop := signal.NotifyContext(cmd.Context(), os.Interrupt)
			defer stop()

			return nil
		},
	}

	cmd.Flags().StringVar(&flagPath, "dir", "", "Directory override, otherwise find closest '.idpzero' directory")

	parent.AddCommand(cmd)
}
