package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/common-nighthawk/go-figure"
	"github.com/fatih/color"
	"github.com/idpzero/idpzero/cmd/initialize"
	"github.com/idpzero/idpzero/cmd/reset"

	//"github.com/idpzero/idpzero/cmd/keys"
	"github.com/idpzero/idpzero/cmd/serve"
	"github.com/idpzero/idpzero/cmd/shared"
	"github.com/idpzero/idpzero/pkg/dbg"
	"github.com/spf13/cobra"
)

var (
	showVersion *bool = new(bool)
	noColor     *bool = new(bool)
)

func init() {

	// add top level commands which add their own sub commands
	resetCmd := reset.New()
	initCmd := initialize.New()
	startCmd := serve.New()

	// root managed flags
	rootCmd.Flags().BoolVar(showVersion, "version", false, "show the version information in output")
	rootCmd.Flags().BoolVar(noColor, "no-color", false, "disable color output")

	// shared across commands
	rootCmd.PersistentFlags().BoolVar(dbg.Debug, "debug", false, "show debug and logging in output")
	rootCmd.PersistentFlags().StringVar(shared.Location, "config", "", "configuration directory (default is .idpzero/ in current or parent heirachy)")

	rootCmd.AddCommand(startCmd, initCmd, resetCmd)
}

var rootCmd = &cobra.Command{
	Use:   "idpzero",
	Short: "Single binary IDP for simplified dev/test experience",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {

		if *noColor {
			color.NoColor = true
		}

		color.Yellow(figure.NewFigure("idpzero", "", true).String())
		fmt.Println("v", color.MagentaString(dbg.Version.Version), "sha", color.MagentaString(dbg.Version.Commit))
		fmt.Println()

		if *dbg.Debug {
			// setup the logger for debugging purposes.
			dbg.Logger = slog.New(
				slog.NewTextHandler(cmd.OutOrStdout(), &slog.HandlerOptions{
					Level: slog.LevelDebug,
				}),
			)
		}
	},
}

func Execute() {

	if err := rootCmd.Execute(); err != nil {
		color.Red(err.Error())
		fmt.Println()
		os.Exit(1)
	}
}
