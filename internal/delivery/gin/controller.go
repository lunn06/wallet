// Package gin contains http controller implementation of delivery layer
package gin

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	sloggin "github.com/samber/slog-gin"

	"github.com/lunn06/wallet/internal/config"
	transactionUc "github.com/lunn06/wallet/internal/domain/usecase/transation"
	walletUc "github.com/lunn06/wallet/internal/domain/usecase/wallet"
)

const basePath = "/api"

// Controller connect http endpoints with domain usecases
type Controller struct {
	logger *slog.Logger
	server *http.Server
	config config.Config

	walletUc      walletUc.Usecase
	transactionUc transactionUc.Usecase
}

func (gc *Controller) Run() error {
	gc.logger.Info("Controller started", slog.String("addr", gc.server.Addr))
	if err := gc.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (gc *Controller) Close(ctx context.Context) error {
	err := gc.server.Shutdown(ctx)
	if err == nil {
		gc.logger.Info("Controller close successfully")
	}

	return err
}

func New(
	cfg config.Config,
	logger *slog.Logger,
	walletUc walletUc.Usecase,
	transactionUc transactionUc.Usecase,
) *Controller {
	controller := Controller{
		logger:        logger,
		config:        cfg,
		walletUc:      walletUc,
		transactionUc: transactionUc,
	}

	r := gin.New()
	r.Use(
		// connect controller's logger to gin handler
		sloggin.New(logger),
		gin.Recovery(),
	)

	controller.setupDocs(r)
	controller.setupEndpoints(r)

	// connect validator to gin handler
	binding.Validator = new(defaultValidator)

	addr := fmt.Sprintf(":%v", controller.config.HTTPServer.Port)
	controller.server = &http.Server{
		Addr:    addr,
		Handler: r,
	}

	controller.logger.Info("Controller created")

	return &controller
}
