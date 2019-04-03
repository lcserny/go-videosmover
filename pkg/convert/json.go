package convert

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

const TIME_FORMAT = "2006-01-02 15:04:05"

func getJSONEncodedString(data interface{}) string {
	resultBytes, err := json.Marshal(data)
	LogError(err)
	return string(resultBytes)
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
