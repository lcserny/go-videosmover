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
	handler := generateHandler(config)
	server := startFileServer(webPath, handler)
	go openBrowser(webPath)
	checkStopServer(server, config)
}

func checkStopServer(server *http.Server, config *models.WebviewConfig) {
	for {
		if MakeTimestamp() > lastRunningPingTimestamp+config.ServerPingTimeoutMs {
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

func generateHandler(config *models.WebviewConfig) *http.ServeMux {
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

	return mux
}

func templateControllers(config *models.WebviewConfig) map[string]handlers.TemplateController {
	templatesMap := make(map[string]handlers.TemplateController)
	searchController := view.NewSearchController(config)
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
