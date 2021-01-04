package model

import (
	"net/http"
)

type Data struct {
	Protocol         string      `json:"protocol"`
	Url              string      `json:"url"`
	App              string      `json:"app"`
	Version          string      `json:"version"`
	SourceApp        string      `json:"source_app"`
	SourceAppVersion string      `json:"source_app_version"`
	Cluster          string      `json:"cluster"`
	UserData         string      `json:"user_data"`
	MetaData         string      `json:"meta_data"`
	Upstream         []*Response `json:"upstream"`
}

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    Data   `json:"data"`
}

type TcpRequest struct {
	Address string `json:"address"`
	Header http.Header `json:"header"`
}

type TcpResponse struct {
	*Response
	Header http.Header `json:"header"`
}