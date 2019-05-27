package search

import (
	"net/http"
	"path/filepath"
	"strings"
	"videosmover/pkg"
)

type VideoResult struct {
	Index     int
	Name      string
	FileName  string
	VideoPath string
	Subtitles []string
}

type ResultWrapper struct {
	Video        VideoResult
	EncodedVideo string
}

type ResultPageData struct {
	Videos []ResultWrapper
}

type controller struct {
	config             *core.WebviewConfig
	codec              core.Codec
	apiReqResProcessor core.WebApiReqResProcessor
}

func NewController(cfg *core.WebviewConfig, codec core.Codec, apiReqResProcessor core.WebApiReqResProcessor) core.WebTemplateController {
	return &controller{config: cfg, codec: codec, apiReqResProcessor: apiReqResProcessor}
}

func (c *controller) ServeTemplate(resp http.ResponseWriter, req *http.Request) (name string, data core.WebTemplateData, render bool) {
	resp.Header().Set("Content-Type", "text/html; charset=utf-8")

	switch strings.ToUpper(req.Method) {
	case http.MethodPost:
		return c.POST(resp, req)
	default:
		return c.GET(resp, req)
	}
}

func (c controller) GET(resp http.ResponseWriter, req *http.Request) (name string, data core.WebTemplateData, render bool) {
	resp.WriteHeader(http.StatusOK)
	return "search", core.WebTemplateData{}, true
}

func (c controller) POST(resp http.ResponseWriter, req *http.Request) (name string, data core.WebTemplateData, render bool) {
	jsonBody, err := c.apiReqResProcessor.ExecutePOST("search", &RequestData{Path: c.config.DownloadsPath}, c.config.VideosMoverAPI)
	if err != nil {
		return c.apiReqResProcessor.Return500("search", err, resp)
	}

	if jsonBody == "" {
		return c.GET(resp, req)
	}

	var searchResponseDataList []core.VideoSearchResult
	if err = c.codec.Decode([]byte(jsonBody), &searchResponseDataList); err != nil {
		return c.apiReqResProcessor.Return500("search", err, resp)
	}

	pageData := ResultPageData{}
	for i, data := range searchResponseDataList {
		fileName := filepath.Base(data.Path)
		fileDir := filepath.Dir(data.Path)
		name := filepath.Base(fileDir)
		if fileDir == filepath.Clean(c.config.DownloadsPath) {
			name = fileName
		}

		searchResult := VideoResult{
			Index:     i,
			Name:      name,
			FileName:  fileName,
			VideoPath: data.Path,
			Subtitles: data.Subtitles,
		}
		encodedContent, _ := c.codec.EncodeString(searchResult)

		pageData.Videos = append(pageData.Videos, ResultWrapper{Video: searchResult, EncodedVideo: encodedContent})
	}

	resp.WriteHeader(http.StatusOK)
	return "search", core.WebTemplateData{Data: pageData}, true
}
