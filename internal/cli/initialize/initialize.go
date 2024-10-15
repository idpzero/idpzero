package initialize

import (
	"github.com/spf13/cobra"
)

func Register(parent *cobra.Command) {

	var cmd = &cobra.Command{
		Use:   "init",
		Short: "Initalize configuration for idpzero",
		Long:  `Setup the configuration and data directory for idpzero`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	parent.AddCommand(cmd)
}
