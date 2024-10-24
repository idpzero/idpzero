package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/idpzero/idpzero/configuration"
)

type VersionInfo struct {
	Version string
	Commit  string
}

func ensureInitialized(conf *configuration.ConfigInformation) {

	conf.PrintStatus()

	if !conf.Initialized() {
		color.Yellow("Configuration not valid. Run 'idpzero init' to initialize configuration")
		fmt.Println()
		os.Exit(1)
	}
}
