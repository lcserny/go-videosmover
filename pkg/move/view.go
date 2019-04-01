package move

import (
	"encoding/json"
	"github.com/lcserny/go-videosmover/pkg/action"
	"github.com/lcserny/go-videosmover/pkg/convert"
	inhttp "github.com/lcserny/go-videosmover/pkg/http"
	. "github.com/lcserny/goutils"
	"net/http"
	"strings"
)

type AjaxController struct {
	config *convert.WebviewConfig
}

func NewAjaxController(config *convert.WebviewConfig) *AjaxController {
	return &AjaxController{config: config}
}

func (sc *AjaxController) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")
	if strings.ToUpper(req.Method) == http.MethodPost {
		moveJsData := req.FormValue("movedata")
		var moveReqDataList []convert.MoveRequestData
		if err := json.Unmarshal([]byte(moveJsData), &moveReqDataList); err != nil {
			LogError(err)
			resp.WriteHeader(http.StatusInternalServerError)
			return
		}

		for i, ele := range moveReqDataList {
			ele.DiskPath = action.GetDiskPath(ele.Type, sc.config)
			moveReqDataList[i] = ele
		}

		jsonBody, err := inhttp.ExecuteVideosMoverPOST("move", &moveReqDataList, sc.config.VideosMoverAPI)
		if err != nil {
			LogError(err)
			resp.WriteHeader(http.StatusInternalServerError)
			return
		}

		var moveResponseDataList []convert.MoveResponseData
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
