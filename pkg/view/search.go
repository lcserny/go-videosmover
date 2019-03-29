package view

import (
	"encoding/json"
	"github.com/lcserny/go-videosmover/pkg/models"
	"net/http"
	"path/filepath"
	"strings"
)

type SearchResult struct {
	Index            int
	Name             string
	FileName         string
	VideoPath        string
	EncodedSubsArray string
}

type SearchResultPageData struct {
	Videos []SearchResult
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
	jsonBody, err := executeVideosMoverPOST("search", &models.SearchRequestData{Path: sc.config.DownloadsPath}, sc.config.VideosMoverAPI)
	if err != nil {
		return return500Error("search", err, resp)
	}

	if jsonBody == "" {
		return sc.GET(resp, req)
	}

	var searchResponseDataList []models.SearchResponseData
	if err = json.Unmarshal([]byte(jsonBody), &searchResponseDataList); err != nil {
		return return500Error("search", err, resp)
	}

	pageData := SearchResultPageData{}
	for i, data := range searchResponseDataList {
		fileName := filepath.Base(data.Path)
		fileDir := filepath.Dir(data.Path)
		name := filepath.Base(fileDir)
		if fileDir == filepath.Clean(sc.config.DownloadsPath) {
			name = fileName
		}

		searchResult := SearchResult{
			Index:            i,
			Name:             name,
			FileName:         fileName,
			VideoPath:        data.Path,
			EncodedSubsArray: encodeToJSONArray(data.Subtitles),
		}
		pageData.Videos = append(pageData.Videos, searchResult)
	}

	resp.WriteHeader(http.StatusOK)
	return "search", pageData, true
}
