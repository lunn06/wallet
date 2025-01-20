package transation_test

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lunn06/wallet/internal/domain/models"
	"github.com/lunn06/wallet/internal/domain/usecase"
	"github.com/lunn06/wallet/internal/dtos"
)

func TestUsecase_Send(t *testing.T) {
	t.Run("success sending", func(t *testing.T) {
		balance, _ := models.NewBalanceFromFloat(rand.Float64() * 100)
		wallet1, err := walletStorage.Insert(context.Background(), models.Wallet{
			Address: uuid.NewString(),
			Balance: balance,
		})
		require.NoError(t, err)

		wallet2, err := walletStorage.Insert(context.Background(), models.Wallet{
			Address: uuid.NewString(),
			Balance: balance,
		})
		require.NoError(t, err)

		_, err = usecaseImpl.Send(context.Background(), dtos.SendRequest{
			FromAddress: wallet1.Address,
			ToAddress:   wallet2.Address,
			Amount:      "3.50",
		})
		require.NoError(t, err)
		now := time.Now()

		transactions, err := transactionStorage.GetLastSuccessful(context.Background(), 1)
		require.NoError(t, err)

		amount, _ := models.NewBalanceFromString("3.50")

		transaction := transactions[0]
		assert.Equal(t, wallet1.Address, transaction.FromAddress)
		assert.Equal(t, wallet2.Address, transaction.ToAddress)
		assert.True(t, transaction.Amount.Equal(amount))
		assert.True(t, now.Sub(transaction.Timestamp).Seconds() < 1)

		wallet11, err := walletStorage.GetByID(context.Background(), wallet1.ID)
		require.NoError(t, err)

		add := wallet11.Balance.Add(amount)
		require.NoError(t, err)
		assert.True(t, wallet1.Balance.Equal(add))

		wallet22, err := walletStorage.GetByID(context.Background(), wallet2.ID)
		require.NoError(t, err)

		sub, err := wallet22.Balance.Sub(amount)
		require.NoError(t, err)
		assert.True(t, wallet1.Balance.Equal(sub))
	})

	t.Run("send with wrong amount", func(t *testing.T) {
		balance, _ := models.NewBalanceFromFloat(rand.Float64() * 100)
		wallet1, err := walletStorage.Insert(context.Background(), models.Wallet{
			Address: uuid.NewString(),
			Balance: balance,
		})
		require.NoError(t, err)

		wallet2, err := walletStorage.Insert(context.Background(), models.Wallet{
			Address: uuid.NewString(),
			Balance: balance,
		})
		require.NoError(t, err)

		_, err = usecaseImpl.Send(context.Background(), dtos.SendRequest{
			FromAddress: wallet1.Address,
			ToAddress:   wallet2.Address,
			Amount:      "wrong",
		})
		assert.ErrorContains(t, err, usecase.ErrInvalid.String())
	})

	t.Run("send with wrong address", func(t *testing.T) {
		balance, _ := models.NewBalanceFromFloat(rand.Float64() * 100)
		wallet1, err := walletStorage.Insert(context.Background(), models.Wallet{
			Address: uuid.NewString(),
			Balance: balance,
		})
		require.NoError(t, err)

		wallet2, err := walletStorage.Insert(context.Background(), models.Wallet{
			Address: uuid.NewString(),
			Balance: balance,
		})
		require.NoError(t, err)

		_, err = usecaseImpl.Send(context.Background(), dtos.SendRequest{
			FromAddress: "wrong",
			ToAddress:   wallet2.Address,
			Amount:      "3.50",
		})
		assert.ErrorContains(t, err, usecase.ErrOnGet.String())

		_, err = usecaseImpl.Send(context.Background(), dtos.SendRequest{
			FromAddress: wallet1.Address,
			ToAddress:   "wrong",
			Amount:      "3.50",
		})
		assert.ErrorContains(t, err, usecase.ErrOnGet.String())
	})

	t.Run("send with same address", func(t *testing.T) {
		balance, _ := models.NewBalanceFromFloat(rand.Float64() * 100)
		wallet, err := walletStorage.Insert(context.Background(), models.Wallet{
			Address: uuid.NewString(),
			Balance: balance,
		})
		require.NoError(t, err)

		_, err = usecaseImpl.Send(context.Background(), dtos.SendRequest{
			FromAddress: wallet.Address,
			ToAddress:   wallet.Address,
			Amount:      "3.50",
		})
		assert.ErrorContains(t, err, usecase.ErrInvalid.String())
	})
}
