package move

import (
	"github.com/lcserny/goutils"
	"net/http"
	"strings"
	"videosmover/pkg"
	"videosmover/pkg/action"
	"videosmover/pkg/web"
)

type ajaxController struct {
	config *web.WebviewConfig
	codec  core.Codec
}

func NewAjaxController(config *web.WebviewConfig, codec core.Codec) http.Handler {
	return &ajaxController{config: config, codec: codec}
}

func (c *ajaxController) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")
	if strings.ToUpper(req.Method) == http.MethodPost {
		moveJsData := req.FormValue("movedata")
		var moveReqDataList []RequestData
		if err := c.codec.Decode([]byte(moveJsData), &moveReqDataList); err != nil {
			goutils.LogError(err)
			resp.WriteHeader(http.StatusInternalServerError)
			return
		}

		for i, ele := range moveReqDataList {
			ele.DiskPath = action.GetDiskPath(ele.Type, c.config)
			moveReqDataList[i] = ele
		}

		jsonBody, err := web.ExecuteVideosMoverPOST("move", &moveReqDataList, c.config.VideosMoverAPI)
		if err != nil {
			goutils.LogError(err)
			resp.WriteHeader(http.StatusInternalServerError)
			return
		}

		var moveResponseDataList []ResponseData
		if err = c.codec.Decode([]byte(jsonBody), &moveResponseDataList); err != nil {
			goutils.LogError(err)
			resp.WriteHeader(http.StatusInternalServerError)
			return
		}

		responseBytes, err := c.codec.EncodeBytes(moveResponseDataList)
		if err != nil {
			goutils.LogError(err)
			resp.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, err = resp.Write(responseBytes)
		goutils.LogError(err)
	}
}
