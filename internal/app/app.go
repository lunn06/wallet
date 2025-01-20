package app

import (
	"context"
	"log/slog"

	"github.com/lunn06/wallet/internal/config"
)

type Initializer interface {
	Initialize(ctx context.Context) error
}

type App struct {
	logger      *slog.Logger
	config      config.Config
	provider    *Provider
	initializer Initializer

	controller Controller
}

func New(config config.Config, logger *slog.Logger) (*App, error) {
	app := App{
		config: config,
		logger: logger,
	}

	if err := app.init(); err != nil {
		return nil, err
	}

	return &app, nil
}

func (app *App) Run() error {
	if err := app.controller.Run(); err != nil {
		return err
	}

	return nil
}

// Close method initiate graceful shutdown for all
// components via Provider
func (app *App) Close(ctx context.Context) error {
	err := app.provider.Close(ctx)
	if err != nil {
		return err
	}

	app.logger.Info("App close successfully")

	return nil
}

func (app *App) init() error {
	inits := []func() error{
		app.initProvider,
		app.initController,
		app.initInitializer,
	}

	for _, f := range inits {
		if err := f(); err != nil {
			return err
		}
	}

	app.logger.Info("App initialized")

	return app.initializer.Initialize(context.Background())
}

func (app *App) initProvider() error {
	provider, err := NewProvider(app.config, app.logger)
	if err != nil {
		return err
	}

	app.provider = provider
	app.logger.Info("Provider inited")

	return nil
}

func (app *App) initInitializer() error {
	app.initializer = app.provider.initializer

	app.logger.Info("Initializer inited")

	return nil
}

func (app *App) initController() error {
	app.controller = app.provider.controller

	app.logger.Info("Controller inited")

	return nil
}
