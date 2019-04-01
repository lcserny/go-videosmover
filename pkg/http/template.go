package http

import "net/http"

type TemplateController interface {
	ServeTemplate(http.ResponseWriter, *http.Request) (string, interface{}, bool)
}
