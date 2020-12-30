package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	//"github.com/mileusna/useragent"
	"github.com/opentracing/opentracing-go"
	olog "github.com/opentracing/opentracing-go/log"
	"golang.org/x/net/context"
	"hidevops.io/hiboot/pkg/app"
	webctx "hidevops.io/hiboot/pkg/app/web/context"
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hiboot/pkg/starter/httpclient"
	"hidevops.io/hiboot/pkg/starter/jaeger"
)

// controller
type Controller struct {
	// embedded at.RestController
	at.RestController
	at.RequestMapping `value:"/"`

	client httpclient.Client

	UpstreamUrls string   `value:"${upstream.urls}"`
	Upstreams    []string `value:"${upstream.urls}"`
	AppName      string   `value:"${app.name}"`
	Version      string   `value:"${app.version}"`
	ClusterName  string   `value:"${cluster.name:my-cluster}"`
	UserData     string   `value:"${user.data:solarmesh}"`
}

// Init inject helloServiceClient
func newController(httpClient httpclient.Client) *Controller {
	return &Controller{
		client: httpClient,
	}
}

func init() {
	app.Register(newController)
}

type ResponseData struct {
	Url              string
	App              string
	Version          string
	SourceApp        string
	SourceAppVersion string
	Cluster          string
	UserData         string
	MetaData         string
	//UserAgent ua.UserAgent
	Header   http.Header
	Upstream []*GetResponse
}

type GetResponse struct {
	Code    int
	Message string
	Data    ResponseData
}

// GET /
func (c *Controller) Get(_ struct {
	at.GetMapping `value:"/"`
}, span *jaeger.ChildSpan, ctx webctx.Context) (response *GetResponse) {
	response = new(GetResponse)

	if span.Span != nil {
		defer span.Finish()
		greeting := span.BaggageItem(c.AppName)
		if greeting == "" {
			greeting = "Hello"
		}
	}
	var newSpan opentracing.Span

	//response.Data.UserAgent = ua.Parse(ctx.GetHeader("User-Agent"))
	fiApp := ctx.GetHeader("fi-app")
	fiVer := ctx.GetHeader("fi-ver")
	fiCluster := ctx.GetHeader("fi-cluster")
	fiCode, _ := strconv.Atoi(ctx.GetHeader("fi-code"))
	fiDelay, _ := strconv.Atoi(ctx.GetHeader("fi-delay"))
	response.Data.Url = ctx.Host() + ctx.Path()
	response.Data.App = c.AppName
	response.Data.Version = c.Version
	response.Data.Cluster = c.ClusterName
	response.Data.UserData = c.UserData
	response.Data.Header = ctx.Request().Header

	log.Infof("Upstreams: %v", c.Upstreams)
	log.Infof("UpstreamUrls: %v", c.UpstreamUrls)

	upstreamUrls := strings.SplitN(c.UpstreamUrls, ",", -1)
	log.Debugf("len of urls: %v", len(c.UpstreamUrls))

	// TODO: it is a patch, to be fixed
	if c.UpstreamUrls == "" && len(c.Upstreams) != 0 {
		upstreamUrls = append(upstreamUrls, c.Upstreams...)
	}

	urlLens := len(upstreamUrls)
	if urlLens == 0 || urlLens != 0 && upstreamUrls[0] == "${upstream.urls}" {
		response.Data.MetaData = " -> " + c.AppName
	} else {
		response.Data.MetaData = " -> " + c.AppName + " -> "
		for _, upstreamUrl := range upstreamUrls {
			if upstreamUrl != "" {
				upstreamResponse := new(GetResponse)
				resp, err := c.client.Get(upstreamUrl, ctx.Request().Header, func(req *http.Request) {
					if span.Span != nil {
						newSpan = span.Inject(context.Background(), "GET", upstreamUrl, req)
					}
				})
				var newResp string
				if err != nil {
					errMsg := err.Error()
					upstreamResponse.Data.Url = upstreamUrl
					upstreamResponse.Code = 500
					upstreamResponse.Message = errMsg
					log.Error(errMsg)
				} else {
					byteResp, _ := ioutil.ReadAll(resp.Body)
					_ = json.Unmarshal(byteResp, upstreamResponse)
				}
				upstreamResponse.Data.SourceApp = c.AppName
				upstreamResponse.Data.SourceAppVersion = c.Version
				response.Data.Upstream = append(response.Data.Upstream, upstreamResponse)

				if newSpan != nil {
					newSpan.LogFields(
						olog.String("event", upstreamUrl),
						olog.String("value", newResp),
					)
				}
			}
		}
	}

	response.Code = 200
	response.Message = "Success"
	if fiApp == c.AppName {

		hasFaultInjection := false

		if fiVer != "" && fiVer == c.Version {
			hasFaultInjection = true
		}

		if fiCluster != "" && fiCluster == c.ClusterName {
			hasFaultInjection = true
		}

		if fiVer == "" || hasFaultInjection {
			faultInjectionMessage := ""
			if fiCode != 0 {
				ctx.StatusCode(fiCode)
				faultInjectionMessage += fmt.Sprintf(" with code: %d,", fiCode)
			}
			if fiDelay != 0 {
				time.Sleep(time.Duration(fiDelay) * time.Millisecond)
				faultInjectionMessage += fmt.Sprintf(" with delay: %d ms,", fiDelay)
			}
			response.Code = fiCode
			response.Message = fmt.Sprintf("Fault Injection %v", faultInjectionMessage)
		}
	}

	respStr, _ := json.Marshal(response)
	log.Info(string(respStr))

	if span.Span != nil {
		span.LogFields(
			olog.String("event", c.AppName),
			olog.String("value", string(respStr)),
		)
	}

	return
}
