package cli

import (
	"fmt"
	"os"

	"github.com/idpzero/idpzero/internal/cli/initialize"
	"github.com/idpzero/idpzero/internal/cli/start"
	"github.com/spf13/cobra"
)

func init() {
	initialize.Register(rootCmd)
	start.Register(rootCmd)
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
