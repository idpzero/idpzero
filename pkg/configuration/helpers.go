package configuration

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/idpzero/idpzero/pkg/console"
	"github.com/idpzero/idpzero/pkg/dbg"
)

func PrintStatus(cfg *ConfigurationManager) {

	confMsg := "Configuration File Exists"
	kconfMsg := "Keys Configuration File Exists"

	if *dbg.Debug {
		confMsg = fmt.Sprintf("%s (%s)", confMsg, cfg.GetServerPath())
		kconfMsg = fmt.Sprintf("%s (%s)", kconfMsg, cfg.GetKeysPath())
	}

	keysIcon := console.IconCheck
	keysInit, err := cfg.IsKeysInitialized()

	if err != nil {
		keysIcon = console.IconCross
		color.Red(err.Error())
	} else if !keysInit {
		keysIcon = console.IconCross
	}

	console.PrintCheck(keysIcon, kconfMsg)

	serverIcon := console.IconCheck
	serverInit, err := cfg.IsServerInitialized()

	if err != nil {
		serverIcon = console.IconCross
		color.Red(err.Error())
	} else if !serverInit {
		serverIcon = console.IconCross
	}

	console.PrintCheck(serverIcon, confMsg)

}
