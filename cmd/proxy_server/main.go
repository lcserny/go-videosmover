package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/lcserny/goutils"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"videosmover/pkg/config"
	"videosmover/pkg/ext/json"
	"videosmover/pkg/web"
)

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		_, _ = fmt.Fprintln(os.Stderr, "ERROR: Please provide `config` flag")
		os.Exit(1)
	}

	cfgPath := flag.String("config", "", "path to proxy server config")
	flag.Parse()

	jsonCodec := json.NewJsonCodec()
	apiRequester := web.NewApiRequester(jsonCodec)
	c := config.MakeProxyConfig(*cfgPath, jsonCodec)
	goutils.InitFileLogger(c.LogFile)

	mux := http.NewServeMux()
	for _, binCmd := range c.Bin {
		mux.Handle(fmt.Sprintf("/exec-bin/%s", binCmd.Uri), web.NewBinExecutor(&binCmd, jsonCodec, apiRequester))
	}
	addInternalHandlers(mux)

	goutils.LogInfo(fmt.Sprintf("Started server on port %s...", c.Port))
	goutils.LogFatal(http.ListenAndServe(fmt.Sprintf(":%s", c.Port), mux))
}

func addInternalHandlers(mux *http.ServeMux) {
	addShutdownEndpoint(mux)
}

func addShutdownEndpoint(mux *http.ServeMux) {
	if runtime.GOOS != "windows" {
		goutils.LogInfo(fmt.Sprintf("Shutdown available for windows only, OS found: %s", runtime.GOOS))
		return
	}

	mux.HandleFunc("/shutdown", func(writer http.ResponseWriter, request *http.Request) {
		values := request.URL.Query()
		secondsInt := "0"
		if seconds, exists := values["seconds"]; exists {
			secondsInt = seconds[0]
		}

		var cmdErr bytes.Buffer
		cmd := exec.Command("cmd", "/C", "shutdown", "-s", "-t", secondsInt)
		cmd.Stderr = &cmdErr
		if err := cmd.Run(); err != nil {
			cmdErr.WriteString(err.Error())
		}
	})
}
