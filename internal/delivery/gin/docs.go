package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/go-swagno/swagno"
	"github.com/go-swagno/swagno-gin/swagger"
	"github.com/go-swagno/swagno/components/endpoint"
	"github.com/go-swagno/swagno/components/http/response"
	"github.com/go-swagno/swagno/components/mime"
	"github.com/go-swagno/swagno/components/parameter"

	"github.com/lunn06/wallet/internal/dtos"
)

func (gc *Controller) setupDocs(r *gin.Engine) {
	sw := swagno.New(swagno.Config{
		Title:       "Wallet",
		Version:     "v0.0.0",
		Description: "Wallet Backend",
		Host:        gc.config.HTTPServer.Address,
		Path:        basePath,
		Contact: &swagno.Contact{
			Name:  "Егор Титов",
			Url:   "https://github.com/lunn06/wallet",
			Email: "kuunii06@gmail.com",
		},
	})

	endpoints := []*endpoint.EndPoint{
		endpoint.New(
			endpoint.POST,
			"/send",
			endpoint.WithBody(dtos.SendRequest{}),
			endpoint.WithSuccessfulReturns([]response.Response{
				response.New(dtos.SendResponse{}, "200", ""),
			}),
			endpoint.WithErrors([]response.Response{
				response.New(dtos.ErrorResp{}, "400", "INVALID_REQUEST_BODY"),
				response.New(dtos.ErrorResp{}, "400", "CLIENT_ERROR"),
				response.New(dtos.ErrorResp{}, "403", "LACK_OF_CURRENCY"),
				response.New(dtos.ErrorResp{}, "404", "NOT_FOUND"),
				response.New(dtos.ErrorResp{}, "500", "INTERNAL_SERVER_ERROR"),
				response.New(dtos.ErrorResp{}, "501", "NOT_IMPLEMENTED"),
			}),
			endpoint.WithConsume([]mime.MIME{mime.JSON}),
			endpoint.WithProduce([]mime.MIME{mime.JSON}),
			endpoint.WithSummary("Send balance between wallets"),
		),

		endpoint.New(
			endpoint.GET,
			"/transaction",
			endpoint.WithParams(
				parameter.IntParam(
					"count",
					parameter.Query,
					parameter.WithRequired(),
					parameter.WithDefault(5),
				),
			),
			endpoint.WithSuccessfulReturns([]response.Response{
				response.New(dtos.GetLastResponse{}, "200", ""),
			}),
			endpoint.WithErrors([]response.Response{
				response.New(dtos.ErrorResp{}, "400", "INVALID_REQUEST_BODY"),
				response.New(dtos.ErrorResp{}, "400", "CLIENT_ERROR"),
				response.New(dtos.ErrorResp{}, "404", "NOT_FOUND"),
				response.New(dtos.ErrorResp{}, "500", "INTERNAL_SERVER_ERROR"),
				response.New(dtos.ErrorResp{}, "501", "NOT_IMPLEMENTED"),
			}),
			endpoint.WithConsume([]mime.MIME{mime.JSON}),
			endpoint.WithProduce([]mime.MIME{mime.JSON}),
			endpoint.WithSummary("Get last transactions"),
		),

		endpoint.New(
			endpoint.GET,
			"/wallet/{address}/balance",
			endpoint.WithParams(
				parameter.StrParam(
					"address",
					parameter.Path,
					parameter.WithRequired(),
				),
			),
			endpoint.WithSuccessfulReturns([]response.Response{
				response.New(dtos.GetBalanceResponse{}, "200", ""),
			}),
			endpoint.WithErrors([]response.Response{
				response.New(dtos.ErrorResp{}, "400", "INVALID_REQUEST_BODY"),
				response.New(dtos.ErrorResp{}, "400", "CLIENT_ERROR"),
				response.New(dtos.ErrorResp{}, "404", "NOT_FOUND"),
				response.New(dtos.ErrorResp{}, "500", "INTERNAL_SERVER_ERROR"),
				response.New(dtos.ErrorResp{}, "501", "NOT_IMPLEMENTED"),
			}),
			endpoint.WithConsume([]mime.MIME{mime.JSON}),
			endpoint.WithProduce([]mime.MIME{mime.JSON}),
			endpoint.WithSummary("Get wallet balance"),
		),
	}

	sw.AddEndpoints(endpoints)

	r.GET(basePath+"/swagger/*any", swagger.SwaggerHandler(sw.MustToJson(), swagger.Config{Prefix: "/api/swagger"}))
}
