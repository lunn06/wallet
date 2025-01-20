package models

import (
	"errors"

	"github.com/shopspring/decimal"
)

// Balance represents currency type with inner Decimal type
// to ensure high precision
type Balance struct {
	d decimal.Decimal
}

func (b Balance) Decimal() decimal.Decimal {
	return b.d
}

func NewBalanceFromString(s string) (Balance, error) {
	d, err := decimal.NewFromString(s)
	if err != nil {
		return Balance{}, err
	}

	if d.IsNegative() {
		return Balance{}, errors.New("balance must be positive")
	}

	return Balance{d: d}, nil
}

func NewBalanceFromDecimal(amount decimal.Decimal) (Balance, error) {
	if amount.IsNegative() {
		return Balance{}, errors.New("balance must be positive")
	}
	return Balance{amount}, nil
}

func NewBalanceFromFloat(f float64) (Balance, error) {
	if f < 0 {
		return Balance{}, errors.New("balance must be positive")
	}
	return Balance{decimal.NewFromFloat(f)}, nil
}

func (b Balance) Add(other Balance) Balance {
	return Balance{b.d.Add(other.d)}
}

func (b Balance) Sub(other Balance) (Balance, error) {
	res := b.d.Sub(other.d)
	if res.IsNegative() {
		return Balance{}, errors.New("result of Sub() must be positive")
	}
	return Balance{res}, nil
}

func (b Balance) GreaterOrEqual(other Balance) bool {
	return b.d.GreaterThanOrEqual(other.d)
}

func (b Balance) Less(other Balance) bool {
	return b.d.LessThan(other.d)
}

func (b Balance) Equal(other Balance) bool {
	return b.d.Equal(other.d)
}

func (b Balance) Float64() float64 {
	return b.d.InexactFloat64()
}

func (b Balance) String() string {
	return b.d.String()
}
