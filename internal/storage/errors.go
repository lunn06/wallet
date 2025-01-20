package storage

import "github.com/joomcode/errorx"

func IsInternalErr(err error) bool {
	return errorx.HasTrait(err, Internal)
}

func IsExternalErr(err error) bool {
	return errorx.HasTrait(err, External)
}

func IsNotFoundErr(err error) bool {
	return errorx.HasTrait(err, errorx.NotFound())
}

func IsDuplicateErr(err error) bool {
	return errorx.HasTrait(err, errorx.Duplicate())
}

var (
	DBErrors = errorx.NewNamespace("database")

	External           = errorx.RegisterTrait("external")
	ErrNotFound        = DBErrors.NewType("not_found", External, errorx.NotFound())
	ErrInvalid         = DBErrors.NewType("invalid", External)
	ErrUniqueViolation = DBErrors.NewType("unique_violation", External, errorx.Duplicate())

	Internal             = errorx.RegisterTrait("internal")
	ErrFailedToInsert    = DBErrors.NewType("failed_to_insert", Internal)
	ErrFailedToUpdate    = DBErrors.NewType("failed_to_update", Internal)
	ErrFailedToDelete    = DBErrors.NewType("failed_to_delete", Internal)
	ErrFailedStmtBuild   = DBErrors.NewType("failed_to_build statement", Internal)
	UnhandledErr         = DBErrors.NewType("unhandled", Internal)
	ErrFailedToMarshal   = DBErrors.NewType("failed to marshal", Internal)
	ErrFailedToUnmarshal = DBErrors.NewType("failed_to_unmarshal", Internal)
)
