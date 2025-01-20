package gin

import (
	"github.com/gin-gonic/gin"
)

func (gc *Controller) setupEndpoints(r *gin.Engine) {
	base := r.Group(basePath)

	base.POST("/send", gc.Send)
	base.GET("/transactions", gc.GetLast)
	base.GET("/wallet/:address/balance", gc.GetBalance)
}
