package mock

import (
	"context"
	"slices"

	"github.com/lunn06/wallet/internal/domain/models"
	storageLayer "github.com/lunn06/wallet/internal/storage"
)

type TransactionStorage struct {
	in []models.Transaction
}

func (ts *TransactionStorage) GetByID(ctx context.Context, id int) (models.Transaction, error) {
	if ts.in == nil {
		ts.in = make([]models.Transaction, 0)
	}
	for _, t := range ts.in {
		if t.ID == id {
			return t, nil
		}
	}

	return models.Transaction{}, storageLayer.ErrNotFound.New("transaction not found, id = %d", id)
}

func (ts *TransactionStorage) GetLastSuccessful(ctx context.Context, limit int) ([]models.Transaction, error) {
	if ts.in == nil {
		ts.in = make([]models.Transaction, 0)
	}
	if limit < 1 {
		return nil, storageLayer.ErrInvalid.New("limit must be greater than zero")
	}

	slices.SortFunc(ts.in, func(e models.Transaction, e2 models.Transaction) int {
		if e.Timestamp.After(e2.Timestamp) {
			return 1
		} else if e.Timestamp.Before(e2.Timestamp) {
			return -1
		}

		return 0
	})

	return ts.in[:limit:limit], nil
}

func (ts *TransactionStorage) Insert(ctx context.Context, transaction models.Transaction) (models.Transaction, error) {
	if ts.in == nil {
		ts.in = make([]models.Transaction, 0)
	}
	var index int
	for _, t := range ts.in {
		index = max(t.ID, index)
	}

	transaction.ID = index + 1
	ts.in = append(ts.in, transaction)

	return transaction, nil
}
