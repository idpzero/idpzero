package cli

import (
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/common-nighthawk/go-figure"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	debug    *bool   = new(bool)
	version  *bool   = new(bool)
	noColor  *bool   = new(bool)
	location *string = new(string)
	logger   *slog.Logger
)

func init() {

	keyCmd.AddCommand(addKeyCmd)
	rootCmd.PersistentFlags().BoolVar(debug, "debug", false, "show debug and logging in output")
	rootCmd.PersistentFlags().BoolVar(version, "version", false, "show the version information in output")
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

		if *version {
			color.Yellow(figure.NewFigure("idpzero", "", true).String())
			fmt.Println()
		}

		// default to discard logs
		output := io.Discard

		if *debug {
			output = os.Stdout
		}

		// setup the logger
		logger = slog.New(
			slog.NewTextHandler(output, &slog.HandlerOptions{
				//AddSource: true,
				Level: slog.LevelDebug,
			}),
		)

	},
}

func Execute() {

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
