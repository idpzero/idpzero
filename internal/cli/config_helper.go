package cli

import (
	"fmt"

	"github.com/idpzero/idpzero/internal/config"
	"github.com/idpzero/idpzero/internal/style"
)

func configDebug(cfg *config.ConfigInformation) {

	fmt.Println("Verifying IDP configuration...")

	if cfg.Directory().Exists() {
		fmt.Print(style.GreenTextStyle.Render(" ✓ "))
	} else {
		fmt.Print(style.ErrorTextStyle.Render(" x "))
	}

	fmt.Println("Configuration Directory Exists")

	if cfg.Config().Exists() {
		fmt.Print(style.GreenTextStyle.Render(" ✓ "))
	} else {
		fmt.Print(style.ErrorTextStyle.Render(" x "))
	}

	fmt.Println("Configuration File Exists")
	fmt.Println()

}
