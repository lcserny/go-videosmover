package output

import (
	"github.com/lcserny/goutils"
	"net/http"
	"strconv"
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

func (c ajaxController) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")
	if strings.ToUpper(req.Method) == http.MethodPost {
		// TODO: validate request
		reqDataType := req.FormValue("type")
		diskPath := action.GetDiskPath(reqDataType, c.config)

		reqDataSkipCache := req.FormValue("skipcache")
		skipCache, err := strconv.ParseBool(reqDataSkipCache)
		if err != nil {
			goutils.LogError(err)
			resp.WriteHeader(http.StatusInternalServerError)
			return
		}

		reqDataSkipOnlineSearch := req.FormValue("skiponlinesearch")
		skipOnlineSearch, err := strconv.ParseBool(reqDataSkipOnlineSearch)
		if err != nil {
			goutils.LogError(err)
			resp.WriteHeader(http.StatusInternalServerError)
			return
		}

		reqDataName := req.FormValue("name")

		jsonBody, err := c.apiRequester.ExecutePOST("output", &RequestData{
			Name:             reqDataName,
			Type:             reqDataType,
			SkipCache:        skipCache,
			SkipOnlineSearch: skipOnlineSearch,
			DiskPath:         diskPath,
		}, c.config.VideosMoverAPI)
		if err != nil {
			goutils.LogError(err)
			resp.WriteHeader(http.StatusInternalServerError)
			return
		}

		var outputResponseData ResponseData
		if err = c.codec.Decode([]byte(jsonBody), &outputResponseData); err != nil {
			goutils.LogError(err)
			resp.WriteHeader(http.StatusInternalServerError)
			return
		}

		responseBytes, err := c.codec.EncodeBytes(outputResponseData)
		if err != nil {
			goutils.LogError(err)
			resp.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, err = resp.Write(responseBytes)
		goutils.LogError(err)
	}
}
