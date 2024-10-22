package cli

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/idpzero/idpzero/internal/config"
)

func configDebug(cfg *config.ConfigInformation) {

	fmt.Println("Verifying IDP configuration...")

	printCheck(cfg.Directory().Exists(), "Configuration Directory Exists")
	printCheck(cfg.Config().Exists(), "Configuration File Exists")

	fmt.Println()

}

func printCheck(passed bool, msg string) {
	var mark = color.RedString(" x ")
	if passed {
		mark = color.GreenString(" ✓ ")
	}

	fmt.Println(mark, msg)
}
