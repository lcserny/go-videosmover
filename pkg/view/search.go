package view

import "net/http"

type SearchController struct {
}

func (sc *SearchController) ServeTemplate(resp http.ResponseWriter, req *http.Request) (name string, data interface{}, render bool) {
	resp.Header().Set("Content-Type", "text/html; charset=utf-8")
	resp.WriteHeader(http.StatusOK)

	if req.Method == http.MethodPost {
		return sc.POST(resp, req)
	}
	return sc.GET(resp, req)
}

func (sc *SearchController) GET(resp http.ResponseWriter, req *http.Request) (name string, data interface{}, render bool) {

	return "search", nil, true
}

func (sc *SearchController) POST(resp http.ResponseWriter, req *http.Request) (name string, data interface{}, render bool) {

	return "search", nil, true
}
