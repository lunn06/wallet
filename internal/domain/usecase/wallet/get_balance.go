package wallet

import (
	"context"

	"github.com/lunn06/wallet/internal/domain/usecase"
	"github.com/lunn06/wallet/internal/dtos"
)

// GetBalance describes getting balance of target wallet
func (wuc Usecase) GetBalance(ctx context.Context, dto dtos.GetBalanceRequest) (dtos.GetBalanceResponse, error) {
	wallet, err := wuc.interactor.GetByAddress(ctx, dto.Address)
	if err != nil {
		return dtos.GetBalanceResponse{}, usecase.ErrOnGet.Wrap(err, "failed to get wallet by address")
	}

	return dtos.GetBalanceResponse{
		Balance: wallet.Balance.String(),
	}, nil
}
