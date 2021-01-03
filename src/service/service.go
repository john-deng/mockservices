package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/opentracing/opentracing-go"
	olog "github.com/opentracing/opentracing-go/log"
	"golang.org/x/net/context"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hiboot/pkg/starter/httpclient"
	"hidevops.io/hiboot/pkg/starter/jaeger"
	"solarmesh.io/mockservices/src/model"
	"solarmesh.io/mockservices/src/rpc/client"
)

// MockService
type MockService struct {

	client httpclient.Client

	mockGRpcClientService *client.MockGRpcClientService

	UpstreamUrls string   `value:"${upstream.urls}"`
	Upstreams    []string `value:"${upstream.urls}"`
	AppName      string   `value:"${app.name}"`
	Version      string   `value:"${app.version}"`
	ClusterName  string   `value:"${cluster.name:my-cluster}"`
	UserData     string   `value:"${user.data:solarmesh}"`
}

func newMockService(httpClient httpclient.Client, mockGRpcClientService *client.MockGRpcClientService) *MockService {
	return &MockService{
		client: httpClient,
		mockGRpcClientService: mockGRpcClientService,
	}
}

func init() {
	app.Register(newMockService)
}

func (c *MockService) SendRequest(protocol string, span *jaeger.ChildSpan, header http.Header) (response *model.GetResponse, err error) {
	response = new(model.GetResponse)

	if span != nil && span.Span != nil {
		defer span.Finish()
		greeting := span.BaggageItem(c.AppName)
		if greeting == "" {
			greeting = "Hello"
		}
	}
	var newSpan opentracing.Span

	//response.Data.UserAgent = ua.Parse(ctx.GetHeader("User-Agent"))
	fiApp := header.Get("fi-app")
	fiVer := header.Get("fi-ver")
	fiCluster := header.Get("fi-cluster")
	fiCode, _ := strconv.Atoi(header.Get("fi-code"))
	fiDelay, _ := strconv.Atoi(header.Get("fi-delay"))
	response.Data.Protocol = protocol
	response.Data.App = c.AppName
	response.Data.Version = c.Version
	response.Data.Cluster = c.ClusterName
	response.Data.UserData = c.UserData

	upstreamUrls := c.parseUpstream()

	urlLens := len(upstreamUrls)
	if urlLens == 0 || urlLens != 0 && upstreamUrls[0] == "${upstream.urls}" {
		response.Data.MetaData = " -> " + c.AppName
	} else {
		response.Data.MetaData = " -> " + c.AppName + " -> "
		for _, upstreamUrl := range upstreamUrls {
			if upstreamUrl != "" {
				upstreamResponse := new(model.GetResponse)
				upstreamResponse.Data.Url = upstreamUrl
				var newResp string

				u, err := url.Parse(upstreamUrl)
				if err != nil {
					log.Warnf("Bad URL format: %v", upstreamUrl)
					continue
				}
				var resp *http.Response

				switch u.Scheme {
				case "http", "https":
					resp, err = c.client.Get(upstreamUrl, header, func(req *http.Request) {
						if span.Span != nil {
							newSpan = span.Inject(context.Background(), "GET", upstreamUrl, req)
						}
					})
					if err == nil {
						byteResp, _ := ioutil.ReadAll(resp.Body)
						_ = json.Unmarshal(byteResp, upstreamResponse)
					}
				case "grpc":
					ctx, cancel := context.WithTimeout(context.Background(), time.Second)
					mockResponse, err := c.mockGRpcClientService.Send(ctx, u.Host, header)
					if err == nil {
						b, e := json.Marshal(mockResponse)
						if e == nil {
							e = json.Unmarshal(b, upstreamResponse)
						}
					}
					cancel()
				}

				if err != nil {
					errMsg := err.Error()
					upstreamResponse.Code = 500
					upstreamResponse.Message = errMsg
					log.Error(errMsg)
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

		hasFaultInjection := true

		if fiVer != "" && fiVer != c.Version {
			hasFaultInjection = false
		}

		if fiCluster != "" && fiCluster != c.ClusterName {
			hasFaultInjection = false
		}

		if hasFaultInjection {
			faultInjectionMessage := ""
			if fiDelay != 0 {
				time.Sleep(time.Duration(fiDelay) * time.Millisecond)
				faultInjectionMessage += fmt.Sprintf(" with delay %d ms,", fiDelay)
			}
			if fiCode != 0 {
				response.Code = fiCode
				faultInjectionMessage += fmt.Sprintf(" with HTTP status code %d,", fiCode)
			}

			if fiCode != 0 || fiDelay != 0 {
				response.Message = fmt.Sprintf("Fault Injection%v", faultInjectionMessage)
			}
		}
	}

	respStr, _ := json.Marshal(response)
	log.Info(string(respStr))

	if span != nil && span.Span != nil {
		span.LogFields(
			olog.String("event", c.AppName),
			olog.String("value", string(respStr)),
		)
	}
	return
}

func (c *MockService) parseUpstream() []string {
	log.Infof("Upstreams: %v", c.Upstreams)
	log.Infof("UpstreamUrls: %v", c.UpstreamUrls)

	upstreamUrls := strings.SplitN(c.UpstreamUrls, ",", -1)
	log.Debugf("len of urls: %v", len(c.UpstreamUrls))

	// TODO: it is a patch, to be fixed
	if c.UpstreamUrls == "" && len(c.Upstreams) != 0 {
		upstreamUrls = append(upstreamUrls, c.Upstreams...)
	}
	return upstreamUrls
}