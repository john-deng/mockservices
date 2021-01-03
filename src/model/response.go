package model

import "net/http"

type ResponseData struct {
	Protocol 		 string
	Url              string
	App              string
	Version          string
	SourceApp        string
	SourceAppVersion string
	Cluster          string
	UserData         string
	MetaData         string
	//UserAgent ua.UserAgent
	Header   http.Header
	Upstream []*GetResponse
}

type GetResponse struct {
	Code    int
	Message string
	Data    ResponseData
}
