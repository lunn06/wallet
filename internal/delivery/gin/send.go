package gin

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lunn06/wallet/internal/dtos"
)

func (gc *Controller) Send(c *gin.Context) {
	var dto dtos.SendRequest
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, dtos.ErrorResp{Error: "INVALID_REQUEST_BODY"})
		return
	}

	response, err := gc.transactionUc.Send(c, dto)
	if err != nil {
		gc.logger.Error("error", "cause", err.Error())
		code, errDto := handleErr(err)
		c.JSON(code, errDto)
		return
	}

	c.JSON(http.StatusOK, response)
}
