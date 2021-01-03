package controller

import (
	"solarmesh.io/mockservices/src/model"
	"solarmesh.io/mockservices/src/rpc/client"
	"solarmesh.io/mockservices/src/service"

	"hidevops.io/hiboot/pkg/app"
	webctx "hidevops.io/hiboot/pkg/app/web/context"
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/hiboot/pkg/starter/httpclient"
	"hidevops.io/hiboot/pkg/starter/jaeger"
)

// controller
type Controller struct {
	// embedded at.RestController
	at.RestController
	at.RequestMapping `value:"/"`

	client httpclient.Client

	mockGRpcClientService *client.MockGRpcClientService
	mockService *service.MockService
}

// Init inject helloServiceClient
func newController(httpClient httpclient.Client,
	mockGRpcClientService *client.MockGRpcClientService,
	mockService *service.MockService) *Controller {
	return &Controller{
		client: httpClient,
		mockGRpcClientService: mockGRpcClientService,
		mockService: mockService,
	}
}

func init() {
	app.Register(newController)
}

// GET /
func (c *Controller) Get(_ struct {
	at.GetMapping `value:"/"`
}, span *jaeger.ChildSpan, ctx webctx.Context) (response *model.GetResponse) {
	var err error
	response, err = c.mockService.SendRequest("HTTP", span, ctx.Request().Header)
	c.response(err, response, ctx)
	return
}

func (c *Controller) response(err error, response *model.GetResponse, ctx webctx.Context) {
	if err == nil {
		response.Data.Url = ctx.Host() + ctx.Path()
		ctx.StatusCode(response.Code)
		for k, v := range ctx.Request().Header {
			ctx.ResponseWriter().Header().Set(k, v[0])
		}
	}
}
