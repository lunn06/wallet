package transation

import (
	"context"

	"github.com/lunn06/wallet/internal/domain/usecase"
	"github.com/lunn06/wallet/internal/dtos"
)

const defaultLimit = 5

// GetLast describe getting last successful transactions
func (tuc Usecase) GetLast(ctx context.Context, dto dtos.GetLastRequest) (dtos.GetLastResponse, error) {
	if dto.Count < 1 {
		dto.Count = defaultLimit
	}

	last, err := tuc.transactionInteractor.GetLastSuccessful(ctx, dto.Count)
	if err != nil {
		return dtos.GetLastResponse{}, usecase.ErrOnGet.Wrap(err, "failed to get last transactions")
	}

	// to transfer transaction info
	// copy models.Transaction to dtos.Transaction
	transactionsDtos := make([]dtos.Transaction, len(last))
	for i, t := range last {
		transactionsDtos[i] = dtos.Transaction{
			ID:          t.ID,
			FromAddress: t.FromAddress,
			ToAddress:   t.ToAddress,
			Amount:      t.Amount.String(),
			Timestamp:   t.Timestamp,
		}
	}

	return dtos.GetLastResponse{
		Transactions: transactionsDtos,
	}, nil
}
