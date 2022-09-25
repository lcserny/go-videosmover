package main

import (
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"time"
	core "videosmover/pkg"
	"videosmover/pkg/config"
	"videosmover/pkg/ext/json"
	"videosmover/pkg/move"
	"videosmover/pkg/output"
	"videosmover/pkg/search"
	"videosmover/pkg/web"

	"github.com/lcserny/goutils"
)

func main() {
	// validate startup
	args := os.Args[1:]
	if len(args) != 1 {
		_, _ = fmt.Fprintln(os.Stderr, "ERROR: Please provide `config` flag")
		os.Exit(1)
	}

	cfgPath := flag.String("config", "", "path to webview config")
	flag.Parse()

	// var pingTimestamp int64
	jsonCodec := json.NewJsonCodec()
	cfg := config.MakeWebviewConfig(*cfgPath, jsonCodec)
	goutils.InitFileLogger(cfg.LogFile)
	apiRequester := web.NewApiRequester(jsonCodec)
	webPath := fmt.Sprintf(":%s", cfg.Port)

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
	htmlServer := http.FileServer(http.Dir(cfg.HtmlFilesPath))
	mux.Handle("/static/", http.StripPrefix("/static/", htmlServer))
	templates := template.Must(template.ParseGlob(filepath.Join(cfg.HtmlFilesPath, "*.gohtml")))
	for pat, tmplController := range tmplControllers {
		mux.HandleFunc(pat, func(resp http.ResponseWriter, req *http.Request) {
			if tmplName, tmplData, renderTmpl := tmplController.ServeTemplate(resp, req); renderTmpl {
				if cfg.AutoDarkModeEnable {
					tmplData.DarkMode = time.Now().Hour() > cfg.AutoDarkModeHourStart || time.Now().Hour() < cfg.AutoDarkModeHourEnd
				}
				goutils.LogFatal(templates.ExecuteTemplate(resp, fmt.Sprintf("%s.gohtml", tmplName), tmplData))
			}
		})
	}
	for pat, controller := range ajaxHandlers {
		mux.Handle(pat, controller)
	}

	// start server
	goutils.LogInfo(fmt.Sprintf("Started server on port %s...", cfg.Port))
	if err := http.ListenAndServe(webPath, mux); err != http.ErrServerClosed {
		goutils.LogFatal(err)
	}
}
