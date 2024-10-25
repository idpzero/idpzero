package main

import (
	"github.com/idpzero/idpzero/cmd"
	"github.com/idpzero/idpzero/pkg/dbg"
)

var (
	version = "dev"
	commit  = "none"
)

func main() {
	// Set the version and commit from built info
	dbg.Version = dbg.VersionInfo{Version: version, Commit: commit}
	cmd.Execute()
}
