package dbg

import (
	"io"
	"log/slog"
)

var (
	Debug   *bool        = new(bool)
	Logger  *slog.Logger = slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelDebug}))
	Version VersionInfo  = VersionInfo{Version: "dev", Commit: "none"}
)

type VersionInfo struct {
	Version string
	Commit  string
}
