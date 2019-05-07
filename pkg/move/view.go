package move

import (
	"github.com/lcserny/goutils"
	"net/http"
	"strings"
	"videosmover/pkg/action"
	"videosmover/pkg/json"
	"videosmover/pkg/web"
)

type AjaxController struct {
	config *web.WebviewConfig
}

func NewAjaxController(config *web.WebviewConfig) *AjaxController {
	return &AjaxController{config: config}
}

func (sc *AjaxController) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")
	if strings.ToUpper(req.Method) == http.MethodPost {
		moveJsData := req.FormValue("movedata")
		var moveReqDataList []RequestData
		if err := json.Decode([]byte(moveJsData), &moveReqDataList); err != nil {
			goutils.LogError(err)
			resp.WriteHeader(http.StatusInternalServerError)
			return
		}

		for i, ele := range moveReqDataList {
			ele.DiskPath = action.GetDiskPath(ele.Type, sc.config)
			moveReqDataList[i] = ele
		}

		jsonBody, err := web.ExecuteVideosMoverPOST("move", &moveReqDataList, sc.config.VideosMoverAPI)
		if err != nil {
			goutils.LogError(err)
			resp.WriteHeader(http.StatusInternalServerError)
			return
		}

		var moveResponseDataList []ResponseData
		if err = json.Decode([]byte(jsonBody), &moveResponseDataList); err != nil {
			goutils.LogError(err)
			resp.WriteHeader(http.StatusInternalServerError)
			return
		}

		responseBytes, err := json.EncodeBytes(moveResponseDataList)
		if err != nil {
			goutils.LogError(err)
			resp.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, err = resp.Write(responseBytes)
		goutils.LogError(err)
	}
}
