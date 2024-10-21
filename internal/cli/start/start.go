package start

import (
	"log/slog"
	"os"
	"os/signal"

	"github.com/idpzero/idpzero/internal/discovery"
	"github.com/idpzero/idpzero/internal/idp"
	"github.com/idpzero/idpzero/internal/server"
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

			config := idp.IDPConfiguration{}
			config.Server = idp.ServerConfig{}
			config.Server.Port = 4379
			config.Server.Issuer = "https://idpzero.local"
			config.Server.KeyPhrase = "secret"
			key1, err := idp.NewRSAKey()
			if err != nil {
				return err
			}
			config.Server.SigningKeys = append(config.Server.SigningKeys, *key1)
			config.Clients = []idp.ClientConfig{}

			idpStore, err := idp.NewStorage(logger)
			idpStore.SetConfig(&config)

			if err != nil {
				return err
			}

			s, err := server.NewServer(logger, config, idpStore)

			if err != nil {
				return err
			}

			return s.Run(ctx)
		},
	}

	cmd.Flags().StringVar(&flagPath, "dir", "", "Directory override, otherwise find closest '.idpzero' directory")

	parent.AddCommand(cmd)
}

func getConfigInfo(path string) (*discovery.ConfigurationInfo, error) {
	if path == "" {
		return discovery.Discover()
	}
	return discovery.EnsureDirectory(path)
}
