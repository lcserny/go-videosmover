package view

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/lcserny/go-videosmover/pkg/handlers"
	"github.com/lcserny/go-videosmover/pkg/models"
	"io/ioutil"
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
	bodyReq, err := generateActionRequest("search", &models.SearchRequestData{Path: sc.config.DownloadsPath})
	if err != nil {
		return return500Error("search", err, resp)
	}

	apiResp, err := http.Post(sc.config.VideosMoverAPI, "application/json", bytes.NewBufferString(bodyReq))
	if err != nil {
		return return500Error("search", err, resp)
	}
	defer apiResp.Body.Close()

	body, _ := ioutil.ReadAll(apiResp.Body)
	if apiResp.StatusCode != http.StatusOK {
		return return500Error("search", errors.New(string(body)), resp)
	}

	var jsonResp handlers.ResponseJsonData
	if err = json.Unmarshal(body, &jsonResp); err != nil {
		return return500Error("search", err, resp)
	}

	var searchResponseDataList []models.SearchResponseData
	if err = json.Unmarshal([]byte(jsonResp.Body), &searchResponseDataList); err != nil {
		return return500Error("search", err, resp)
	}

	pageData := SearchResultPageData{}
	for _, data := range searchResponseDataList {
		pageData.Videos = append(pageData.Videos, data.Path)
	}

	resp.WriteHeader(http.StatusOK)
	return "search", pageData, true
}
