package convert

import (
	"encoding/json"
	"fmt"
	"github.com/lcserny/goutils"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"testing"
)

func GetJSONEncodedString(data interface{}) string {
	resultBytes, err := json.Marshal(data)
	goutils.LogError(err)
	return string(resultBytes)
}

func GetJSONBytesForTest(t *testing.T, data interface{}) []byte {
	bytes, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Couldn't decode request: %+v", err)
	}
	return bytes
}

func GetJSONStringForTest(t *testing.T, data interface{}) string {
	bytes, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Couldn't decode response: %+v", err)
	}
	return string(bytes)
}

func RemoveTmpStoredJsonPayload(tempJsonFile *os.File) {
	err := tempJsonFile.Close()
	goutils.LogErrorWithMessage(fmt.Sprintf("Couldn't close tmpFile: %s", tempJsonFile.Name()), err)

	err = os.Remove(tempJsonFile.Name())
	goutils.LogErrorWithMessage(fmt.Sprintf("Couldn't remove tmpFile: %s", tempJsonFile.Name()), err)
}

func TmpStoreJsonPayload(jsonData interface{}) *os.File {
	tempFile, err := ioutil.TempFile(os.TempDir(), "vms-")
	goutils.LogError(errors.Wrap(err, "Couldn't create tmpFile"))

	jsonString, err := json.Marshal(jsonData)
	goutils.LogError(errors.Wrap(err, "Couldn't convert data to bytes"))
	if err == nil {
		_, err = tempFile.Write([]byte(jsonString))
	}
	goutils.LogErrorWithMessage(fmt.Sprintf("Couldn't write to tmpFile: %s", tempFile.Name()), err)

	return tempFile
}
