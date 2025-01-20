package transation_test

import (
	"testing"

	"github.com/lunn06/wallet/internal/domain/usecase/transation"
	"github.com/lunn06/wallet/internal/storage/mock"
)

var (
	usecaseImpl        transation.Usecase
	transactionStorage mock.TransactionStorage
	walletStorage      mock.WalletStorage
)

func TestMain(m *testing.M) {
	transactionStorage = mock.TransactionStorage{}
	walletStorage = mock.WalletStorage{}
	usecaseImpl = transation.NewUsecase(&transactionStorage, &walletStorage)

	m.Run()
}
