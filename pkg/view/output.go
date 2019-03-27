package view

import (
	"encoding/json"
	"github.com/lcserny/go-videosmover/pkg/models"
	. "github.com/lcserny/goutils"
	"net/http"
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
		// TODO: validate request params! and write to the response the list of names
		formBytes, err := json.Marshal(req.FormValue("data"))
		LogError(err)
		resp.Write(formBytes)
	}
}
