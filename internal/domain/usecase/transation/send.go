package transation

import (
	"context"
	"time"

	"github.com/lunn06/wallet/internal/domain/models"
	"github.com/lunn06/wallet/internal/domain/usecase"
	"github.com/lunn06/wallet/internal/dtos"
)

// Send describes transferring balance between wallets
func (tuc Usecase) Send(ctx context.Context, dto dtos.SendRequest) (respDto dtos.SendResponse, err error) {
	if dto.FromAddress == dto.ToAddress {
		return respDto, usecase.ErrInvalid.New("invalid dto with same addresses")
	}

	from, err := tuc.walletInteractor.GetByAddress(ctx, dto.FromAddress)
	if err != nil {
		return dtos.SendResponse{}, usecase.ErrOnGet.Wrap(err, "failed to get wallet by address")
	}

	to, err := tuc.walletInteractor.GetByAddress(ctx, dto.ToAddress)
	if err != nil {
		return dtos.SendResponse{}, usecase.ErrOnGet.Wrap(err, "failed to get wallet by address")
	}

	amountBalance, err := models.NewBalanceFromString(dto.Amount)
	if err != nil {
		return dtos.SendResponse{}, usecase.ErrInvalid.Wrap(err, "invalid amount")
	}

	// to create and insert Transaction with success status
	// this defer func see on predefined err variable to define success status
	defer func() {
		tr := models.Transaction{
			FromAddress: dto.FromAddress,
			ToAddress:   dto.ToAddress,
			Amount:      amountBalance,
			Timestamp:   time.Now().UTC(),
			Successful:  err == nil,
		}
		if _, inErr := tuc.transactionInteractor.Insert(ctx, tr); inErr != nil {
			err = usecase.ErrOnInsert.Wrap(inErr, "failed to insert transaction")
		}
	}()

	if from.Balance.Less(amountBalance) {
		return dtos.SendResponse{}, usecase.ErrLackOfCurrency.New("underdraft from-wallet balance")
	}

	from.Balance, _ = from.Balance.Sub(amountBalance)
	to.Balance = to.Balance.Add(amountBalance)

	if err := tuc.walletInteractor.UpdateBalance(ctx, from); err != nil {
		return dtos.SendResponse{}, usecase.ErrOnUpdate.Wrap(err, "failed to update from-wallet")
	}
	if err := tuc.walletInteractor.UpdateBalance(ctx, to); err != nil {
		// if it can't update to-wallet, make the rollback here
		from.Balance = from.Balance.Add(amountBalance)
		if err := tuc.walletInteractor.UpdateBalance(ctx, from); err != nil {
			return dtos.SendResponse{}, usecase.ErrOnRollback.Wrap(err, "failed to rollback from-wallet")
		}
		return dtos.SendResponse{}, usecase.ErrOnUpdate.Wrap(err, "failed to update to-wallet")
	}

	return dtos.SendResponse{}, nil
}
