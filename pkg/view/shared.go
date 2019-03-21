package view

import (
	"encoding/json"
	. "github.com/lcserny/goutils"
	"net/http"
)

type VideosMoverAPIRequest struct {
	Action  string      `json:"action"`
	Payload interface{} `json:"payload"`
}

func generateActionRequest(action string, payload interface{}) (string, error) {
	apiRequest := VideosMoverAPIRequest{Action: action, Payload: payload}
	bytes, err := json.Marshal(apiRequest)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func return500Error(tmpl string, err error, resp http.ResponseWriter) (string, interface{}, bool) {
	resp.WriteHeader(http.StatusInternalServerError)
	LogError(err)
	return tmpl, nil, false
}
