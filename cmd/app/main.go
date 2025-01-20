// cmd/app is app entry point that init a logger,
// read config startup and shutdown app
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/lunn06/wallet/internal/app"
	"github.com/lunn06/wallet/internal/config"
	"github.com/lunn06/wallet/pkg/zapslog"
)

func main() {
	debug := os.Getenv("LOG_LEVEL") == "debug"
	logger, sync := zapslog.Init(debug)
	defer sync()

	cfg, err := config.ReadConfig("configs/main.yaml")
	if err != nil {
		logger.Error("can't open config file", "err", err)
		return
	}

	a, err := app.New(cfg, logger)
	if err != nil {
		logger.Error("can't init app", "err", err)
		return
	}

	logger.Info("Starting app")
	go func() {
		if err := a.Run(); err != nil {
			logger.Error("can't start app", "err", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit // wait for SIGINT or SIGTERM signal

	logger.Info("Shutting down app")

	ctx := context.Background()
	if err := a.Close(ctx); err != nil {
		logger.Error("can't close app", "err", err)
		return
	}

}
