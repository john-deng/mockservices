package tcp

import (

	"hidevops.io/hiboot/pkg/at"
)

type Server struct {
	Enabled bool `json:"enabled"`
	Port string `json:"port" default:"8585"`
}

type Client struct {
	Enabled bool `json:"enabled"`
	Port string `json:"port" default:"8585"`
}

type properties struct {
	at.ConfigurationProperties `value:"tcp"`

	Server Server `json:"server"`
	Client Client `json:"client"`
}
