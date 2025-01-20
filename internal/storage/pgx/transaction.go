package pgx

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/im"
	"github.com/stephenafamo/bob/dialect/psql/sm"

	"github.com/lunn06/wallet/internal/domain/models"
	storageLayer "github.com/lunn06/wallet/internal/storage"
	pgxmodels "github.com/lunn06/wallet/internal/storage/pgx/models"
)

type TransactionStorage struct {
	*Storage
}

func (ts TransactionStorage) GetLastSuccessful(ctx context.Context, limit int) ([]models.Transaction, error) {
	if limit < 1 {
		return nil, storageLayer.ErrInvalid.New("limit must be greater than zero")
	}

	transactions := make([]models.Transaction, 0, limit)

	// access to pgxpool via embed Storage
	if err := ts.Do(func(db *pgxpool.Pool) error {
		var dbTransaction pgxmodels.Transaction

		// SELECT * FROM dbTransaction.TableName() WHERE successful = true ORDER BY timestamp DESC LIMIT $1
		cte := psql.Select(
			sm.From(dbTransaction.TableName()),
			sm.Where(psql.Quote("successful").EQ(psql.Arg(true))),
			sm.OrderBy(psql.Quote("timestamp")).Desc(),
			sm.Limit(limit),
		)
		stmt, args, err := cte.Build(ctx)
		if err != nil {
			return storageLayer.ErrFailedStmtBuild.WrapWithNoMessage(err)
		}

		rows, err := db.Query(ctx, stmt, args...)
		if err != nil {
			return handleError(err, "error on GetLastSuccessful transaction")
		}
		defer rows.Close()

		for rows.Next() {
			// Marshall query output to pgxmodels.Transaction
			dbTransaction, err = pgx.RowToStructByName[pgxmodels.Transaction](rows)
			if err != nil {
				return storageLayer.ErrFailedToMarshal.Wrap(err, "pgx.Transactions = %v", dbTransaction)
			}

			transaction, err := dbTransaction.ToDomain()
			if err != nil {
				return storageLayer.ErrFailedToUnmarshal.Wrap(err, "pgx.Transaction = %v", dbTransaction)
			}
			transactions = append(transactions, transaction)
		}
		if err = rows.Err(); err != nil {
			return handleError(err, "error on GetLastSuccessful transaction")
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return transactions, nil
}

func (ts TransactionStorage) GetByID(ctx context.Context, id int) (models.Transaction, error) {
	var transaction models.Transaction

	// access to pgxpool via embed Storage
	if err := ts.Do(func(db *pgxpool.Pool) error {
		var dbTransaction pgxmodels.Transaction

		// SELECT * FROM dbTransaction.TableName() WHERE id=$1 LIMIT 1
		cte := psql.Select(
			sm.From(dbTransaction.TableName()),
			sm.Where(psql.Quote("id").EQ(psql.Arg(id))),
			sm.Limit(1),
		)
		stmt, args, err := cte.Build(ctx)
		if err != nil {
			return storageLayer.ErrFailedStmtBuild.Wrap(err, "id = %d", id)
		}

		rows, err := db.Query(ctx, stmt, args...)
		if err != nil {
			return handleError(err, "id = %d", id)
		}
		defer rows.Close()

		// if no rows it means that Transaction not found or err
		if ok := rows.Next(); !ok {
			if err := rows.Err(); err != nil {
				return handleError(err, "id = %d", id)
			}
			return storageLayer.ErrNotFound.New("transaction not found, id = %d", id)
		}

		// Marshall query output to pgxmodels.Transaction
		dbTransaction, err = pgx.RowToStructByName[pgxmodels.Transaction](rows)
		if err != nil {
			return storageLayer.ErrFailedToMarshal.Wrap(err, "pgx.Transaction = %v", dbTransaction)
		}

		transaction, err = dbTransaction.ToDomain()
		if err != nil {
			return storageLayer.ErrFailedToUnmarshal.Wrap(err, "pgx.Transaction = %v", dbTransaction)
		}

		return nil
	}); err != nil {
		return models.Transaction{}, err
	}

	return transaction, nil
}

func (ts TransactionStorage) Insert(ctx context.Context, transaction models.Transaction) (models.Transaction, error) {
	// access to pgxpool via embed Storage
	if err := ts.Do(func(db *pgxpool.Pool) error {
		newDBTransaction, err := pgxmodels.TransactionFromDomain(transaction)
		if err != nil {
			return err
		}

		// INSERT INTO newDBWTransaction.TableName() VALUES newDBTransaction.ValuesWithoutID() RETURNING id
		cte := psql.Insert(
			im.Into(newDBTransaction.TableName(), newDBTransaction.FieldsWithoutID()...),
			im.Values(psql.Arg(newDBTransaction.ValuesWithoutID()...)),
			im.Returning(psql.Quote("id")),
		)
		stmt, args, err := cte.Build(ctx)
		if err != nil {
			return storageLayer.ErrFailedStmtBuild.WrapWithNoMessage(err)
		}

		rows, err := db.Query(ctx, stmt, args...)
		if err != nil {
			return handleError(err, "error on insert transaction")
		}
		defer rows.Close()

		// if no rows it means that Transaction not found or err
		if ok := rows.Next(); !ok {
			if err := rows.Err(); err != nil {
				return handleError(err, "error on insert transaction")
			}
		}

		// Marshall query output to pgxmodels.Transaction
		err = rows.Scan(&newDBTransaction.ID)
		if err != nil {
			return storageLayer.ErrFailedToUnmarshal.Wrap(err, "error on insert transaction")
		}

		transaction, err = newDBTransaction.ToDomain()
		if err != nil {
			return storageLayer.ErrFailedToUnmarshal.Wrap(err, "pgx.Transaction = %v", newDBTransaction)
		}

		return nil
	}); err != nil {
		return models.Transaction{}, err
	}

	return transaction, nil
}
