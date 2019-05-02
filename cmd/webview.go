package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/lcserny/goutils"
	"github.com/pkg/browser"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"
	"videosmover/pkg/move"
	"videosmover/pkg/output"
	"videosmover/pkg/search"
	"videosmover/pkg/web"
)

func main() {
	// validate startup
	args := os.Args[1:]
	if len(args) != 1 {
		_, _ = fmt.Fprintln(os.Stderr, "ERROR: Please provide `configPath` flag")
		return
	}

	goutils.InitFileLogger("vm-webview.log")

	cfgPath := flag.String("configPath", "", "path to webview config files")
	flag.Parse()

	var pingTimestamp int64
	config := web.GenerateWebviewConfig(*cfgPath, fmt.Sprintf("config_%s.json", runtime.GOOS))
	webPath := fmt.Sprintf("localhost:%s", config.Port)

	// define template controllers
	tmplControllers := make(map[string]web.TemplateController)
	searchController := search.NewController(config)
	tmplControllers["/"] = searchController
	tmplControllers["/search"] = searchController

	// define AJAX handlers
	ajaxHandlers := make(map[string]http.Handler)
	ajaxHandlers["/ajax/output"] = output.NewAjaxController(config)
	ajaxHandlers["/ajax/move"] = move.NewAjaxController(config)

	// init web handler
	mux := http.NewServeMux()
	mux.HandleFunc("/running", func(writer http.ResponseWriter, request *http.Request) {
		pingTimestamp = goutils.MakeTimestamp()
	})
	staticServer := http.FileServer(http.Dir(filepath.Join(config.HtmlFilesPath, "static")))
	mux.Handle("/static/", http.StripPrefix("/static/", staticServer))
	templates := template.Must(template.ParseGlob(filepath.Join(config.HtmlFilesPath, "*.gohtml")))
	for pat, tmplController := range tmplControllers {
		mux.HandleFunc(pat, func(resp http.ResponseWriter, req *http.Request) {
			if tmplName, tmplData, renderTmpl := tmplController.ServeTemplate(resp, req); renderTmpl {
				goutils.LogFatal(templates.ExecuteTemplate(resp, fmt.Sprintf("%s.gohtml", tmplName), tmplData))
			}
		})
	}
	for pat, controller := range ajaxHandlers {
		mux.Handle(pat, controller)
	}

	// start server
	server := &http.Server{Addr: webPath, Handler: mux}
	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			goutils.LogFatal(err)
		}
		os.Exit(0)
	}()

	// open browser
	go goutils.LogFatal(browser.OpenURL(fmt.Sprintf("http://%s", webPath)))

	// check shutdown server
	for range time.NewTicker(time.Second).C {
		if (pingTimestamp != 0) && (goutils.MakeTimestamp() > pingTimestamp+config.ServerPingTimeoutMs) {
			goutils.LogFatal(server.Shutdown(context.Background()))
		}
	}
}
