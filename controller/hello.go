package controller

import (
	"fmt"
	"strconv"

	"github.com/mileusna/useragent"
	olog "github.com/opentracing/opentracing-go/log"
	"hidevops.io/hiboot/pkg/app"
	webctx "hidevops.io/hiboot/pkg/app/web/context"
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hiboot/pkg/starter/httpclient"
	"hidevops.io/hiboot/pkg/starter/jaeger"
)

// controller
type helloController struct {
	// embedded at.RestController
	at.RestController

	client  httpclient.Client
	AppName string `value:"${app.name}-${app.version}"`
	UserData string `value:"${user.data}"`
}

// Init inject helloServiceClient
func newHelloController(httpClient httpclient.Client) *helloController {
	return &helloController{
		client: httpClient,
	}
}

// GET /hello/http/{http}
func (c *helloController) GetByHttp(name string, span *jaeger.ChildSpan, ctx webctx.Context) (response string) {
	if span.Span != nil {
		defer span.Finish()
		greeting := span.BaggageItem("greeting")
		if greeting == "" {
			greeting = "Hello"
		}
	}
	u := ua.Parse(ctx.GetHeader("User-Agent"))
	fiSvc := ctx.GetHeader("fi-svc")
	fiCode, _ := strconv.Atoi(ctx.GetHeader("fi-code"))

	if fiSvc == c.AppName {
		ctx.StatusCode(fiCode)
		response = fmt.Sprintf("[%v.%v][%v] Hello %v with %d", c.AppName, c.UserData, u.OS + "-" + u.Name , name, fiCode)
	} else {
		response = fmt.Sprintf("[%v.%v][%v] Hello %v ", c.AppName, c.UserData, u.OS + "-" + u.Name , name)
	}

	log.Info(response)

	if span.Span != nil {
		span.LogFields(
			olog.String("event", "string-format"),
			olog.String("value", name),
		)
	}
	return
}

func init() {
	app.Register(newHelloController)
}
