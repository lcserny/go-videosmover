package web

import (
	"bytes"
	"fmt"
	"github.com/lcserny/goutils"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"os"
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

	encodeBytes, _ := be.codec.EncodeBytes(requestData.Payload)
	tmpFile := be.tmpStorePayload(encodeBytes)
	defer be.removeTmpStoredPayload(tmpFile)

	var cmdOut bytes.Buffer
	var cmdErr bytes.Buffer
	cmd := exec.Command(be.cmd.Path, "-config="+be.cmd.ConfigPath, "-action="+requestData.Action, "-payloadFile="+tmpFile.Name())
	cmd.Stdout = &cmdOut
	cmd.Stderr = &cmdErr
	if err := cmd.Run(); err != nil {
		cmdErr.WriteString(err.Error())
	}

	//w.Write([]byte(string(cmdOut.Bytes()), string(cmdErr.Bytes())))
}

func (be javaExecutor) removeTmpStoredPayload(tempFile *os.File) {
	err := tempFile.Close()
	goutils.LogErrorWithMessage(fmt.Sprintf("Couldn't close tmpFile: %s", tempFile.Name()), err)

	err = os.Remove(tempFile.Name())
	goutils.LogErrorWithMessage(fmt.Sprintf("Couldn't remove tmpFile: %s", tempFile.Name()), err)
}

func (be javaExecutor) tmpStorePayload(bytes []byte) *os.File {
	tempFile, err := ioutil.TempFile(os.TempDir(), "vms-")
	goutils.LogError(errors.Wrap(err, "Couldn't create tmpFile"))

	_, err = tempFile.Write(bytes)
	goutils.LogErrorWithMessage(fmt.Sprintf("Couldn't write to tmpFile: %s", tempFile.Name()), err)

	return tempFile
}
