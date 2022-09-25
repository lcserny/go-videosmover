package move

import (
	"github.com/lcserny/goutils"
	"io/ioutil"
	"net/http"
	"strings"
	"videosmover/pkg"
	"videosmover/pkg/action"
)

type ajaxController struct {
	config       *core.WebviewConfig
	codec        core.Codec
	apiRequester core.WebAPIRequester
}

func NewAjaxController(config *core.WebviewConfig, codec core.Codec, apiRequester core.WebAPIRequester) http.Handler {
	return &ajaxController{config: config, codec: codec, apiRequester: apiRequester}
}

func (c *ajaxController) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")
	if strings.ToUpper(req.Method) == http.MethodPost {
		reqBytes, err := ioutil.ReadAll(req.Body)
		if err != nil {
			goutils.LogError(err)
			resp.WriteHeader(http.StatusInternalServerError)
			return
		}

		var moveReqDataList []RequestData
		if err := c.codec.Decode(reqBytes, &moveReqDataList); err != nil {
			goutils.LogError(err)
			resp.WriteHeader(http.StatusInternalServerError)
			return
		}

		for i, ele := range moveReqDataList {
			ele.DiskPath = action.GetDiskPath(ele.Type, c.config)
			moveReqDataList[i] = ele
		}

		jsonBody, err := c.apiRequester.ExecutePOST("move", &moveReqDataList, c.config.VideosMoverAPI)
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
