package output

import (
	"github.com/lcserny/goutils"
	"net/http"
	"strconv"
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

		jsonBody, err := web.ExecuteVideosMoverPOST("output", &RequestData{
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
