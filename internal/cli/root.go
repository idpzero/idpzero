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
	Version     VersionInfo = VersionInfo{Version: "dev", Commit: "none"}
	debug       *bool       = new(bool)
	showVersion *bool       = new(bool)
	noColor     *bool       = new(bool)
	location    *string     = new(string)
	logger      *slog.Logger
	// key
	kid     *string = new(string)
	use     *string = new(string)
	replace *bool   = new(bool)
)

func init() {

	addKeyCmd.Flags().StringVar(kid, "kid", "", "key identifier")
	addKeyCmd.Flags().StringVar(use, "use", "sig", "usage type for key")
	addKeyCmd.Flags().BoolVar(replace, "replace", false, "replace the key if it already exists")
	addKeyCmd.MarkFlagRequired("kid")
	removeKeyCmd.Flags().StringVar(kid, "kid", "", "key identifier")
	removeKeyCmd.MarkFlagRequired("kid")

	keyCmd.AddCommand(addKeyCmd, removeKeyCmd)
	rootCmd.PersistentFlags().BoolVar(debug, "debug", false, "show debug and logging in output")
	rootCmd.PersistentFlags().BoolVar(showVersion, "version", false, "show the version information in output")
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

		if *showVersion { //revive:disable:unexported-return
			color.Yellow(figure.NewFigure("idpzero", "", true).String())
			fmt.Println("v", color.MagentaString(Version.Version), "sha", color.MagentaString(Version.Commit))
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
