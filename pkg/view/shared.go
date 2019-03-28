package view

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/lcserny/go-videosmover/pkg/handlers"
	"github.com/lcserny/go-videosmover/pkg/models"
	. "github.com/lcserny/goutils"
	"io/ioutil"
	"net/http"
	"strings"
)

type VideosMoverAPIRequest struct {
	Action  string      `json:"action"`
	Payload interface{} `json:"payload"`
}

func generateActionRequest(action string, payload interface{}) (string, error) {
	apiRequest := VideosMoverAPIRequest{Action: action, Payload: payload}
	jsonBytes, err := json.Marshal(apiRequest)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

func return500Error(tmpl string, err error, resp http.ResponseWriter) (string, interface{}, bool) {
	resp.WriteHeader(http.StatusInternalServerError)
	LogError(err)
	return tmpl, nil, false
}

func executeVideosMoverPOST(action string, payload interface{}, videosMoverAPI string) (string, error) {
	apiReq, err := generateActionRequest(action, payload)
	if err != nil {
		return "", err
	}

	apiResp, err := http.Post(videosMoverAPI, "application/json", bytes.NewBufferString(apiReq))
	if err != nil {
		return "", err
	}
	defer apiResp.Body.Close()

	apiBody, _ := ioutil.ReadAll(apiResp.Body)
	var jsonResp handlers.ResponseJsonData
	if err = json.Unmarshal(apiBody, &jsonResp); err != nil {
		return "", err
	}

	if apiResp.StatusCode != http.StatusOK || jsonResp.Code != http.StatusOK {
		return "", errors.New(jsonResp.Error)
	}

	return jsonResp.Body, nil
}

func getDiskPath(videoType string, config *models.WebviewConfig) string {
	loweredType := strings.ToLower(videoType)
	if loweredType != models.UNKNOWN {
		diskPath := config.MoviesPath
		if loweredType != models.MOVIE {
			diskPath = config.TvSeriesPath
		}
		return diskPath
	}
	return ""
}
