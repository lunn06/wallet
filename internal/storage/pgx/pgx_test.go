package pgx_test

import (
	"context"
	"log/slog"
	"strconv"
	"testing"

	"github.com/lunn06/wallet/internal/config"
	"github.com/lunn06/wallet/internal/storage/pgx"
	"github.com/lunn06/wallet/internal/utils/pgsql"
)

var (
	storage            *pgx.Storage
	walletStorage      pgx.WalletStorage
	transactionStorage pgx.TransactionStorage
)

// TestMain define a startup and shutdown resources for tests
func TestMain(m *testing.M) {
	cfg, err := config.ReadConfig("../../../configs/main.yaml")
	if err != nil {
		panic(err)
	}

	logger := slog.Default()
	dns := pgsql.BuildDns(
		cfg.Database.Host,
		strconv.Itoa(int(cfg.Database.Port)),
		cfg.Database.SSLMode,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
	)
	storage, err = pgx.NewStorage(dns, logger)
	if err != nil {
		panic(err)
	}

	walletStorage = pgx.WalletStorage{storage}
	transactionStorage = pgx.TransactionStorage{storage}

	m.Run() // run all tests

	if err := storage.Close(context.Background()); err != nil {
		panic(err)
	}
}
