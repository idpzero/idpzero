package store

import (
	"database/sql"
	"embed"
	"fmt"

	"github.com/idpzero/idpzero/pkg/dbg"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func Migrate(db *sql.DB) error {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("sqlite"); err != nil {
		return err
	}

	if err := goose.Up(db, "migrations"); err != nil {
		return err
	}

	return nil
}

func init() {
	goose.SetLogger(gooseLogger{})
}

type gooseLogger struct{}

func (gl gooseLogger) Printf(format string, v ...interface{}) {
	dbg.Logger.Info(fmt.Sprintf(format, v...))
}

func (gl gooseLogger) Fatalf(format string, v ...interface{}) {
	dbg.Logger.Info(fmt.Sprintf(format, v...))
}
