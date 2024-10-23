package main

import (
	"github.com/idpzero/idpzero/internal/cli"
)

var (
	version = "dev"
	commit  = "none"
)

func main() {
	// Set the version and commit from built info
	cli.Version = cli.VersionInfo{Version: version, Commit: commit}
	cli.Execute()
}
