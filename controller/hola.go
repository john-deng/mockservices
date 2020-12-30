package controller

import (
	"fmt"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/hiboot/pkg/log"
)

// controller
type holaController struct {
	// embedded at.RestController
	at.RestController

	AppName string `value:"${app.name}-${app.version}"`
}

// Init inject holaServiceClient
func newHolaController() *holaController {
	return &holaController{

	}
}

// GET /hola/http/{name}
func (c *holaController) GetByHttp(name string) (response string) {
	response = "Hola " + name + " de HTTP"
	response = fmt.Sprintf("[%v] Hola %v", c.AppName, name)
	log.Info(response)
	return
}

func init() {

	app.Register(newHolaController)
}
