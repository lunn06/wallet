package app

import (
	"context"
	"log/slog"
	"strconv"

	"github.com/lunn06/wallet/internal/config"
	gincontroller "github.com/lunn06/wallet/internal/delivery/gin"
	"github.com/lunn06/wallet/internal/domain/usecase/transation"
	"github.com/lunn06/wallet/internal/domain/usecase/wallet"
	"github.com/lunn06/wallet/internal/storage/pgx"
	"github.com/lunn06/wallet/internal/utils/pgsql"
)

// Provider is DI container that initialize
// concrete storages, usecases and controller
type Provider struct {
	logger      *slog.Logger
	controller  Controller
	initializer Initializer
	graceful    *Graceful
}

func NewProvider(cfg config.Config, logger *slog.Logger) (*Provider, error) {
	dns := pgsql.BuildDns(
		cfg.Database.Host,
		strconv.Itoa(int(cfg.Database.Port)),
		cfg.Database.SSLMode,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
	)
	storage, err := pgx.NewStorage(dns, logger)
	if err != nil {
		return nil, err
	}
	walletStorage := pgx.WalletStorage{storage}
	transactionStorage := pgx.TransactionStorage{storage}

	walletUc := wallet.NewUsecase(walletStorage)
	transactionUc := transation.NewUsecase(transactionStorage, walletStorage)

	controller := gincontroller.New(
		cfg,
		logger,
		walletUc,
		transactionUc,
	)

	graceful := NewGraceful(controller, storage)

	return &Provider{
		logger,
		controller,
		&walletUc,
		graceful,
	}, nil
}

func (p *Provider) Close(ctx context.Context) error {
	if err := p.graceful.Shutdown(ctx); err != nil {
		p.logger.Info("failed to shutdown gracefully")
		return err
	}

	p.logger.Info("Provider closed successfully")

	return nil
}
