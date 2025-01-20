package gin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"github.com/lunn06/wallet/internal/dtos"
)

func (gc *Controller) GetLast(c *gin.Context) {
	countStr, _ := c.GetQuery("count")
	count, _ := strconv.Atoi(countStr)

	dto := dtos.GetLastRequest{
		Count: count,
	}

	if err := binding.Validator.ValidateStruct(dto); err != nil {
		c.JSON(http.StatusBadRequest, dtos.ErrorResp{Error: "INVALID_REQUEST_BODY"})
		return
	}

	response, err := gc.transactionUc.GetLast(c, dto)
	if err != nil {
		gc.logger.Error("error", "cause", err.Error())
		code, errDto := handleErr(err)
		c.JSON(code, errDto)
		return
	}

	c.JSON(http.StatusOK, response)
}
