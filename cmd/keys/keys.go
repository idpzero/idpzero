package keys

import (
	"github.com/idpzero/idpzero/cmd/shared"
	"github.com/idpzero/idpzero/pkg/configuration"
	"github.com/spf13/cobra"
)

var (
	kid     *string = new(string)
	use     *string = new(string)
	replace *bool   = new(bool)
	conf    *configuration.ConfigInformation
)

func init() {

	// common flags
	keyCmd.PersistentFlags().StringVar(kid, "kid", "", "key identifier")
	keyCmd.MarkFlagRequired("kid")

	// add Key
	addKeyCmd.Flags().StringVar(use, "use", "sig", "usage type for key")
	addKeyCmd.Flags().BoolVar(replace, "replace", false, "replace the key if it already exists")

	// register the sub commands
	keyCmd.AddCommand(addKeyCmd, removeKeyCmd)
}

func New() *cobra.Command {
	return keyCmd
}

var keyCmd = &cobra.Command{
	Use:   "keys",
	Short: "Manage keys in the configuration file",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {

		// get the config dir to use from the path or discovery
		c, err := configuration.Resolve(*shared.Location)

		if err != nil {
			return err
		}

		if err := shared.EnsureInitialized(c); err != nil {
			return err
		}

		// set the config for the command and sub commands to use.
		conf = c

		return nil
	},
}
