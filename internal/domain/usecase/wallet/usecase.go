package wallet

import (
	"context"

	"github.com/lunn06/wallet/internal/domain/models"
)

// Defining interactor interface, that define necessary to usecase methods

type walletInteractor interface {
	GetByAddress(ctx context.Context, address string) (models.Wallet, error)
	Insert(ctx context.Context, wallet models.Wallet) (models.Wallet, error)
}

// Usecase contains interactors interfaces
type Usecase struct {
	interactor walletInteractor
}

func NewUsecase(interactor walletInteractor) Usecase {
	if interactor == nil {
		panic("interactor can not be nil")
	}
	return Usecase{interactor: interactor}
}
