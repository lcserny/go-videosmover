package web

import (
	"bytes"
	"errors"
	"github.com/lcserny/goutils"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	"videosmover/pkg"
)

type CloudDatabase interface {
	Init()
}

type apiRequester struct {
	timeFormat string
	codec      core.Codec
}

func NewApiRequester(codec core.Codec) core.WebApiReqResProcessor {
	return &apiRequester{timeFormat: "2006-01-02 15:04:05", codec: codec}
}

func (ar apiRequester) generateActionRequest(action string, payload interface{}) (string, error) {
	apiRequest := core.WebAPIRequest{Action: action, Payload: payload}
	s, err := ar.codec.EncodeString(apiRequest)
	if err != nil {
		return "", err
	}
	return s, nil
}

func (ar apiRequester) Return500(tmpl string, err error, resp http.ResponseWriter) (string, core.WebTemplateData, bool) {
	resp.WriteHeader(http.StatusInternalServerError)
	goutils.LogError(err)
	return tmpl, core.WebTemplateData{}, false
}

func (ar apiRequester) ExecutePOST(action string, payload interface{}, videosMoverAPI string) (string, error) {
	apiReq, err := ar.generateActionRequest(action, payload)
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
	if err = ar.codec.Decode(apiBody, &responseData); err != nil {
		return "", err
	}

	if apiResp.StatusCode != http.StatusOK || responseData.Code != http.StatusOK {
		return "", errors.New(responseData.Error)
	}

	return responseData.Body, nil
}

func (ar apiRequester) ProcessBody(content, err string) []byte {
	if strings.Contains(content, "ERROR") {
		err = content
		content = ""
	}

	code := 200
	if len(err) > 0 {
		code = 500
	}

	responseData := &ResponseData{
		Code:  code,
		Error: err,
		Date:  time.Now().Format(ar.timeFormat),
		Body:  content,
	}

	respBytes, _ := ar.codec.EncodeBytes(responseData)
	return respBytes
}

func (ar apiRequester) ProcessError(err string) []byte {
	responseData := &ResponseData{
		Code:  500,
		Error: err,
		Date:  time.Now().Format(ar.timeFormat),
		Body:  "",
	}

	respBytes, _ := ar.codec.EncodeBytes(responseData)
	return respBytes
}
