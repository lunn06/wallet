package models

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/lunn06/wallet/internal/domain/models"
)

type Wallet struct {
	ID      int         `db:"id"`
	Address pgtype.UUID `db:"address"`
	Balance Balance     `db:"balance"`
}

func (w Wallet) TableName() string {
	return "wallets"
}

func (w Wallet) Fields() []string {
	return []string{"id", "address", "balance"}
}

func (w Wallet) FieldsWithoutID() []string {
	return w.Fields()[1:]
}

func (w Wallet) Values() []any {
	return []any{w.ID, w.Address, w.Balance}
}

func (w Wallet) ValuesWithoutID() []any {
	return w.Values()[1:]
}

func (w Wallet) ToDomain() (models.Wallet, error) {
	pguuid, err := w.Address.UUIDValue()
	if err != nil {
		return models.Wallet{}, err
	}

	balance, err := models.NewBalanceFromDecimal(w.Balance.Decimal)
	if err != nil {
		return models.Wallet{}, err
	}

	walletAddress, _ := uuid.FromBytes(pguuid.Bytes[:])
	return models.Wallet{
		ID:      w.ID,
		Address: walletAddress.String(),
		Balance: balance,
	}, nil
}

func WalletFromDomain(domain models.Wallet) (Wallet, error) {
	var dbUUID pgtype.UUID
	if err := dbUUID.Scan(domain.Address); err != nil {
		return Wallet{}, err
	}

	return Wallet{
		ID:      domain.ID,
		Address: dbUUID,
		Balance: BalanceFromDomain(domain.Balance),
	}, nil
}
