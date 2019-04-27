package main

import (
	"flag"
	"fmt"
	"github.com/lcserny/go-videosmover/pkg/web"
	"github.com/lcserny/goutils"
	"net/http"
	"os"
	"runtime"
)

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		_, _ = fmt.Fprintln(os.Stderr, "ERROR: Please provide `configPath` flag")
		return
	}

	goutils.InitFileLogger("vm-proxyserver.log")

	cfgPath := flag.String("configPath", "", "path to proxy server configs folder")
	flag.Parse()

	c := web.GenerateProxyConfig(*cfgPath, fmt.Sprintf("config_%s.json", runtime.GOOS))

	mux := http.NewServeMux()
	for _, binCmd := range c.Bin {
		mux.Handle(fmt.Sprintf("/exec-bin/%s", binCmd.Uri), &web.BinJsonExecuteHandler{&binCmd})
	}

	goutils.LogInfo(fmt.Sprintf("Started server on port %s...", c.Port))
	goutils.LogFatal(http.ListenAndServe(fmt.Sprintf(":%s", c.Port), mux))
}
