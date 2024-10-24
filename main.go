package main

import (
	"github.com/idpzero/idpzero/cmd"
)

var (
	version = "dev"
	commit  = "none"
)

func main() {
	// Set the version and commit from built info
	cmd.Version = cmd.VersionInfo{Version: version, Commit: commit}
	cmd.Execute()
}
