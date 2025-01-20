package wallet

import (
	"context"
	"log/slog"

	"github.com/lunn06/wallet/internal/domain/models"
)

// Defining interactor interface, that define necessary to usecase methods

type walletInteractor interface {
	GetByAddress(ctx context.Context, address string) (models.Wallet, error)
	Insert(ctx context.Context, wallet models.Wallet) (models.Wallet, error)
}

// Usecase contains interactors interfaces
type Usecase struct {
	logger     *slog.Logger
	interactor walletInteractor
}

func NewUsecase(interactor walletInteractor, logger *slog.Logger) Usecase {
	if interactor == nil {
		panic("interactor can not be nil")
	}
	return Usecase{interactor: interactor, logger: logger}
}
