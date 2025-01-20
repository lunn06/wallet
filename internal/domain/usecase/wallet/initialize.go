package wallet

import (
	"context"

	"github.com/google/uuid"

	"github.com/lunn06/wallet/internal/domain/models"
	"github.com/lunn06/wallet/internal/domain/usecase"
)

// Initialize implements app.Initializer interface to define startup behavior
func (wuc Usecase) Initialize(ctx context.Context) error {
	for i := 0; i < 10; i++ {
		balance, _ := models.NewBalanceFromFloat(100.)
		wallet := models.Wallet{
			Address: uuid.NewString(),
			Balance: balance,
		}
		if _, err := wuc.interactor.Insert(ctx, wallet); err != nil {
			return usecase.ErrOnInsert.Wrap(err, "failed to insert wallet")
		}
	}

	return nil
}
