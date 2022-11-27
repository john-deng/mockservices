package tcp

import (
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/at"
	"github.com/john-deng/mockservices/src/service/tcp/server"
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
