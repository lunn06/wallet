package pgx_test

import (
	"context"
	"math"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lunn06/wallet/internal/domain/models"
	storageLayer "github.com/lunn06/wallet/internal/storage"
	pgxmodels "github.com/lunn06/wallet/internal/storage/pgx/models"
)

const (
	walletSelectQuery = "SELECT * FROM wallets WHERE address = $1"
	walletInsertQuery = "INSERT INTO wallets (address, balance) VALUES ($1, $2) RETURNING id"
	walletDeleteQuery = "DELETE FROM wallets WHERE address = $1"
)

func TestWalletStorage_GetByID(t *testing.T) {
	t.Run("get wallet by id", func(t *testing.T) {
		balance, _ := models.NewBalanceFromFloat(math.Abs(gofakeit.Float64()))
		wallet := models.Wallet{
			Address: uuid.NewString(),
			Balance: balance,
		}

		defer storage.Do(func(db *pgxpool.Pool) error {
			db.Exec(context.Background(), walletDeleteQuery, wallet.Address)

			return nil
		})

		err := storage.Do(func(db *pgxpool.Pool) error {
			newDBWallet, _ := pgxmodels.WalletFromDomain(wallet)

			values := newDBWallet.ValuesWithoutID()
			err := db.QueryRow(context.Background(), walletInsertQuery, values...).Scan(&wallet.ID)

			return err
		})
		require.NoError(t, err)

		result, err := walletStorage.GetByID(context.Background(), wallet.ID)
		assert.NoError(t, err)

		assert.Equal(t, wallet.ID, result.ID)
		assert.Equal(t, wallet.Address, result.Address)

		assert.True(t, result.Balance.Equal(result.Balance))
	})
	t.Run("wallet not found", func(t *testing.T) {
		id := 999_999_999

		_, err := walletStorage.GetByID(context.Background(), id)
		assert.ErrorContains(t, err, storageLayer.ErrNotFound.String())
	})
}

func TestWalletStorage_GetByAddress(t *testing.T) {
	t.Run("get wallet by address", func(t *testing.T) {
		balance, _ := models.NewBalanceFromFloat(math.Abs(gofakeit.Float64()))
		wallet := models.Wallet{
			Address: uuid.NewString(),
			Balance: balance,
		}

		defer storage.Do(func(db *pgxpool.Pool) error {
			db.Exec(context.Background(), walletDeleteQuery, wallet.Address)

			return nil
		})

		err := storage.Do(func(db *pgxpool.Pool) error {
			newDBUser, _ := pgxmodels.WalletFromDomain(wallet)

			values := newDBUser.ValuesWithoutID()
			err := db.QueryRow(context.Background(), walletInsertQuery, values...).Scan(&wallet.ID)

			return err
		})
		require.NoError(t, err)

		result, err := walletStorage.GetByAddress(context.Background(), wallet.Address)
		assert.NoError(t, err)

		assert.Equal(t, wallet.ID, result.ID)
		assert.Equal(t, wallet.Address, result.Address)

		assert.True(t, result.Balance.Equal(result.Balance))
	})
	t.Run("junior not found", func(t *testing.T) {
		_, err := walletStorage.GetByAddress(context.Background(), uuid.NewString())
		assert.ErrorContains(t, err, storageLayer.ErrNotFound.String())
	})
}

func TestWalletStorage_Insert(t *testing.T) {
	t.Run("insert valid wallet", func(t *testing.T) {
		balance, _ := models.NewBalanceFromFloat(math.Abs(gofakeit.Float64()))
		wallet := models.Wallet{
			Address: uuid.NewString(),
			Balance: balance,
		}

		defer storage.Do(func(db *pgxpool.Pool) error {
			db.Exec(context.Background(), walletDeleteQuery, wallet.Address)
			return nil
		})

		result1, err := walletStorage.Insert(context.Background(), wallet)

		assert.NoError(t, err)
		assert.Equal(t, wallet.Address, result1.Address)
		assert.True(t, wallet.Balance.Equal(result1.Balance))

		var dbWallet pgxmodels.Wallet
		err = storage.Do(func(db *pgxpool.Pool) error {
			dbWallet, _ = pgxmodels.WalletFromDomain(wallet)
			return db.QueryRow(
				context.Background(),
				walletSelectQuery,
				dbWallet.Address,
			).Scan(&dbWallet.ID, &dbWallet.Address, &dbWallet.Balance)
		})

		result2, err := dbWallet.ToDomain()
		assert.NoError(t, err)
		assert.Equal(t, result1.ID, result2.ID)
		assert.Equal(t, wallet.Address, result2.Address)
		assert.True(t, wallet.Balance.Equal(result2.Balance))
	})
}

func TestWalletStorage_UpdateBalance(t *testing.T) {
	t.Run("update wallet", func(t *testing.T) {
		balance, _ := models.NewBalanceFromFloat(math.Abs(gofakeit.Float64()))
		wallet := models.Wallet{
			Address: uuid.NewString(),
			Balance: balance,
		}

		defer storage.Do(func(db *pgxpool.Pool) error {
			db.Exec(context.TODO(), walletDeleteQuery, wallet.Address)
			return nil
		})

		err := storage.Do(func(db *pgxpool.Pool) error {
			newDBWallet, _ := pgxmodels.WalletFromDomain(wallet)
			values := newDBWallet.ValuesWithoutID()
			err := db.QueryRow(context.Background(), walletInsertQuery, values...).Scan(&wallet.ID)
			return err
		})
		require.NoError(t, err)

		amount, _ := models.NewBalanceFromFloat(10.)
		wallet.Balance = wallet.Balance.Add(amount)

		err = walletStorage.UpdateBalance(context.Background(), wallet)

		assert.NoError(t, err)

		var dbWallet pgxmodels.Wallet
		err = storage.Do(func(db *pgxpool.Pool) error {
			dbWallet, _ = pgxmodels.WalletFromDomain(wallet)
			return db.QueryRow(
				context.Background(),
				walletSelectQuery,
				dbWallet.Address,
			).Scan(&dbWallet.ID, &dbWallet.Address, &dbWallet.Balance)
		})

		result, err := dbWallet.ToDomain()
		assert.NoError(t, err)
		assert.Equal(t, wallet.ID, result.ID)
		assert.Equal(t, wallet.Address, result.Address)
		assert.True(t, wallet.Balance.Equal(result.Balance))
	})
	t.Run("wallet not found", func(t *testing.T) {
		balance, _ := models.NewBalanceFromFloat(math.Abs(gofakeit.Float64()))
		wallet := models.Wallet{
			ID:      999_999_999,
			Address: uuid.NewString(),
			Balance: balance,
		}

		err := walletStorage.UpdateBalance(context.Background(), wallet)
		assert.ErrorContains(t, err, storageLayer.ErrNotFound.String())
	})
}
