package cli

import (
	"github.com/spf13/cobra"
)

var keyCmd = &cobra.Command{
	Use:   "keys",
	Short: "Manage keys in the configuration file",
}
