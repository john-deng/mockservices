package model

type ResponseData struct {
	Protocol         string         `json:"protocol"`
	Url              string         `json:"url"`
	App              string         `json:"app"`
	Version          string         `json:"version"`
	SourceApp        string         `json:"source_app"`
	SourceAppVersion string         `json:"source_app_version"`
	Cluster          string         `json:"cluster"`
	UserData         string         `json:"user_data"`
	MetaData         string         `json:"meta_data"`
	Upstream         []*GetResponse `json:"upstream"`
}

type GetResponse struct {
	Code    int          `json:"code"`
	Message string       `json:"message"`
	Data    ResponseData `json:"data"`
}
