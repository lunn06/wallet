package pgx_test

import (
	"context"
	"math"
	"testing"
	"time"

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
	transactionSelectQuery = "SELECT * FROM transactions WHERE id = $1"
	transactionDeleteQuery = "DELETE FROM transactions WHERE id = $1"

	transactionInsertQuery = `
		INSERT INTO 
			transactions(from_address, to_address, amount, timestamp, successful) 
 		VALUES ($1, $2, $3, $4, $5) RETURNING id
	`
)

func TestTransactionStorage_GetByID(t *testing.T) {
	t.Run("get transaction by id", func(t *testing.T) {
		// Prepare input wallets
		var (
			wallet1 models.Wallet
			wallet2 models.Wallet
		)
		// Insert input wallets
		err := storage.Do(func(db *pgxpool.Pool) error {
			balance, _ := models.NewBalanceFromFloat(math.Abs(gofakeit.Float64()))
			wallet1 = models.Wallet{
				Address: uuid.NewString(),
				Balance: balance,
			}

			wallet2 = models.Wallet{
				Address: uuid.NewString(),
				Balance: balance,
			}

			newDBWallet1, _ := pgxmodels.WalletFromDomain(wallet1)
			newDBWallet2, _ := pgxmodels.WalletFromDomain(wallet2)

			values1 := newDBWallet1.ValuesWithoutID()
			err := db.QueryRow(context.Background(), walletInsertQuery, values1...).Scan(&wallet1.ID)
			if err != nil {
				return err
			}

			values2 := newDBWallet2.ValuesWithoutID()
			err = db.QueryRow(context.Background(), walletInsertQuery, values2...).Scan(&wallet2.ID)
			if err != nil {
				return err
			}

			return nil
		})
		require.NoError(t, err)

		balance, _ := models.NewBalanceFromFloat(math.Abs(gofakeit.Float64()))
		transaction := models.Transaction{
			FromAddress: wallet1.Address,
			ToAddress:   wallet2.Address,
			Amount:      balance,
			Timestamp:   time.Now(),
			Successful:  true,
		}

		// Insert transaction and scan its id
		storage.Do(func(db *pgxpool.Pool) error {
			dbTransaction, err := pgxmodels.TransactionFromDomain(transaction)
			values := dbTransaction.ValuesWithoutID()
			err = db.QueryRow(context.Background(), transactionInsertQuery, values...).Scan(&transaction.ID)
			return err
		})
		require.NoError(t, err)

		// defer cleanup func
		defer storage.Do(func(db *pgxpool.Pool) error {
			db.Exec(context.Background(), transactionDeleteQuery, transaction.ID)
			db.Exec(context.Background(), walletDeleteQuery, wallet1.ID)
			db.Exec(context.Background(), walletDeleteQuery, wallet2.ID)

			return nil
		})

		// Get transaction by id
		result, err := transactionStorage.GetByID(context.Background(), transaction.ID)
		assert.NoError(t, err)

		// Comparison predefined and returned fields
		assert.Equal(t, transaction.ID, result.ID)
		assert.Equal(t, transaction.FromAddress, result.FromAddress)
		assert.Equal(t, transaction.ToAddress, result.ToAddress)
		assert.True(t, transaction.Timestamp.Sub(result.Timestamp).Seconds() < 1)
		assert.Equal(t, transaction.Successful, result.Successful)
		assert.True(t, transaction.Amount.Equal(result.Amount))
	})
	t.Run("transaction not found", func(t *testing.T) {
		id := 999_999_999

		_, err := transactionStorage.GetByID(context.Background(), id)
		assert.ErrorContains(t, err, storageLayer.ErrNotFound.String())
	})
}

