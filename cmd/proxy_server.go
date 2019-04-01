package main

import (
	"flag"
	"fmt"
	"github.com/lcserny/go-videosmover/pkg/generate"
	inhttp "github.com/lcserny/go-videosmover/pkg/http"
	. "github.com/lcserny/goutils"
	"net/http"
	"os"
	"runtime"
)

var proxyServerConfigsPathFlag = flag.String("configPath", "", "path to proxy server configs folder")

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		_, _ = fmt.Fprintln(os.Stderr, "ERROR: Please provide `configPath` flag")
		return
	}

	InitFileLogger("vm-proxyserver.log")

	flag.Parse()
	serverConfig := generate.NewServerConfig(*proxyServerConfigsPathFlag, fmt.Sprintf("config_%s.json", runtime.GOOS))

	mux := http.NewServeMux()
	mux.Handle("/exec-java/videos-mover", inhttp.NewJavaJsonExecuteHandler(serverConfig))
	mux.Handle("/exec-bin/videos-mover", inhttp.NewBinJsonExecuteHandler(serverConfig))

	LogInfo(fmt.Sprintf("Started server on %s port %s...", serverConfig.Host, serverConfig.Port))
	LogFatal(http.ListenAndServe(fmt.Sprintf("%s:%s", serverConfig.Host, serverConfig.Port), mux))
}
