package gin

import (
	"net/http"

	"github.com/joomcode/errorx"

	"github.com/lunn06/wallet/internal/domain/usecase"
	"github.com/lunn06/wallet/internal/dtos"
)

// handleErr cast error type to errorx.Error
// to handle type of error and return matching http status
// and dto with error message
func handleErr(err error) (int, dtos.ErrorResp) {
	errx := errorx.Cast(err)
	switch {
	case usecase.IsLackOfCurrencyErr(errx):
		return http.StatusForbidden, dtos.ErrorResp{Error: "LACK_OF_CURRENCY"}
	case usecase.IsServerErr(errx):
		return http.StatusInternalServerError, dtos.ErrorResp{Error: "INTERNAL_SERVER_ERROR"}
	case usecase.IsNotFoundErr(errx):
		return http.StatusNotFound, dtos.ErrorResp{Error: "NOT_FOUND"}
	case usecase.IsDuplicateErr(errx):
		return http.StatusConflict, dtos.ErrorResp{Error: "DUPLICATE_ERROR"}
	case usecase.IsClientErr(errx):
		return http.StatusBadRequest, dtos.ErrorResp{Error: "CLIENT_ERROR"}
	default:
		return http.StatusNotImplemented, dtos.ErrorResp{Error: "NOT_IMPLEMENTED"}
	}
}
