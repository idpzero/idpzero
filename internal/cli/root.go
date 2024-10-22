package cli

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/common-nighthawk/go-figure"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	noBanner *bool   = new(bool)
	noColor  *bool   = new(bool)
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

	keyCmd.AddCommand(addKeyCmd)

	rootCmd.PersistentFlags().BoolVar(noBanner, "no-banner", false, "hide the banner and version information")
	rootCmd.PersistentFlags().BoolVar(noColor, "no-color", false, "disable color output")
	rootCmd.PersistentFlags().StringVar(location, "config", "", "configuration directory (default is .idpzero/ in current or parent heirachy)")
	rootCmd.AddCommand(startCmd, initializeCmd, keyCmd)

}

var rootCmd = &cobra.Command{
	Use:   "idpzero",
	Short: "Single binary IDP for simplified dev/test experience",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {

		if *noColor {
			color.NoColor = true
		}

		if !*noBanner {
			color.Yellow(figure.NewFigure("idpzero", "", true).String())
			fmt.Println()
		}

	},
}

func Execute() {

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
