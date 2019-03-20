package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/lcserny/go-videosmover/pkg/handlers"
	"github.com/lcserny/go-videosmover/pkg/view"
	. "github.com/lcserny/goutils"
	"github.com/pkg/browser"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type TemplatedController func(writer http.ResponseWriter, request *http.Request) (tmplName string, tmplData interface{}, renderTmpl bool)

var (
	wvConfigsPath            = flag.String("configPath", "", "path to webview config files")
	lastRunningPingTimestamp = MakeTimestamp()

	templatedViewsMap = map[string]TemplatedController{
		"/":              view.SearchController,
		"/search":        view.SearchController,
		"/searchResults": view.SearchResultsController,
	}
)

// FIXME: why is it using 30% CPU? profile it
func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		_, _ = fmt.Fprintln(os.Stderr, "ERROR: Please provide `configPath` flag")
		return
	}

	InitCurrentPathFileLogger("vm-webview.log")

	flag.Parse()
	config := handlers.GenerateWebviewConfig(*wvConfigsPath, fmt.Sprintf("config_%s.json", runtime.GOOS))

	webPath := fmt.Sprintf("%s:%s", config.Host, config.Port)
	handler := generateHandler(config.HtmlFilesPath)
	server := startFileServer(webPath, handler)
	go openBrowser(webPath)
	checkStopServer(server, config)
}

func checkStopServer(server *http.Server, config *handlers.WebviewConfig) {
	for {
		if MakeTimestamp() > lastRunningPingTimestamp+config.ServerPingTimeoutMs {
			LogInfo(fmt.Sprintf("No ping received in %d ms, stopping server", config.ServerPingTimeoutMs))
			LogFatal(server.Shutdown(context.TODO()))
		}
	}
}

func startFileServer(webPath string, handler *http.ServeMux) *http.Server {
	server := &http.Server{Addr: webPath, Handler: handler}
	go func() {
		LogInfo(fmt.Sprintf("Started server on %s...", webPath))
		LogFatal(server.ListenAndServe())
	}()
	return server
}

// TODO: add more paths if needed
func generateHandler(htmlFilesPath string) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/running", func(writer http.ResponseWriter, request *http.Request) {
		lastRunningPingTimestamp = MakeTimestamp()
	})

	staticServer := http.FileServer(http.Dir(filepath.Join(htmlFilesPath, "static")))
	mux.Handle("/static/", http.StripPrefix("/static/", staticServer))

	templates := template.Must(template.ParseGlob(filepath.Join(htmlFilesPath, "*.gohtml")))
	for pat, tmplView := range templatedViewsMap {
		mux.HandleFunc(pat, func(writer http.ResponseWriter, request *http.Request) {
			tmplName, tmplData, renderTmpl := tmplView(writer, request)
			if renderTmpl {
				LogFatal(templates.ExecuteTemplate(writer, tmplName, tmplData))
			}
		})
	}

	return mux
}

func openBrowser(webPath string) {
	if !strings.HasPrefix(webPath, "http") {
		webPath = fmt.Sprintf("http://%s", webPath)
	}
	LogFatal(browser.OpenURL(webPath))
}
