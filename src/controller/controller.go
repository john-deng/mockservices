package controller

import (
	"github.com/hidevopsio/hiboot/pkg/app"
	webctx "github.com/hidevopsio/hiboot/pkg/app/web/context"
	"github.com/hidevopsio/hiboot/pkg/at"
	"github.com/hidevopsio/hiboot/pkg/starter/jaeger"
	"github.com/john-deng/mockservices/src/model"
	"github.com/john-deng/mockservices/src/service"
)

// controller
type controller struct {
	// embedded at.RestController
	at.RestController
	at.RequestMapping `value:"/"`

	mockService *service.MockService
}

// newController inject mockService
func newController(mockService *service.MockService) *controller {
	return &controller{
		mockService: mockService,
	}
}

func init() {
	app.Register(newController)
}

// GET /
func (c *controller) Get(_ struct {
	at.GetMapping `value:"/"`
}, span *jaeger.ChildSpan, ctx webctx.Context) (response *model.Response) {
	var err error
	response, err = c.mockService.SendRequest("HTTP", span, ctx.Request().Header)
	if err == nil {
		response.Data.Url = ctx.Host() + ctx.Path()
		ctx.StatusCode(response.Code)
	}
	return
}
