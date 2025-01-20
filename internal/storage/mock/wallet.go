package mock

import (
	"context"

	"github.com/lunn06/wallet/internal/domain/models"
	storageLayer "github.com/lunn06/wallet/internal/storage"
)

type WalletStorage struct {
	in []models.Wallet
}

func (ws *WalletStorage) GetByID(ctx context.Context, id int) (models.Wallet, error) {
	if ws.in == nil {
		ws.in = make([]models.Wallet, 0)
	}
	for _, w := range ws.in {
		if w.ID == id {
			return w, nil
		}
	}

	return models.Wallet{}, storageLayer.ErrNotFound.New("wallet not found, id = %d", id)
}

func (ws *WalletStorage) GetByAddress(ctx context.Context, address string) (models.Wallet, error) {
	if ws.in == nil {
		ws.in = make([]models.Wallet, 0)
	}
	for _, w := range ws.in {
		if w.Address == address {
			return w, nil
		}
	}

	return models.Wallet{}, storageLayer.ErrNotFound.New("address = %s", address)
}

func (ws *WalletStorage) Insert(ctx context.Context, wallet models.Wallet) (models.Wallet, error) {
	if ws.in == nil {
		ws.in = make([]models.Wallet, 0, 1)
	}
	var index int
	for _, w := range ws.in {
		index = max(w.ID, index)
	}

	wallet.ID = index + 1
	ws.in = append(ws.in, wallet)

	return wallet, nil
}

func (ws *WalletStorage) UpdateBalance(ctx context.Context, wallet models.Wallet) error {
	if ws.in == nil {
		ws.in = make([]models.Wallet, 0)
	}
	for i, w := range ws.in {
		if w.ID == wallet.ID {
			ws.in[i].Balance = wallet.Balance
		}
	}

	return nil
}
