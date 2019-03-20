package view

import "net/http"

// TODO: if you want to encapsulate better, make a handler class and use interface in `templatedViewsMap`
func SearchController(writer http.ResponseWriter, request *http.Request) (tmplName string, tmplData interface{}, renderTmpl bool) {
	return "search", nil, true
}

func SearchResultsController(writer http.ResponseWriter, request *http.Request) (tmplName string, tmplData interface{}, renderTmpl bool) {
	return "search_results", nil, true
}
