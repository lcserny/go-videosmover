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

// TODO: when finished, change config for windows
type javaExecutor struct {
	cmd   *core.CmdHandlerConfig
	codec core.Codec
}

func NewJavaExecutor(cmd *core.CmdHandlerConfig, codec core.Codec) http.Handler {
	return &javaExecutor{cmd: cmd, codec: codec}
}

func (be javaExecutor) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch strings.ToUpper(r.Method) {
	case "POST":
		be.servePOST(w, r)
		break
	default:
		goutils.LogWarning("Invalid http method. JavaJsonExecuteHandler doesn't support: " + r.Method)
		return
	}
}

func (be javaExecutor) servePOST(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	bodyBytes, _ := ioutil.ReadAll(r.Body)
	var requestData RequestData
	if err := be.codec.Decode(bodyBytes, &requestData); err != nil {
		errorMessage := "Couldn't decode JSON data provided"
		_, _ = w.Write([]byte(errorMessage))
		goutils.LogErrorWithMessage(errorMessage, err)
		return
	}

	jsonPayload, _ := be.codec.EncodeBytes(requestData.Payload)
	base64Json := base64.StdEncoding.EncodeToString([]byte(jsonPayload))

	var cmdOut bytes.Buffer
	var cmdErr bytes.Buffer
	cmd := exec.Command(be.cmd.Path, "-config="+be.cmd.ConfigPath, "-action="+requestData.Action, "-payload="+base64Json)
	cmd.Stdout = &cmdOut
	cmd.Stderr = &cmdErr
	if err := cmd.Run(); err != nil {
		cmdErr.WriteString(err.Error())
	}

	//w.Write([]byte(string(cmdOut.Bytes()), string(cmdErr.Bytes())))
}
