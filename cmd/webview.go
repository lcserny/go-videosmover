package main

import (
	"flag"
	"fmt"
	"github.com/lcserny/go-videosmover/pkg/handlers"
	"github.com/lcserny/goutils"
	"github.com/pkg/browser"
	"net/http"
	"os"
	"runtime"
	"strings"
)

var wvConfigsPath = flag.String("configPath", "", "path to webview config files")

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		_, _ = fmt.Fprintln(os.Stderr, "ERROR: Please provide `configPath` flag")
		return
	}

	goutils.InitCurrentPathFileLogger("vm-webview.log")

	flag.Parse()
	config := handlers.GenerateWebviewConfig(*wvConfigsPath, fmt.Sprintf("config_%s.json", runtime.GOOS))

	webPath := fmt.Sprintf("%s:%s", config.Host, config.Port)
	go openBrowser(webPath)
	startFileServer(webPath, config.HtmlFilesPath)
}

func startFileServer(webPath, htmlDir string) {
	http.Handle("/", http.FileServer(http.Dir(htmlDir)))
	goutils.LogInfo(fmt.Sprintf("Started server on %s...", webPath))
	goutils.LogFatal(http.ListenAndServe(webPath, nil))
}

func openBrowser(webPath string) {
	if !strings.HasPrefix(webPath, "http") {
		webPath = fmt.Sprintf("http://%s", webPath)
	}
	goutils.LogFatal(browser.OpenURL(webPath))
}
