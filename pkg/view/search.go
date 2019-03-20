package view

import "net/http"

// TODO: if you wnat to encapsulate better, make a handler class and use interface in `templatedViewsMap`
func SearchController(writer http.ResponseWriter, request *http.Request) (tmplName string, tmplData interface{}) {
	return "search.gohtml", nil
}

func SearchResultsController(writer http.ResponseWriter, request *http.Request) (tmplName string, tmplData interface{}) {
	return "search_results.gohtml", nil
}
