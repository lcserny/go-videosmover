package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/lcserny/go-videosmover/pkg/convert"
	"github.com/lcserny/go-videosmover/pkg/generate"
	inhttp "github.com/lcserny/go-videosmover/pkg/http"
	"github.com/lcserny/go-videosmover/pkg/move"
	"github.com/lcserny/go-videosmover/pkg/output"
	"github.com/lcserny/go-videosmover/pkg/search"
	. "github.com/lcserny/goutils"
	"github.com/pkg/browser"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var (
	wvConfigsPath            = flag.String("configPath", "", "path to webview config files")
	lastRunningPingTimestamp int64
)

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		_, _ = fmt.Fprintln(os.Stderr, "ERROR: Please provide `configPath` flag")
		return
	}

	InitFileLogger("vm-webview.log")

	flag.Parse()
	config := generate.NewWebviewConfig(*wvConfigsPath, fmt.Sprintf("config_%s.json", runtime.GOOS))

	webPath := fmt.Sprintf("%s:%s", config.Host, config.Port)
	handler := generateHandler(config)
	server := startFileServer(webPath, handler)
	go openBrowser(webPath)
	checkStopServer(server, config)
}

func checkStopServer(server *http.Server, config *convert.WebviewConfig) {
	for {
		if (lastRunningPingTimestamp != 0) && (MakeTimestamp() > lastRunningPingTimestamp+config.ServerPingTimeoutMs) {
			LogFatal(server.Shutdown(context.TODO()))
		}
		time.Sleep(time.Second)
	}
}

func startFileServer(webPath string, handler *http.ServeMux) *http.Server {
	server := &http.Server{Addr: webPath, Handler: handler}
	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			LogFatal(err)
		}
		os.Exit(0)
	}()
	return server
}

func generateHandler(config *convert.WebviewConfig) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/running", func(writer http.ResponseWriter, request *http.Request) {
		lastRunningPingTimestamp = MakeTimestamp()
	})

	staticServer := http.FileServer(http.Dir(filepath.Join(config.HtmlFilesPath, "static")))
	mux.Handle("/static/", http.StripPrefix("/static/", staticServer))

	templates := template.Must(template.ParseGlob(filepath.Join(config.HtmlFilesPath, "*.gohtml")))
	for pat, tmplController := range templateControllers(config) {
		mux.HandleFunc(pat, func(resp http.ResponseWriter, req *http.Request) {
			if tmplName, tmplData, renderTmpl := tmplController.ServeTemplate(resp, req); renderTmpl {
				LogFatal(templates.ExecuteTemplate(resp, fmt.Sprintf("%s.gohtml", tmplName), tmplData))
			}
		})
	}

	for pat, controller := range ajaxHandlers(config) {
		mux.Handle(pat, controller)
	}

	return mux
}

func templateControllers(config *convert.WebviewConfig) map[string]inhttp.TemplateController {
	templatesMap := make(map[string]inhttp.TemplateController)
	searchController := search.NewController(config)
	templatesMap["/"] = searchController
	templatesMap["/search"] = searchController
	// TODO: add more if needed

	return templatesMap
}

func ajaxHandlers(config *convert.WebviewConfig) map[string]http.Handler {
	templatesMap := make(map[string]http.Handler)
	templatesMap["/ajax/output"] = output.NewAjaxController(config)
	templatesMap["/ajax/move"] = move.NewAjaxController(config)
	// TODO: add more if needed

	return templatesMap
}

func openBrowser(webPath string) {
	if !strings.HasPrefix(webPath, "http") {
		webPath = fmt.Sprintf("http://%s", webPath)
	}
	LogFatal(browser.OpenURL(webPath))
}
