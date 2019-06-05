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
	"time"
	"videosmover/pkg"
	"videosmover/pkg/config"
	"videosmover/pkg/ext/json"
	"videosmover/pkg/move"
	"videosmover/pkg/output"
	"videosmover/pkg/search"
	"videosmover/pkg/web"
)

func main() {
	// validate startup
	args := os.Args[1:]
	if len(args) != 1 {
		_, _ = fmt.Fprintln(os.Stderr, "ERROR: Please provide `config` flag")
		return
	}

	cfgPath := flag.String("config", "", "path to webview config")
	flag.Parse()

	var pingTimestamp int64
	jsonCodec := json.NewJsonCodec()
	cfg := config.MakeWebviewConfig(*cfgPath, jsonCodec)
	goutils.InitFileLogger(cfg.LogFile)
	apiRequester := web.NewApiRequester(jsonCodec)
	webPath := fmt.Sprintf("localhost:%s", cfg.Port)

	// define template controllers
	tmplControllers := make(map[string]core.WebTemplateController)
	searchController := search.NewController(cfg, jsonCodec, apiRequester)
	tmplControllers["/"] = searchController
	tmplControllers["/search"] = searchController

	// define AJAX handlers
	ajaxHandlers := make(map[string]http.Handler)
	ajaxHandlers["/ajax/output"] = output.NewAjaxController(cfg, jsonCodec, apiRequester)
	ajaxHandlers["/ajax/move"] = move.NewAjaxController(cfg, jsonCodec, apiRequester)

	// init web handler
	mux := http.NewServeMux()
	mux.HandleFunc("/running", func(writer http.ResponseWriter, request *http.Request) {
		pingTimestamp = goutils.MakeTimestamp()
	})
	htmlServer := http.FileServer(http.Dir(cfg.HtmlFilesPath))
	mux.Handle("/static/", http.StripPrefix("/static/", htmlServer))
	templates := template.Must(template.ParseGlob(filepath.Join(cfg.HtmlFilesPath, "*.gohtml")))
	for pat, tmplController := range tmplControllers {
		mux.HandleFunc(pat, func(resp http.ResponseWriter, req *http.Request) {
			if tmplName, tmplData, renderTmpl := tmplController.ServeTemplate(resp, req); renderTmpl {
				if cfg.AutoDarkModeEnable {
					tmplData.DarkMode = time.Now().Hour() > cfg.AutoDarkModeHourMax || time.Now().Hour() < cfg.AutoDarkModeHourMin
				}
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
		if (pingTimestamp != 0) && (goutils.MakeTimestamp() > pingTimestamp+cfg.ServerPingTimeoutMs) {
			goutils.LogFatal(server.Shutdown(context.Background()))
		}
	}
}
