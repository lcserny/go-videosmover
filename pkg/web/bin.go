package web

import (
	"bytes"
	"encoding/json"
	"github.com/lcserny/goutils"
	"net/http"
	"os/exec"
	"strings"
	"videosmover/pkg/convert"
)

type BinJsonExecuteHandler struct {
	Cmd *cmdHandlerConfig
}

func (h *BinJsonExecuteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch strings.ToUpper(r.Method) {
	case "POST":
		h.servePOST(w, r)
		break
	default:
		goutils.LogWarning("Invalid http method. BinJsonExecuteHandler doesn't support: " + r.Method)
		return
	}
}

func (h *BinJsonExecuteHandler) servePOST(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)
	var jsonData RequestJsonData
	err := decoder.Decode(&jsonData)
	if err != nil {
		errorMessage := "Couldn't decode JSON data provided"
		_, _ = w.Write(getErrorJsonResponseAsBytes(errorMessage))
		goutils.LogErrorWithMessage(errorMessage, err)
		return
	}

	tempJsonFile := convert.TmpStoreJsonPayload(jsonData.Payload)
	defer convert.RemoveTmpStoredJsonPayload(tempJsonFile)

	var cmdOut bytes.Buffer
	var cmdErr bytes.Buffer
	cmd := exec.Command(h.Cmd.Path, "-configs="+h.Cmd.ConfigPath, "-action="+jsonData.Action, "-payloadFile="+tempJsonFile.Name())
	cmd.Stdout = &cmdOut
	cmd.Stderr = &cmdErr
	err = cmd.Run()
	if err != nil {
		cmdErr.WriteString(err.Error())
	}

	w.Write(getJsonResponseFromAsBytes(string(cmdOut.Bytes()), string(cmdErr.Bytes())))
}
