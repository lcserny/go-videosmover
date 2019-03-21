package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/lcserny/go-videosmover/pkg/handlers"
	"github.com/lcserny/go-videosmover/pkg/models"
	"github.com/lcserny/go-videosmover/pkg/view"
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

type TemplatedHandler interface {
	ServeTemplate(resp http.ResponseWriter, req *http.Request) (name string, data interface{}, render bool)
}

var (
	wvConfigsPath            = flag.String("configPath", "", "path to webview config files")
	lastRunningPingTimestamp = MakeTimestamp()
)

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		_, _ = fmt.Fprintln(os.Stderr, "ERROR: Please provide `configPath` flag")
		return
	}

	InitFileLogger("vm-webview.log")

	flag.Parse()
	config := handlers.GenerateWebviewConfig(*wvConfigsPath, fmt.Sprintf("config_%s.json", runtime.GOOS))

	webPath := fmt.Sprintf("%s:%s", config.Host, config.Port)
	handler := generateHandler(config.HtmlFilesPath)
	server := startFileServer(webPath, handler)
	go openBrowser(webPath)
	checkStopServer(server, config)
}

func checkStopServer(server *http.Server, config *models.WebviewConfig) {
	for {
		if MakeTimestamp() > lastRunningPingTimestamp+config.ServerPingTimeoutMs {
			LogInfo(fmt.Sprintf("No ping received in %d ms, stopping server", config.ServerPingTimeoutMs))
			LogFatal(server.Shutdown(context.TODO()))
		}
		time.Sleep(time.Second)
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

func generateHandler(htmlFilesPath string) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/running", func(writer http.ResponseWriter, request *http.Request) {
		lastRunningPingTimestamp = MakeTimestamp()
	})

	staticServer := http.FileServer(http.Dir(filepath.Join(htmlFilesPath, "static")))
	mux.Handle("/static/", http.StripPrefix("/static/", staticServer))

	for pat, tmplHandler := range getTemplatedHandlers() {
		mux.HandleFunc(pat, func(resp http.ResponseWriter, req *http.Request) {
			if tmplName, tmplData, renderTmpl := tmplHandler.ServeTemplate(resp, req); renderTmpl {
				tmpl := template.Must(template.ParseFiles(
					filepath.Join(htmlFilesPath, "layout.gohtml"),
					filepath.Join(htmlFilesPath, fmt.Sprintf("%s.gohtml", tmplName))),
				)
				LogFatal(tmpl.Execute(resp, tmplData))
			}
		})
	}

	return mux
}

func getTemplatedHandlers() map[string]TemplatedHandler {
	templatesMap := make(map[string]TemplatedHandler)
	searchController := &view.SearchController{}
	templatesMap["/"] = searchController
	templatesMap["/search"] = searchController
	// TODO: add more if needed
	return templatesMap
}

func openBrowser(webPath string) {
	if !strings.HasPrefix(webPath, "http") {
		webPath = fmt.Sprintf("http://%s", webPath)
	}
	LogFatal(browser.OpenURL(webPath))
}
