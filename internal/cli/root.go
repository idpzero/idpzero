package cli

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

var (
	location *string = new(string)
	logger   *slog.Logger
)

func init() {

	logger = slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			//AddSource: true,
			Level: slog.LevelDebug,
		}),
	)

	rootCmd.PersistentFlags().StringVar(location, "config", "", "configuration directory (default is .idpzero/ in current or parent heirachy)")
	rootCmd.AddCommand(startCmd, initializeCmd)
}

var rootCmd = &cobra.Command{
	Use:   "idpzero",
	Short: "Single binary IDP for simplified dev/test experience",
}

func Execute() {

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
