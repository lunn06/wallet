package pgx

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/im"
	"github.com/stephenafamo/bob/dialect/psql/sm"
	"github.com/stephenafamo/bob/dialect/psql/um"

	"github.com/lunn06/wallet/internal/domain/models"
	storageLayer "github.com/lunn06/wallet/internal/storage"
	pgxmodels "github.com/lunn06/wallet/internal/storage/pgx/models"
)

type WalletStorage struct {
	*Storage
}

func (ws WalletStorage) GetByID(ctx context.Context, id int) (models.Wallet, error) {
	var wallet models.Wallet

	// access to pgxpool via embed Storage
	if err := ws.Do(func(db *pgxpool.Pool) error {
		var dbWallet pgxmodels.Wallet

		// SELECT * FROM dbWallet.TableName() WHERE id=$1 LIMIT 1
		cte := psql.Select(
			sm.From(dbWallet.TableName()),
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

		if ok := rows.Next(); !ok {
			if err := rows.Err(); err != nil {
				return handleError(err, "id = %d", id)
			}
			return storageLayer.ErrNotFound.New("wallet not found, id = %d", id)
		}

		// Marshall query output to pgxmodels.Wallet
		dbWallet, err = pgx.RowToStructByName[pgxmodels.Wallet](rows)
		if err != nil {
			return storageLayer.ErrFailedToMarshal.Wrap(err, "pgx.User = %v", dbWallet)
		}

		wallet, err = dbWallet.ToDomain()
		if err != nil {
			return storageLayer.ErrFailedToUnmarshal.Wrap(err, "pgx.User = %v", dbWallet)
		}

		return nil
	}); err != nil {
		return models.Wallet{}, err
	}

	return wallet, nil
}

func (ws WalletStorage) GetByAddress(ctx context.Context, address string) (models.Wallet, error) {
	var wallet models.Wallet

	// access to pgxpool via embed Storage
	if err := ws.Do(func(db *pgxpool.Pool) error {
		var dbWallet pgxmodels.Wallet

		// SELECT * FROM dbWallet.TableName() WHERE address=$1 LIMIT 1
		cte := psql.Select(
			sm.From(dbWallet.TableName()),
			sm.Where(psql.Quote("address").EQ(psql.Arg(address))),
			sm.Limit(1),
		)
		stmt, args, err := cte.Build(ctx)
		if err != nil {
			return storageLayer.ErrFailedStmtBuild.Wrap(err, "address = %s", address)
		}

		rows, err := db.Query(ctx, stmt, args...)
		if err != nil {
			return handleError(err, "address = %s", address)
		}
		defer rows.Close()

		if ok := rows.Next(); !ok {
			if err := rows.Err(); err != nil {
				return handleError(err, "address = %s", address)
			}
			return storageLayer.ErrNotFound.New("address = %s", address)
		}

		// Marshall query output to pgxmodels.Wallet
		dbWallet, err = pgx.RowToStructByName[pgxmodels.Wallet](rows)
		if err != nil {
			return storageLayer.ErrFailedToMarshal.Wrap(err, "pgx.Wallet = %v", dbWallet)
		}

		wallet, err = dbWallet.ToDomain()
		if err != nil {
			return storageLayer.ErrFailedToUnmarshal.Wrap(err, "pgx.Wallet = %v", dbWallet)
		}

		return nil
	}); err != nil {
		return models.Wallet{}, err
	}

	return wallet, nil
}

func (ws WalletStorage) Insert(ctx context.Context, wallet models.Wallet) (models.Wallet, error) {
	// access to pgxpool via embed Storage
	if err := ws.Do(func(db *pgxpool.Pool) error {
		newDBWallet, err := pgxmodels.WalletFromDomain(wallet)
		if err != nil {
			return err
		}

		// INSERT INTO newDBWallet.TableName() VALUES newDBWallet.ValuesWithoutID() RETURNING id
		cte := psql.Insert(
			im.Into(newDBWallet.TableName(), newDBWallet.FieldsWithoutID()...),
			im.Values(psql.Arg(newDBWallet.ValuesWithoutID()...)),
			im.Returning(psql.Quote("id")),
		)
		stmt, args, err := cte.Build(ctx)
		if err != nil {
			return storageLayer.ErrFailedStmtBuild.WrapWithNoMessage(err)
		}

		rows, err := db.Query(ctx, stmt, args...)
		if err != nil {
			return handleError(err, "error on insert wallet")
		}
		defer rows.Close()

		if ok := rows.Next(); !ok {
			if err := rows.Err(); err != nil {
				return handleError(err, "error on insert wallet")
			}
			return storageLayer.ErrFailedToInsert.Wrap(err, "error on insert wallet")
		}

		// Marshall query output to pgxmodels.Wallet
		err = rows.Scan(&newDBWallet.ID)
		if err != nil {
			return storageLayer.ErrFailedToUnmarshal.Wrap(err, "error on insert wallet")
		}

		wallet, err = newDBWallet.ToDomain()
		if err != nil {
			return storageLayer.ErrFailedToUnmarshal.Wrap(err, "error on insert wallet")
		}

		return nil
	}); err != nil {
		return models.Wallet{}, err
	}

	return wallet, nil
}

func (ws WalletStorage) UpdateBalance(ctx context.Context, changedWallet models.Wallet) error {
	// access to pgxpool via embed Storage
	if err := ws.Do(func(db *pgxpool.Pool) error {
		changedDBWallet, err := pgxmodels.WalletFromDomain(changedWallet)
		if err != nil {
			return err
		}

		// UPDATE changedDBWallet.TableName() SET balance = changedDBWallet.Balance WHERE id = changedDBWallet.ID
		cte := psql.Update(
			um.Table(changedDBWallet.TableName()),
			um.Where(psql.Quote("id").EQ(psql.Arg(changedDBWallet.ID))),
			um.SetCol("balance").ToArg(changedDBWallet.Balance),
		)
		stmt, args, err := cte.Build(ctx)
		if err != nil {
			return storageLayer.ErrFailedStmtBuild.Wrap(err, "id = %d", changedWallet.ID)
		}

		command, err := db.Exec(ctx, stmt, args...)
		if err != nil {
			return handleError(err, "id = %d", changedWallet.ID)
		}

		// If zero rows affected it means that wallet not found
		if command.RowsAffected() == 0 {
			return storageLayer.ErrNotFound.New("id = %d", changedWallet.ID)
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
