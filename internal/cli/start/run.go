package start

import (
	"context"
	"log/slog"
)

func run(ctx context.Context, logger *slog.Logger) error {

	return nil
	// configInfo, err := getConfigInfo(flagPath)
	// if err != nil {
	// 	if err == discovery.ErrDiscoveryFailed {

	// 		fmt.Println(style.ErrorTextStyle.Render("Unable to find a '.idpzero' directory in the current or parent directories and no override privided."))
	// 		fmt.Println()
	// 		fmt.Println(style.WarningTextStyle.Render("First time running? Run 'idpzero init' to initialize default configuration."))
	// 		fmt.Println()
	// 	}
	// 	return nil
	// }

	// var serverConfig *configuration.Document = &configuration.Document{}
	// if err := configuration.Parse(serverConfig, configInfo.ConfigPath()); err != nil {
	// 	return err
	// }

	// store, err := storage.NewStorage(logger, serverConfig)
	// if err != nil {
	// 	return err
	// }

	// server, err := server.NewServer(logger, serverConfig, store)

	// if err != nil {
	// 	return err
	// }

	// server.Start()
	// <-ctx.Done()

	// // create a shutdown context which times out after 5 seconds
	// const timeoutSeconds = 5
	// shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), timeoutSeconds*time.Second)
	// defer cancelShutdown()

	// logger.With(slog.Int("timout_seconds", timeoutSeconds)).Debug("Initiating shutdown")
	// if err := server.Shutdown(shutdownCtx); err != nil {
	// 	return err
	// }

	// logger.Debug("Shutdown complete. Bye!")

	// return nil

}
