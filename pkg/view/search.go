package view

import (
	"github.com/lcserny/go-videosmover/pkg/models"
	"net/http"
	"strings"
)

type SearchResultPageData struct {
	Videos []string
}

type SearchController struct {
	config *models.WebviewConfig
}

func NewSearchController(config *models.WebviewConfig) *SearchController {
	return &SearchController{config: config}
}

func (sc *SearchController) ServeTemplate(resp http.ResponseWriter, req *http.Request) (name string, data interface{}, render bool) {
	resp.Header().Set("Content-Type", "text/html; charset=utf-8")

	switch strings.ToUpper(req.Method) {
	case http.MethodPost:
		return sc.POST(resp, req)
	default:
		return sc.GET(resp, req)
	}
}

func (sc *SearchController) GET(resp http.ResponseWriter, req *http.Request) (name string, data interface{}, render bool) {
	resp.WriteHeader(http.StatusOK)
	return "search", nil, true
}

func (sc *SearchController) POST(resp http.ResponseWriter, req *http.Request) (name string, data interface{}, render bool) {
	resp.WriteHeader(http.StatusOK)
	// TODO: exec POST on `cfg.videosMoverAPI` passing needed data
	var videos []string
	videos = append(videos, "/some/path/here.mp4")
	return "search", SearchResultPageData{Videos: videos}, true
}
