package models

import (
	"github.com/shopspring/decimal"

	"github.com/lunn06/wallet/internal/domain/models"
)

type Balance struct {
	decimal.Decimal
}

func (b Balance) ToDomain() (models.Balance, error) {
	return models.NewBalanceFromDecimal(b.Decimal)
}

func BalanceFromDomain(b models.Balance) Balance {
	return Balance{b.Decimal()}
}
