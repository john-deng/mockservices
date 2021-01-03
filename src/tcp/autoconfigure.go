package tcp

import (
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/at"
	"solarmesh.io/mockservices/src/tcp/server"
)

var Profile = "tcp"

type configuration struct {
	at.AutoConfiguration

	properties *Properties
	mockTcpServer *server.MockTcpServer
}

func newConfiguration(properties *Properties,
	mockTcpServer *server.MockTcpServer) *configuration {
	return &configuration{
		properties:    properties,
		mockTcpServer: mockTcpServer,
	}
}

func init() {
	app.IncludeProfiles("tcp")
	app.Register(newConfiguration)
}

func (c *configuration) Server() (sever *Server)  {

	c.mockTcpServer.Listen(c.properties.Server.Port)

	return &Server{}
}
