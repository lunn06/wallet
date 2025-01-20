package pgx

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"

	storageLayer "github.com/lunn06/wallet/internal/storage"
)

const (
	uniqueViolationCode = "23505"
	invalidSyntax       = "22P02"
)

// handleError handle pgconn.PgError and wrap it in storageLayer errors
func handleError(err error, message string, args ...any) error {
	var pgxErr *pgconn.PgError
	if !errors.As(err, &pgxErr) {
		return storageLayer.UnhandledErr.WrapWithNoMessage(err)
	}

	switch pgxErr.Code {
	case uniqueViolationCode:
		return storageLayer.ErrUniqueViolation.New("constraint = %s", pgxErr.ConstraintName)
	case invalidSyntax:
		return storageLayer.ErrInvalid.New("constraint = %s", pgxErr.ConstraintName)
	default:
		return storageLayer.UnhandledErr.Wrap(pgxErr, message, args...)
	}
}
