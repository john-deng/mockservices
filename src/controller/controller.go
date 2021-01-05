package controller

import (
	"hidevops.io/hiboot/pkg/app"
	webctx "hidevops.io/hiboot/pkg/app/web/context"
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/hiboot/pkg/starter/jaeger"
	"solarmesh.io/mockservices/src/model"
	"solarmesh.io/mockservices/src/service"
)

// controller
type controller struct {
	// embedded at.RestController
	at.RestController
	at.RequestMapping `value:"/"`

	mockService    *service.MockService
}

// newController inject mockService
func newController(mockService *service.MockService) *controller {
	return &controller{
		mockService:    mockService,
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
		for k, v := range ctx.Request().Header {
			ctx.ResponseWriter().Header().Set(k, v[0])
		}
	}
	return
}

