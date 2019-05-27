package core

import "net/http"

type WebAPIRequest struct {
	Action  string      `json:"action"`
	Payload interface{} `json:"payload"`
}

type WebTemplateData struct {
	Data     interface{}
	DarkMode bool
}

type WebTemplateController interface {
	ServeTemplate(http.ResponseWriter, *http.Request) (string, WebTemplateData, bool)
}

type WebAPIRequester interface {
	Return500(tmpl string, err error, resp http.ResponseWriter) (string, WebTemplateData, bool)
	ExecutePOST(action string, payload interface{}, videosMoverAPI string) (string, error)
}

type WebResponseProcessor interface {
	ProcessBody(content, err string) []byte
	ProcessError(err string) []byte
}

type WebApiReqResProcessor interface {
	WebAPIRequester
	WebResponseProcessor
}

type VideoWebResult struct {
	Title       string
	ReleaseDate string
}

type VideoWebSearcher interface {
	CanSearch() bool
	SearchMovies(name string, year int) ([]*VideoWebResult, error)
	SearchTVSeries(name string, year int) ([]*VideoWebResult, error)
}
