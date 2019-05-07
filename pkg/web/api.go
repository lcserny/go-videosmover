package web

import (
	"bytes"
	"errors"
	"github.com/lcserny/goutils"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	"videosmover/pkg/json"
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
	s, err := json.EncodeString(apiRequest)
	if err != nil {
		return "", err
	}
	return s, nil
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
	var responseData ResponseData
	if err = json.Decode(apiBody, &responseData); err != nil {
		return "", err
	}

	if apiResp.StatusCode != http.StatusOK || responseData.Code != http.StatusOK {
		return "", errors.New(responseData.Error)
	}

	return responseData.Body, nil
}

func getResponseFromAsBytes(body, err string) []byte {
	if strings.Contains(body, "ERROR") {
		err = body
		body = ""
	}

	code := 200
	if len(err) > 0 {
		code = 500
	}

	responseData := &ResponseData{
		Code:  code,
		Error: err,
		Date:  time.Now().Format(TIME_FORMAT),
		Body:  body,
	}

	respBytes, _ := json.EncodeBytes(responseData)
	return respBytes
}

func getErrorResponseAsBytes(err string) []byte {
	responseData := &ResponseData{
		Code:  500,
		Error: err,
		Date:  time.Now().Format(TIME_FORMAT),
		Body:  "",
	}

	respBytes, _ := json.EncodeBytes(responseData)
	return respBytes
}
