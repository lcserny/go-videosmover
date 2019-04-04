package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/lcserny/goutils"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const TIME_FORMAT = "2006-01-02 15:04:05"

type TemplateController interface {
	ServeTemplate(http.ResponseWriter, *http.Request) (string, interface{}, bool)
}

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

func Return500Error(tmpl string, err error, resp http.ResponseWriter) (string, interface{}, bool) {
	resp.WriteHeader(http.StatusInternalServerError)
	goutils.LogError(err)
	return tmpl, nil, false
}

func ExecuteVideosMoverPOST(action string, payload interface{}, videosMoverAPI string) (string, error) {
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
	var jsonResp ResponseJsonData
	if err = json.Unmarshal(apiBody, &jsonResp); err != nil {
		return "", err
	}

	if apiResp.StatusCode != http.StatusOK || jsonResp.Code != http.StatusOK {
		return "", errors.New(jsonResp.Error)
	}

	return jsonResp.Body, nil
}

func getJsonResponseFromAsBytes(body, err string) []byte {
	if strings.Contains(body, "ERROR") {
		err = body
		body = ""
	}

	code := 200
	if len(err) > 0 {
		code = 500
	}

	responseJsonData := &ResponseJsonData{
		Code:  code,
		Error: err,
		Date:  time.Now().Format(TIME_FORMAT),
		Body:  body,
	}

	jsonBytes, _ := json.Marshal(responseJsonData)
	return jsonBytes
}

func getErrorJsonResponseAsBytes(err string) []byte {
	responseJsonData := &ResponseJsonData{
		Code:  500,
		Error: err,
		Date:  time.Now().Format(TIME_FORMAT),
		Body:  "",
	}

	jsonBytes, _ := json.Marshal(responseJsonData)
	return jsonBytes
}
