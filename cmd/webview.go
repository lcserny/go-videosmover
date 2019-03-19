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
	// TODO: figure out how to stop server when tab closed? JS loop?
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
	goutils.LogFatal(server.Shutdown(context.TODO()))
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
