package controller

import (
	"hidevops.io/hiboot/pkg/app/web"
	"net/http"
	"testing"
)

func TestHelloClient(t *testing.T) {

	testApp := web.RunTestApplication(t, newHelloController)
	testApp.Get("/hello/http/{http}").
		WithPath("http", "test").
		Expect().Status(http.StatusOK).
		Body().Contains("hello test")
}
