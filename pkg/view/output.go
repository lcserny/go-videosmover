package view

import (
	"encoding/json"
	"github.com/lcserny/go-videosmover/pkg/models"
	. "github.com/lcserny/goutils"
	"net/http"
	"strconv"
	"strings"
)

type AjaxOutputController struct {
	config *models.WebviewConfig
}

func NewAjaxOutputController(config *models.WebviewConfig) *AjaxOutputController {
	return &AjaxOutputController{config: config}
}

func (sc *AjaxOutputController) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")
	if strings.ToUpper(req.Method) == http.MethodPost {
		// TODO: validate request
		reqDataType := req.FormValue("type")
		diskPath := sc.config.MoviesPath
		if strings.ToLower(reqDataType) != models.MOVIE {
			diskPath = sc.config.TvSeriesPath
		}

		reqDataSkipCache := req.FormValue("skipcache")
		skipCache, err := strconv.ParseBool(reqDataSkipCache)
		if err != nil {
			LogError(err)
			return
		}

		reqDataSkipOnlineSearch := req.FormValue("skiponlinesearch")
		skipOnlineSearch, err := strconv.ParseBool(reqDataSkipOnlineSearch)
		if err != nil {
			LogError(err)
			return
		}

		reqDataName := req.FormValue("name")

		jsonBody, err := executeVideosMoverPOST("output", &models.OutputRequestData{
			Name:             reqDataName,
			Type:             reqDataType,
			SkipCache:        skipCache,
			SkipOnlineSearch: skipOnlineSearch,
			DiskPath:         diskPath,
		}, sc.config.VideosMoverAPI)
		if err != nil {
			LogError(err)
			return
		}

		var outputResponseData models.OutputResponseData
		if err = json.Unmarshal([]byte(jsonBody), &outputResponseData); err != nil {
			LogError(err)
			return
		}

		responseBytes, err := json.Marshal(outputResponseData)
		if err != nil {
			LogError(err)
			return
		}

		_, err = resp.Write(responseBytes)
		LogError(err)
	}
}
