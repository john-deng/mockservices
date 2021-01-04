package tcp

import "hidevops.io/hiboot/pkg/at"

type Server struct {
	Port string `json:"port" default:"8585"`
}

type Properties struct {
	at.ConfigurationProperties `value:"tcp"`

	Server Server `json:"server"`
}
