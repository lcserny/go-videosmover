package shutdown

import (
	"bytes"
	"os/exec"
)

func Shutdown(seconds string) {
	var cmdErr bytes.Buffer
	cmd := exec.Command("cmd", "/C", "shutdown", "-s", "-t", seconds)
	cmd.Stderr = &cmdErr
	if err := cmd.Run(); err != nil {
		cmdErr.WriteString(err.Error())
	}
}
