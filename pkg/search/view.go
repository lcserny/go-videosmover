package search

import (
	"net/http"
	"path/filepath"
	"strings"
	"videosmover/pkg"
	"videosmover/pkg/web"
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

type controller struct {
	config *web.WebviewConfig
	codec  core.Codec
}

func NewController(cfg *web.WebviewConfig, codec core.Codec) web.TemplateController {
	return &controller{config: cfg, codec: codec}
}

func (c *controller) ServeTemplate(resp http.ResponseWriter, req *http.Request) (name string, data interface{}, render bool) {
	resp.Header().Set("Content-Type", "text/html; charset=utf-8")

	switch strings.ToUpper(req.Method) {
	case http.MethodPost:
		return c.POST(resp, req)
	default:
		return c.GET(resp, req)
	}
}

func (c *controller) GET(resp http.ResponseWriter, req *http.Request) (name string, data interface{}, render bool) {
	resp.WriteHeader(http.StatusOK)
	return "search", nil, true
}

func (c *controller) POST(resp http.ResponseWriter, req *http.Request) (name string, data interface{}, render bool) {
	jsonBody, err := web.ExecuteVideosMoverPOST("search", &RequestData{Path: c.config.DownloadsPath}, c.config.VideosMoverAPI)
	if err != nil {
		return web.Return500Error("search", err, resp)
	}

	if jsonBody == "" {
		return c.GET(resp, req)
	}

	var searchResponseDataList []ResponseData
	if err = c.codec.Decode([]byte(jsonBody), &searchResponseDataList); err != nil {
		return web.Return500Error("search", err, resp)
	}

	pageData := ResultPageData{}
	for i, data := range searchResponseDataList {
		fileName := filepath.Base(data.Path)
		fileDir := filepath.Dir(data.Path)
		name := filepath.Base(fileDir)
		if fileDir == filepath.Clean(c.config.DownloadsPath) {
			name = fileName
		}

		encodeString, _ := c.codec.EncodeString(data.Subtitles)

		searchResult := Result{
			Index:            i,
			Name:             name,
			FileName:         fileName,
			VideoPath:        data.Path,
			EncodedSubsArray: encodeString,
		}
		pageData.Videos = append(pageData.Videos, searchResult)
	}

	resp.WriteHeader(http.StatusOK)
	return "search", pageData, true
}
