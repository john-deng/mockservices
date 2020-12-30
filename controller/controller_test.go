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
		Run(t)
	testApp.Get("/").
		Expect().Status(http.StatusOK).
		Body().Contains("solarmesh")
}
