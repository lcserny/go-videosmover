package shutdown

import (
	"bytes"
	"os/exec"
	"runtime"

	"github.com/lcserny/goutils"
)

func Shutdown(seconds string) {
	var cmdErr bytes.Buffer
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/C", "shutdown", "-s", "-t", seconds)
	case "linux":
		cmd = exec.Command("shutdown", seconds)
	default:
		goutils.LogWarning("Unknown OS detected, cannot shutdown.")
	}

	cmd.Stderr = &cmdErr
	if err := cmd.Run(); err != nil {
		goutils.LogError(err)
		cmdErr.WriteString(err.Error())
	}
}
