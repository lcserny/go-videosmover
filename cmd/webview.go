package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/lcserny/go-videosmover/pkg/handlers"
	"github.com/lcserny/goutils"
	"github.com/pkg/browser"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"
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

	mux := generateHandler(config.HtmlFilesPath)
	_ = startFileServer(webPath, mux)
	// TODO: from js send pings to /health which sets the timestamp of ping,
	//  in another place, loop until ping time is older than time.Now() - 30sec?
	//  this js that pings needs to be placed on each html page, so include it somehow? server side includes?
	// stopFileServer(server)
}

func startFileServer(webPath string, handler *http.ServeMux) *http.Server {
	server := &http.Server{Addr: webPath, Handler: handler}
	go func() {
		goutils.LogInfo(fmt.Sprintf("Started server on %s...", webPath))
		goutils.LogFatal(server.ListenAndServe())
	}()
	return server
}

func stopFileServer(server *http.Server) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	goutils.LogFatal(server.Shutdown(ctx))
}

func generateHandler(htmlDir string) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(htmlDir)))
	// TODO: add more paths for ajax calls maybe?
	return mux
}

func openBrowser(webPath string) {
	if !strings.HasPrefix(webPath, "http") {
		webPath = fmt.Sprintf("http://%s", webPath)
	}
	goutils.LogFatal(browser.OpenURL(webPath))
}
