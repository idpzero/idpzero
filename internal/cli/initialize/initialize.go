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

			// https://github.com/charmbracelet/bubbletea
			// https://github.com/charmbracelet/lipgloss

			// 1. Confirm location (cwd / .idpzero)
			// 2. Generate folder and include default config

			return nil
		},
	}

	parent.AddCommand(cmd)
}
