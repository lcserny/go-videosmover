package shutdown

import (
	"bytes"
	"os/exec"

	"github.com/lcserny/goutils"
)

func Shutdown(seconds string) {
	var cmdErr bytes.Buffer
	cmd := exec.Command("cmd", "/C", "shutdown", "-s", "-t", seconds)
	cmd.Stderr = &cmdErr
	if err := cmd.Run(); err != nil {
		goutils.LogError(err)
		cmdErr.WriteString(err.Error())
	}
}
