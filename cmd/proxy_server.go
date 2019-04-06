package main

import (
	"flag"
	"fmt"
	"github.com/lcserny/go-videosmover/pkg/web"
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
	c := web.GenerateProxyConfig(*proxyServerConfigsPathFlag, fmt.Sprintf("config_%s.json", runtime.GOOS))

	mux := http.NewServeMux()
	for _, binCmd := range c.Bin {
		mux.Handle(fmt.Sprintf("/exec-bin/%s", binCmd.Uri), &web.BinJsonExecuteHandler{&binCmd})
	}

	LogInfo(fmt.Sprintf("Started server on port %s...", c.Port))
	LogFatal(http.ListenAndServe(fmt.Sprintf(":%s", c.Port), mux))
}
