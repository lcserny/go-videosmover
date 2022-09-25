package shutdown

import (
	"bytes"
	"os/exec"
	"runtime"
)

func Shutdown(seconds string) {
	var cmdErr bytes.Buffer
	os := runtime.GOOS
	var cmd *exec.Cmd
	switch os {
	case "windows":
		cmd = exec.Command("cmd", "/C", "shutdown", "-s", "-t", seconds)
	case "linux":
		cmd = exec.Command("poweroff")
	}
	cmd.Stderr = &cmdErr
	if err := cmd.Run(); err != nil {
		cmdErr.WriteString(err.Error())
	}
}
