package models

import (
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/lunn06/wallet/internal/domain/models"
)

type Transaction struct {
	ID          int              `db:"id"`
	FromAddress pgtype.UUID      `db:"from_address"`
	ToAddress   pgtype.UUID      `db:"to_address"`
	Amount      Balance          `db:"amount"`
	Timestamp   pgtype.Timestamp `db:"timestamp"`
	Successful  bool             `db:"successful"`
}

func (t Transaction) TableName() string {
	return "transactions"
}

func (t Transaction) Fields() []string {
	return []string{"id", "from_address", "to_address", "amount", "timestamp", "successful"}
}

func (t Transaction) FieldsWithoutID() []string {
	return t.Fields()[1:]
}

func (t Transaction) Values() []any {
	return []any{t.ID, t.FromAddress, t.ToAddress, t.Amount, t.Timestamp, t.Successful}
}

func (t Transaction) ValuesWithoutID() []any {
	return t.Values()[1:]
}

func (t Transaction) ToDomain() (models.Transaction, error) {
	amount, err := t.Amount.ToDomain()
	if err != nil {
		return models.Transaction{}, err
	}
	return models.Transaction{
		ID:          t.ID,
		FromAddress: t.FromAddress.String(),
		ToAddress:   t.ToAddress.String(),
		Amount:      amount,
		Timestamp:   t.Timestamp.Time,
		Successful:  t.Successful,
	}, nil
}

func TransactionFromDomain(domain models.Transaction) (Transaction, error) {
	var fromDBUUID pgtype.UUID
	if err := fromDBUUID.Scan(domain.FromAddress); err != nil {
		return Transaction{}, err
	}

	var toDBUUID pgtype.UUID
	if err := toDBUUID.Scan(domain.ToAddress); err != nil {
		return Transaction{}, err
	}

	dbTimestamp := pgtype.Timestamp{
		Time:             domain.Timestamp,
		InfinityModifier: pgtype.Finite,
		Valid:            true,
	}

	return Transaction{
		ID:          domain.ID,
		FromAddress: fromDBUUID,
		ToAddress:   toDBUUID,
		Amount:      BalanceFromDomain(domain.Amount),
		Timestamp:   dbTimestamp,
		Successful:  domain.Successful,
	}, nil
}
