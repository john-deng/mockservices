package controller

import (
	"net/http"
	"testing"

	"hidevops.io/hiboot/pkg/app/web"
)

func TestBasic(t *testing.T) {

	testApp := web.NewTestApp(t, newController).SetProperty("").Run(t)
	testApp.Get("/").
		Expect().Status(http.StatusOK).
		Body().Contains("Hello")
}


func TestUpstreams(t *testing.T) {

	testApp := web.NewTestApp(t, newController).
		SetProperty("upstream.urls", "http://localhost:8083/,").
		SetProperty("app.name", "mockservices").
		Run(t)

	testApp.Get("/").
		Expect().Status(http.StatusOK).
		Body().Contains("solarmesh")
}


func TestUpstreamsFI(t *testing.T) {

	testApp := web.NewTestApp(t, newController).
		SetProperty("upstream.urls", "http://localhost:8083/,").
		SetProperty("app.name", "mockservices").
		Run(t)

	testApp.Get("/").
		WithHeader("fi-app", "mockservices").
		WithHeader("fi-ver", "v1").
		WithHeader("fi-cluster", "my-cluster").
		WithHeader("fi-code", "503").
		WithHeader("fi-delay", "2").
		Expect().Status(http.StatusServiceUnavailable).
		Body().Contains("solarmesh")
}


func TestGRpcUpstreamsFI(t *testing.T) {

	testApp := web.NewTestApp(t, newController).
		SetProperty("upstream.urls", "grpc://localhost:7575,").
		SetProperty("app.name", "mockservices").
		Run(t)

	testApp.Get("/").
		WithHeader("fi-app", "mockservices").
		WithHeader("fi-ver", "v1").
		WithHeader("fi-cluster", "my-cluster").
		//WithHeader("fi-code", "503").
		//WithHeader("fi-delay", "2").
		Expect().Status(http.StatusOK).
		Body().Contains("Success")
}


func TestTcpUpstreamsFI(t *testing.T) {

	testApp := web.NewTestApp(t, newController).
		SetProperty("upstream.urls", "tcp://localhost:8585,").
		SetProperty("app.name", "mockservices").
		Run(t)

	testApp.Get("/").
		WithHeader("fi-app", "mockservices").
		WithHeader("fi-ver", "v1").
		WithHeader("fi-cluster", "my-cluster").
		//WithHeader("fi-code", "503").
		//WithHeader("fi-delay", "20").
		Expect().Status(http.StatusOK).
		Body().Contains("Success")
}

