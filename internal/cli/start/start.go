package start

import (
	"log/slog"
	"os"
	"os/signal"

	"github.com/idpzero/idpzero/internal/discovery"
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
			ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt)
			defer stop()

			logger := slog.New(
				slog.NewTextHandler(cmd.OutOrStdout(), &slog.HandlerOptions{
					//AddSource: true,
					Level: slog.LevelDebug,
				}),
			)

			return run(ctx, logger)

		},
	}

	cmd.Flags().StringVar(&flagPath, "dir", "", "Directory override, otherwise find closest '.idpzero' directory")

	parent.AddCommand(cmd)
}

func getConfigInfo(path string) (*discovery.ConfigurationInfo, error) {
	if path == "" {
		return discovery.Discover()
	}
	return discovery.Ensure(path)
}
