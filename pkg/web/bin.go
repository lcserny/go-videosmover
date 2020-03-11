package web

import (
	"bytes"
	"encoding/base64"
	"github.com/lcserny/goutils"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
	"videosmover/pkg"
)

type binExecutor struct {
	cmd   *core.CmdHandlerConfig
	codec core.Codec
	wrp   core.WebResponseProcessor
}

func NewBinExecutor(cmd *core.CmdHandlerConfig, codec core.Codec, processor core.WebResponseProcessor) http.Handler {
	return &binExecutor{cmd: cmd, codec: codec, wrp: processor}
}

func (be binExecutor) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch strings.ToUpper(r.Method) {
	case "POST":
		be.servePOST(w, r)
		break
	default:
		goutils.LogWarning("Invalid http method. BinJsonExecuteHandler doesn't support: " + r.Method)
		return
	}
}

func (be binExecutor) servePOST(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	bodyBytes, _ := ioutil.ReadAll(r.Body)
	var requestData RequestData
	if err := be.codec.Decode(bodyBytes, &requestData); err != nil {
		errorMessage := "Couldn't decode JSON data provided"
		_, _ = w.Write(be.wrp.ProcessError(errorMessage))
		goutils.LogErrorWithMessage(errorMessage, err)
		return
	}

	jsonPayload, _ := be.codec.EncodeString(requestData.Payload)
	base64Json := base64.StdEncoding.EncodeToString([]byte(jsonPayload))

	var cmdOut bytes.Buffer
	var cmdErr bytes.Buffer
	cmd := exec.Command(be.cmd.Path, "-config="+be.cmd.ConfigPath, "-action="+requestData.Action, "-payload="+base64Json)
	cmd.Stdout = &cmdOut
	cmd.Stderr = &cmdErr
	if err := cmd.Run(); err != nil {
		cmdErr.WriteString(err.Error())
	}

	w.Write(be.wrp.ProcessBody(string(cmdOut.Bytes()), string(cmdErr.Bytes())))
}