func TestTransactionStorage_GetLastSuccessful(t *testing.T) {
	t.Run("get last successful transactions", func(t *testing.T) {
		// Prepare input wallets
		var (
			wallet1 models.Wallet
			wallet2 models.Wallet
		)
		// Insert input wallets
		err := storage.Do(func(db *pgxpool.Pool) error {
			balance, _ := models.NewBalanceFromFloat(math.Abs(gofakeit.Float64()))
			wallet1 = models.Wallet{
				Address: uuid.NewString(),
				Balance: balance,
			}

			wallet2 = models.Wallet{
				Address: uuid.NewString(),
				Balance: balance,
			}

			newDBWallet1, _ := pgxmodels.WalletFromDomain(wallet1)
			newDBWallet2, _ := pgxmodels.WalletFromDomain(wallet2)

			values1 := newDBWallet1.ValuesWithoutID()
			err := db.QueryRow(context.Background(), walletInsertQuery, values1...).Scan(&wallet1.ID)
			if err != nil {
				return err
			}

			values2 := newDBWallet2.ValuesWithoutID()
			err = db.QueryRow(context.Background(), walletInsertQuery, values2...).Scan(&wallet2.ID)
			if err != nil {
				return err
			}

			return nil
		})
		require.NoError(t, err)

		// Insert 10 transactions
		transactions := make([]models.Transaction, 10)
		storage.Do(func(db *pgxpool.Pool) error {
			for i := 0; i < 10; i++ {
				balance, _ := models.NewBalanceFromFloat(math.Abs(gofakeit.Float64()))
				transaction := models.Transaction{
					FromAddress: wallet1.Address,
					ToAddress:   wallet2.Address,
					Amount:      balance,
					Timestamp:   time.Now().UTC(),
					Successful:  true,
				}
				dbTransaction, err := pgxmodels.TransactionFromDomain(transaction)
				values := dbTransaction.ValuesWithoutID()
				err = db.QueryRow(context.Background(), transactionInsertQuery, values...).Scan(&transaction.ID)
				if err != nil {
					return err
				}

				transactions[i] = transaction
			}

			return nil
		})
		require.NoError(t, err)

		// defer cleanup func
		defer storage.Do(func(db *pgxpool.Pool) error {
			db.Exec(context.Background(), walletDeleteQuery, wallet1.ID)
			db.Exec(context.Background(), walletDeleteQuery, wallet2.ID)
			for _, transaction := range transactions {
				db.Exec(context.Background(), transactionDeleteQuery, transaction.ID)
			}

			return nil
		})

		// getting last successful transactions
		result, err := transactionStorage.GetLastSuccessful(context.Background(), 10)
		assert.NoError(t, err)
		assert.Equal(t, len(transactions), len(result))

		// Comparison predefined and returned transactions(in reverse order)
		for i := range transactions {
			assert.Equal(t, transactions[len(transactions)-1-i].ID, result[i].ID)
			assert.Equal(t, transactions[len(transactions)-1-i].FromAddress, result[i].FromAddress)
			assert.Equal(t, transactions[len(transactions)-1-i].ToAddress, result[i].ToAddress)
			assert.True(t, transactions[len(transactions)-1-i].Timestamp.Sub(result[i].Timestamp).Seconds() < 1)
			assert.Equal(t, transactions[len(transactions)-1-i].Successful, result[i].Successful)
			assert.True(t, transactions[len(transactions)-1-i].Amount.Equal(result[i].Amount))
		}
	})
	t.Run("low limit", func(t *testing.T) {
		_, err := transactionStorage.GetLastSuccessful(context.Background(), 0)
		assert.ErrorContains(t, err, storageLayer.ErrInvalid.String())
	})
}

func TestTransactionStorage_Insert(t *testing.T) {
	t.Run("insert valid wallet", func(t *testing.T) {
		// Prepare input wallets
		var (
			wallet1 models.Wallet
			wallet2 models.Wallet
		)
		// Insert input wallets
		err := storage.Do(func(db *pgxpool.Pool) error {
			balance, _ := models.NewBalanceFromFloat(math.Abs(gofakeit.Float64()))
			wallet1 = models.Wallet{
				Address: uuid.NewString(),
				Balance: balance,
			}

			wallet2 = models.Wallet{
				Address: uuid.NewString(),
				Balance: balance,
			}

			newDBWallet1, _ := pgxmodels.WalletFromDomain(wallet1)
			newDBWallet2, _ := pgxmodels.WalletFromDomain(wallet2)

			values1 := newDBWallet1.ValuesWithoutID()
			err := db.QueryRow(context.Background(), walletInsertQuery, values1...).Scan(&wallet1.ID)
			if err != nil {
				return err
			}

			values2 := newDBWallet2.ValuesWithoutID()
			err = db.QueryRow(context.Background(), walletInsertQuery, values2...).Scan(&wallet2.ID)
			if err != nil {
				return err
			}

			return nil
		})
		require.NoError(t, err)

		balance, _ := models.NewBalanceFromFloat(math.Abs(gofakeit.Float64()))
		transaction := models.Transaction{
			FromAddress: wallet1.Address,
			ToAddress:   wallet2.Address,
			Amount:      balance,
			Timestamp:   time.Now().UTC(),
			Successful:  true,
		}

		// Insert transaction
		result1, err := transactionStorage.Insert(context.Background(), transaction)

		// defer cleanup func
		defer storage.Do(func(db *pgxpool.Pool) error {
			db.Exec(context.Background(), walletDeleteQuery, wallet1.ID)
			db.Exec(context.Background(), walletDeleteQuery, wallet2.ID)
			db.Exec(context.Background(), transactionDeleteQuery, result1.ID)

			return nil
		})

		// Scan inserted fields
		var dbTransaction pgxmodels.Transaction
		err = storage.Do(func(db *pgxpool.Pool) error {
			dbTransaction, _ = pgxmodels.TransactionFromDomain(result1)
			return db.QueryRow(
				context.Background(),
				transactionSelectQuery,
				dbTransaction.ID,
			).Scan(
				&dbTransaction.ID,
				&dbTransaction.FromAddress,
				&dbTransaction.ToAddress,
				&dbTransaction.Amount,
				&dbTransaction.Timestamp,
				&dbTransaction.Successful,
			)
		})

		result2, err := dbTransaction.ToDomain()
		assert.NoError(t, err)

		// Comparison inserted and scanned transactions fields
		assert.Equal(t, result1.ID, result2.ID)
		assert.Equal(t, transaction.FromAddress, result2.FromAddress)
		assert.Equal(t, transaction.ToAddress, result2.ToAddress)
		assert.Equal(t, transaction.Successful, result2.Successful)
		assert.True(t, transaction.Timestamp.Sub(result2.Timestamp).Seconds() < 1)
		assert.True(t, transaction.Amount.Equal(result2.Amount))
	})
}
