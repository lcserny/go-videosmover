package search

import (
	"encoding/json"
	"github.com/lcserny/go-videosmover/pkg/convert"
	"github.com/lcserny/go-videosmover/pkg/web"
	"net/http"
	"path/filepath"
	"strings"
)

type Result struct {
	Index            int
	Name             string
	FileName         string
	VideoPath        string
	EncodedSubsArray string
}

type ResultPageData struct {
	Videos []Result
}

type Controller struct {
	config *web.WebviewConfig
}

func NewController(config *web.WebviewConfig) *Controller {
	return &Controller{config: config}
}

func (sc *Controller) ServeTemplate(resp http.ResponseWriter, req *http.Request) (name string, data interface{}, render bool) {
	resp.Header().Set("Content-Type", "text/html; charset=utf-8")

	switch strings.ToUpper(req.Method) {
	case http.MethodPost:
		return sc.POST(resp, req)
	default:
		return sc.GET(resp, req)
	}
}

func (sc *Controller) GET(resp http.ResponseWriter, req *http.Request) (name string, data interface{}, render bool) {
	resp.WriteHeader(http.StatusOK)
	return "search", nil, true
}

func (sc *Controller) POST(resp http.ResponseWriter, req *http.Request) (name string, data interface{}, render bool) {
	jsonBody, err := web.ExecuteVideosMoverPOST("search", &RequestData{Path: sc.config.DownloadsPath}, sc.config.VideosMoverAPI)
	if err != nil {
		return web.Return500Error("search", err, resp)
	}

	if jsonBody == "" {
		return sc.GET(resp, req)
	}

	var searchResponseDataList []ResponseData
	if err = json.Unmarshal([]byte(jsonBody), &searchResponseDataList); err != nil {
		return web.Return500Error("search", err, resp)
	}

	pageData := ResultPageData{}
	for i, data := range searchResponseDataList {
		fileName := filepath.Base(data.Path)
		fileDir := filepath.Dir(data.Path)
		name := filepath.Base(fileDir)
		if fileDir == filepath.Clean(sc.config.DownloadsPath) {
			name = fileName
		}

		searchResult := Result{
			Index:            i,
			Name:             name,
			FileName:         fileName,
			VideoPath:        data.Path,
			EncodedSubsArray: convert.GetJSONEncodedString(data.Subtitles),
		}
		pageData.Videos = append(pageData.Videos, searchResult)
	}

	resp.WriteHeader(http.StatusOK)
	return "search", pageData, true
}
