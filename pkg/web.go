package core

import "net/http"

type WebAPIRequest struct {
	Action  string      `json:"action"`
	Payload interface{} `json:"payload"`
}

type WebTemplateController interface {
	ServeTemplate(http.ResponseWriter, *http.Request) (string, interface{}, bool)
}

type WebAPIRequester interface {
	Return500(tmpl string, err error, resp http.ResponseWriter) (string, interface{}, bool)
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
