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

type binExecutor struct {
	cmd   *cmdHandlerConfig
	codec core.Codec
	wrp   core.WebResponseProcessor
}

func NewBinExecutor(cmd *cmdHandlerConfig, codec core.Codec, processor core.WebResponseProcessor) http.Handler {
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

	encodeBytes, _ := be.codec.EncodeBytes(requestData.Payload)
	tmpFile := be.tmpStorePayload(encodeBytes)
	defer be.removeTmpStoredPayload(tmpFile)

	var cmdOut bytes.Buffer
	var cmdErr bytes.Buffer
	cmd := exec.Command(be.cmd.Path, "-configs="+be.cmd.ConfigPath, "-action="+requestData.Action, "-payloadFile="+tmpFile.Name())
	cmd.Stdout = &cmdOut
	cmd.Stderr = &cmdErr
	if err := cmd.Run(); err != nil {
		cmdErr.WriteString(err.Error())
	}

	w.Write(be.wrp.ProcessBody(string(cmdOut.Bytes()), string(cmdErr.Bytes())))
}

func (be binExecutor) removeTmpStoredPayload(tempFile *os.File) {
	err := tempFile.Close()
	goutils.LogErrorWithMessage(fmt.Sprintf("Couldn't close tmpFile: %s", tempFile.Name()), err)

	err = os.Remove(tempFile.Name())
	goutils.LogErrorWithMessage(fmt.Sprintf("Couldn't remove tmpFile: %s", tempFile.Name()), err)
}

func (be binExecutor) tmpStorePayload(bytes []byte) *os.File {
	tempFile, err := ioutil.TempFile(os.TempDir(), "vms-")
	goutils.LogError(errors.Wrap(err, "Couldn't create tmpFile"))

	_, err = tempFile.Write(bytes)
	goutils.LogErrorWithMessage(fmt.Sprintf("Couldn't write to tmpFile: %s", tempFile.Name()), err)

	return tempFile
}
