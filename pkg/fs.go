package core

import (
	"fmt"
	"github.com/lcserny/goutils"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
)

func RemoveTmpStoredPayload(tempFile *os.File) {
	err := tempFile.Close()
	goutils.LogErrorWithMessage(fmt.Sprintf("Couldn't close tmpFile: %s", tempFile.Name()), err)

	err = os.Remove(tempFile.Name())
	goutils.LogErrorWithMessage(fmt.Sprintf("Couldn't remove tmpFile: %s", tempFile.Name()), err)
}

func TmpStorePayload(bytes []byte) *os.File {
	tempFile, err := ioutil.TempFile(os.TempDir(), "vms-")
	goutils.LogError(errors.Wrap(err, "Couldn't create tmpFile"))

	_, err = tempFile.Write(bytes)
	goutils.LogErrorWithMessage(fmt.Sprintf("Couldn't write to tmpFile: %s", tempFile.Name()), err)

	return tempFile
}
