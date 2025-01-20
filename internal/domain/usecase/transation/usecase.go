package transation

import (
	"context"

	"github.com/lunn06/wallet/internal/domain/models"
)

// Defining interactors interfaces, that define necessary to usecaseImpl methods

type walletInteractor interface {
	GetByAddress(ctx context.Context, address string) (models.Wallet, error)
	Insert(ctx context.Context, wallet models.Wallet) (models.Wallet, error)
	UpdateBalance(ctx context.Context, wallet models.Wallet) error
}

type transactionInteractor interface {
	GetLastSuccessful(ctx context.Context, limit int) ([]models.Transaction, error)
	Insert(ctx context.Context, transaction models.Transaction) (models.Transaction, error)
}

// Usecase contains interactors interfaces
type Usecase struct {
	transactionInteractor transactionInteractor
	walletInteractor      walletInteractor
}

func NewUsecase(transactionInteractor transactionInteractor, walletInteractor walletInteractor) Usecase {
	if transactionInteractor == nil || walletInteractor == nil {
		panic("interactor can not be nil")
	}
	return Usecase{
		transactionInteractor: transactionInteractor,
		walletInteractor:      walletInteractor,
	}
}
