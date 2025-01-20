package gin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"github.com/lunn06/wallet/internal/dtos"
)

func (gc *Controller) GetBalance(c *gin.Context) {
	address := c.Param("address")
	dto := dtos.GetBalanceRequest{
		Address: address,
	}

	if err := binding.Validator.ValidateStruct(dto); err != nil {
		c.JSON(http.StatusBadRequest, dtos.ErrorResp{Error: "INVALID_REQUEST_BODY"})
		return
	}

	response, err := gc.walletUc.GetBalance(c, dto)
	if err != nil {
		gc.logger.Error("error", "cause", err.Error())
		code, errDto := handleErr(err)
		c.JSON(code, errDto)
		return
	}

	c.JSON(http.StatusOK, response)
}
