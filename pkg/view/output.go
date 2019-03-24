package view

import (
	"github.com/lcserny/go-videosmover/pkg/models"
	"net/http"
	"strings"
)

type AjaxOutputController struct {
	config *models.WebviewConfig
}

func NewOutputController(config *models.WebviewConfig) *AjaxOutputController {
	return &AjaxOutputController{config: config}
}

func (sc *AjaxOutputController) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if strings.ToUpper(req.Method) == http.MethodPost {
		// TODO: write to the response the list of names
	}
}
