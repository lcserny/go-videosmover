package view

import (
	"encoding/json"
	"github.com/lcserny/go-videosmover/pkg/models"
	. "github.com/lcserny/goutils"
	"net/http"
	"strings"
)

type AjaxMoveController struct {
	config *models.WebviewConfig
}

func NewAjaxMoveController(config *models.WebviewConfig) *AjaxMoveController {
	return &AjaxMoveController{config: config}
}

func (sc *AjaxMoveController) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")
	if strings.ToUpper(req.Method) == http.MethodPost {
		moveJsData := req.FormValue("movedata")
		var moveReqDataList []models.MoveRequestData
		if err := json.Unmarshal([]byte(moveJsData), &moveReqDataList); err != nil {
			LogError(err)
			resp.WriteHeader(http.StatusInternalServerError)
			return
		}

		for i, ele := range moveReqDataList {
			ele.DiskPath = getDiskPath(ele.Type, sc.config)
			moveReqDataList[i] = ele
		}

		jsonBody, err := executeVideosMoverPOST("move", &moveReqDataList, sc.config.VideosMoverAPI)
		if err != nil {
			LogError(err)
			resp.WriteHeader(http.StatusInternalServerError)
			return
		}

		var moveResponseDataList []models.MoveResponseData
		if err = json.Unmarshal([]byte(jsonBody), &moveResponseDataList); err != nil {
			LogError(err)
			resp.WriteHeader(http.StatusInternalServerError)
			return
		}

		responseBytes, err := json.Marshal(moveResponseDataList)
		if err != nil {
			LogError(err)
			resp.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, err = resp.Write(responseBytes)
		LogError(err)
	}
}
