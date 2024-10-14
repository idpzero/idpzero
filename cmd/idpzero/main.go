package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/common-nighthawk/go-figure"
	"github.com/idpzero/idpzero/pkg/config"
	"github.com/idpzero/idpzero/pkg/server"
	"github.com/idpzero/idpzero/pkg/storage"
	"github.com/jessevdk/go-flags"
)

type Options struct {
	Path string `short:"p" long:"path" description:"directory to store data. will search up the folder heirachy for a '.idpzero' folder if not provided." required:"false" env:"DATA_DIR"`
}

func run(ctx context.Context,
	args []string,
	stdout io.Writer,
	_ io.Writer,
) error {

	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	// parse the command line arguments
	var options = Options{}
	_, err := flags.ParseArgs(&options, args)
	if err != nil {
		return err
	}

	logger := slog.New(
		slog.NewTextHandler(stdout, &slog.HandlerOptions{
			//AddSource: true,
			Level: slog.LevelDebug,
		}),
	)

	var c config.ConfigurationInfo
	if options.Path == "" {
		path, err := config.Discover()
		if err != nil {
			return err
		}
		c = path
	} else {
		c = config.Load(options.Path)
	}

	store, err := storage.NewStorage(logger, c)
	if err != nil {
		return err
	}

	server, err := server.NewServer(logger, store)

	if err != nil {
		return err
	}

	server.Start()
	<-ctx.Done()

	// create a shutdown context which times out after 5 seconds
	const timeoutSeconds = 5
	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), timeoutSeconds*time.Second)
	defer cancelShutdown()

	logger.With(slog.Int("timout_seconds", timeoutSeconds)).Debug("Initiating shutdown")
	if err := server.Shutdown(shutdownCtx); err != nil {
		return err
	}

	logger.Debug("Shutdown complete. Bye!")

	return nil
}

func main() {

	myFigure := figure.NewColorFigure("idpzero", "", "yellow", true)
	myFigure.Print()
	fmt.Println("")

	ctx := context.Background()
	if err := run(ctx, os.Args, os.Stdout, os.Stderr); err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}

}
