package web

import (
	"bytes"
	"github.com/lcserny/goutils"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
	"videosmover/pkg/fs"
	"videosmover/pkg/json"
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

	bodyBytes, _ := ioutil.ReadAll(r.Body)
	var requestData RequestData
	if err := json.Decode(bodyBytes, &requestData); err != nil {
		errorMessage := "Couldn't decode JSON data provided"
		_, _ = w.Write(getErrorResponseAsBytes(errorMessage))
		goutils.LogErrorWithMessage(errorMessage, err)
		return
	}

	encodeBytes, _ := json.EncodeBytes(requestData.Payload)
	tmpFile := fs.TmpStorePayload(encodeBytes)
	defer fs.RemoveTmpStoredPayload(tmpFile)

	var cmdOut bytes.Buffer
	var cmdErr bytes.Buffer
	cmd := exec.Command(h.Cmd.Path, "-configs="+h.Cmd.ConfigPath, "-action="+requestData.Action, "-payloadFile="+tmpFile.Name())
	cmd.Stdout = &cmdOut
	cmd.Stderr = &cmdErr
	if err := cmd.Run(); err != nil {
		cmdErr.WriteString(err.Error())
	}

	w.Write(getResponseFromAsBytes(string(cmdOut.Bytes()), string(cmdErr.Bytes())))
}
