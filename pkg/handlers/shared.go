package handlers

import (
	"encoding/json"
	"fmt"
	. "github.com/lcserny/goutils"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const TIME_FORMAT = "2006-01-02 15:04:05"

type RequestJsonData struct {
	Action  string      `json:"action"`
	Payload interface{} `json:"payload,string"`
}

type ResponseJsonData struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
	Date  string `json:"date"`
	Body  string `json:"body"`
}

type ProxyServerConfig struct {
	Host                       string `json:"host"`
	Port                       string `json:"port"`
	PathVideosMoverJava        string `json:"path.videosMover.java"`
	PathVideosMoverJavaConfigs string `json:"path.videosMover.java.configs"`
	PathVideosMoverBin         string `json:"path.videosMover.bin"`
	PathVideosMoverBinConfigs  string `json:"path.videosMover.bin.configs"`
}

func GenerateServerConfig(configsPath, configFile string) *ProxyServerConfig {
	configBytes, err := ioutil.ReadFile(filepath.Join(configsPath, configFile))
	LogFatal(err)

	var serverConfig ProxyServerConfig
	err = json.Unmarshal(configBytes, &serverConfig)
	LogFatal(err)

	if serverConfig.Host == "" || serverConfig.Port == "" {
		LogFatal(errors.New("No `host` and/or `port` configured"))
	}

	return &serverConfig
}

func removeTmpStoredJsonPayload(tempJsonFile *os.File) {
	err := tempJsonFile.Close()
	LogErrorWithMessage(fmt.Sprintf("Couldn't close tmpFile: %s", tempJsonFile.Name()), err)

	err = os.Remove(tempJsonFile.Name())
	LogErrorWithMessage(fmt.Sprintf("Couldn't remove tmpFile: %s", tempJsonFile.Name()), err)
}

func tmpStoreJsonPayload(jsonData interface{}) *os.File {
	tempFile, err := ioutil.TempFile(os.TempDir(), "vms-")
	LogErrorWithMessage(fmt.Sprintf("Couldn't create tmpFile: %s", tempFile.Name()), err)

	jsonString, err := json.Marshal(jsonData)
	LogErrorWithMessage("Couldn't convert data to bytes", err)
	if err == nil {
		_, err = tempFile.Write([]byte(jsonString))
	}
	LogErrorWithMessage(fmt.Sprintf("Couldn't write to tmpFile: %s", tempFile.Name()), err)

	return tempFile
}

func getJsonResponseFromAsBytes(body, err string) []byte {
	if strings.Contains(body, "ERROR") {
		err = body
		body = ""
	}

	code := 200
	if len(err) > 0 {
		code = 500
	}

	responseJsonData := &ResponseJsonData{
		Code:  code,
		Error: err,
		Date:  time.Now().Format(TIME_FORMAT),
		Body:  body,
	}

	jsonBytes, _ := json.Marshal(responseJsonData)
	return jsonBytes
}

func getErrorJsonResponseAsBytes(err string) []byte {
	responseJsonData := &ResponseJsonData{
		Code:  500,
		Error: err,
		Date:  time.Now().Format(TIME_FORMAT),
		Body:  "",
	}

	jsonBytes, _ := json.Marshal(responseJsonData)
	return jsonBytes
}
