package controller

import (
	"hidevops.io/hiboot/pkg/app/web"
	"net/http"
	"testing"
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
		SetProperty("app.name", "solar-mock-app").
		Run(t)

	testApp.Get("/").
		Expect().Status(http.StatusOK).
		Body().Contains("solarmesh")
}


func TestUpstreamsFI(t *testing.T) {

	testApp := web.NewTestApp(t, newController).
		SetProperty("upstream.urls", "http://localhost:8083/,").
		SetProperty("app.name", "solar-mock-app").
		Run(t)

	testApp.Get("/").
		WithHeader("fi-app", "solar-mock-app").
		WithHeader("fi-ver", "v1").
		WithHeader("fi-cluster", "my-cluster").
		WithHeader("fi-code", "503").
		WithHeader("fi-delay", "2").
		Expect().Status(http.StatusServiceUnavailable).
		Body().Contains("solarmesh")
}

