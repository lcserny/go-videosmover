package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/lcserny/go-videosmover/pkg/models"
	. "github.com/lcserny/goutils"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

type BinJsonExecuteHandler struct {
	videosMoverPath        string
	videosMoverConfigsPath string
}

func NewBinJsonExecuteHandler(serverConfig *models.ProxyServerConfig) *BinJsonExecuteHandler {
	return &BinJsonExecuteHandler{
		videosMoverPath:        serverConfig.PathVideosMoverBin,
		videosMoverConfigsPath: serverConfig.PathVideosMoverBinConfigs,
	}
}

func (h *BinJsonExecuteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch strings.ToUpper(r.Method) {
	case "GET":
		h.serveGET(w, r)
		break
	case "POST":
		h.servePOST(w, r)
		break
	default:
		LogWarning("Invalid http method. BinJsonExecuteHandler doesn't support: " + r.Method)
		return
	}
}

func (h *BinJsonExecuteHandler) serveGET(w http.ResponseWriter, r *http.Request) {
	LogInfo("Entered in GET request")

	time.Sleep(5 * time.Second)

	LogInfo("Exited GET request")
}

func (h *BinJsonExecuteHandler) servePOST(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)
	var jsonData RequestJsonData
	err := decoder.Decode(&jsonData)
	if err != nil {
		errorMessage := "Couldn't decode JSON data provided"
		_, _ = w.Write(getErrorJsonResponseAsBytes(errorMessage))
		LogErrorWithMessage(errorMessage, err)
		return
	}

	tempJsonFile := tmpStoreJsonPayload(jsonData.Payload)
	defer removeTmpStoredJsonPayload(tempJsonFile)

	var cmdOut bytes.Buffer
	var cmdErr bytes.Buffer
	cmd := exec.Command(h.videosMoverPath, "-configs="+h.videosMoverConfigsPath, "-action="+jsonData.Action, "-payloadFile="+tempJsonFile.Name())
	cmd.Stdout = &cmdOut
	cmd.Stderr = &cmdErr
	_ = cmd.Run()

	_, _ = w.Write(getJsonResponseFromAsBytes(string(cmdOut.Bytes()), string(cmdErr.Bytes())))
}
