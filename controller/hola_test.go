package controller

import (
	"hidevops.io/hiboot/pkg/app/web"
	"net/http"
	"testing"
)

func TestHolaClient(t *testing.T) {


	testApp := web.RunTestApplication(t, newHolaController)
	testApp.Get("/hola/http/{http}").
		WithPath("http", "test").
		Expect().Status(http.StatusOK).
		Body().Contains("hola test")
}
