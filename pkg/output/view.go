package output

import (
	"encoding/json"
	"github.com/lcserny/go-videosmover/pkg/action"
	"github.com/lcserny/go-videosmover/pkg/convert"
	inhttp "github.com/lcserny/go-videosmover/pkg/http"
	. "github.com/lcserny/goutils"
	"net/http"
	"strconv"
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
		// TODO: validate request
		reqDataType := req.FormValue("type")
		diskPath := action.GetDiskPath(reqDataType, sc.config)

		reqDataSkipCache := req.FormValue("skipcache")
		skipCache, err := strconv.ParseBool(reqDataSkipCache)
		if err != nil {
			LogError(err)
			resp.WriteHeader(http.StatusInternalServerError)
			return
		}

		reqDataSkipOnlineSearch := req.FormValue("skiponlinesearch")
		skipOnlineSearch, err := strconv.ParseBool(reqDataSkipOnlineSearch)
		if err != nil {
			LogError(err)
			resp.WriteHeader(http.StatusInternalServerError)
			return
		}

		reqDataName := req.FormValue("name")

		jsonBody, err := inhttp.ExecuteVideosMoverPOST("output", &convert.OutputRequestData{
			Name:             reqDataName,
			Type:             reqDataType,
			SkipCache:        skipCache,
			SkipOnlineSearch: skipOnlineSearch,
			DiskPath:         diskPath,
		}, sc.config.VideosMoverAPI)
		if err != nil {
			LogError(err)
			resp.WriteHeader(http.StatusInternalServerError)
			return
		}

		var outputResponseData convert.OutputResponseData
		if err = json.Unmarshal([]byte(jsonBody), &outputResponseData); err != nil {
			LogError(err)
			resp.WriteHeader(http.StatusInternalServerError)
			return
		}

		responseBytes, err := json.Marshal(outputResponseData)
		if err != nil {
			LogError(err)
			resp.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, err = resp.Write(responseBytes)
		LogError(err)
	}
}
