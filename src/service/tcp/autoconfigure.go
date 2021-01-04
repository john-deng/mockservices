package tcp

import (
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/at"
	"solarmesh.io/mockservices/src/service/tcp/server"
)

const Profile string = "tcp"

type configuration struct {
	at.AutoConfiguration

	properties *properties
}

func newConfiguration(properties *properties) *configuration {
	return &configuration{properties: properties}
}

func init() {
	app.Register(newConfiguration)
}

type ServerListener struct {

}

func (c *configuration) ServerListener(mockServer *server.MockServer) (serverListener *ServerListener) {
	if c.properties.Server.Enabled {
		mockServer.Listen(c.properties.Server.Port)
	}
	return
}
