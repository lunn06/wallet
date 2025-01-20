package usecase

import (
	"github.com/joomcode/errorx"

	storageImpl "github.com/lunn06/wallet/internal/storage"
)

func IsNotFoundErr(err *errorx.Error) bool {
	return storageImpl.IsNotFoundErr(err.Cause())
}

func IsDuplicateErr(err *errorx.Error) bool {
	return storageImpl.IsDuplicateErr(err.Cause())
}

func IsLackOfCurrencyErr(err *errorx.Error) bool {
	return err.IsOfType(ErrLackOfCurrency)
}

func IsClientErr(err *errorx.Error) bool {
	return errorx.HasTrait(err, Client) || storageImpl.IsExternalErr(err.Cause())
}

func IsServerErr(err *errorx.Error) bool {
	return errorx.HasTrait(err, Server) && storageImpl.IsInternalErr(err.Cause())
}

var (
	// DomainErrors define namespace for service layer errors
	DomainErrors = errorx.NewNamespace("domain")

	// Client is errorx trait for external errors
	Client            = errorx.RegisterTrait("client")
	ErrInvalid        = DomainErrors.NewType("invalid", Client)
	ErrLackOfCurrency = DomainErrors.NewType("lack_of_currency")

	// Server is errorx trait for internal errors
	Server        = errorx.RegisterTrait("server")
	ErrOnGet      = DomainErrors.NewType("failed_to_get", Server)
	ErrOnInsert   = DomainErrors.NewType("failed_to_insert", Server)
	ErrOnRollback = DomainErrors.NewType("failed_to_rollback", Server)
	ErrOnUpdate   = DomainErrors.NewType("failed_to_update token", Server)
)
