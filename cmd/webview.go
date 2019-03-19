package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	. "github.com/lcserny/goutils"
	"github.com/pkg/browser"
	"github.com/pkg/errors"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type WebviewConfig struct {
	Host                string `json:"host"`
	Port                string `json:"port"`
	HtmlFilesPattern    string `json:"htmlFilesPattern"`
	ServerPingTimeoutMs int64  `json:"serverPingTimeoutMs"`
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

	InitCurrentPathFileLogger("vm-webview.log")

	flag.Parse()
	config := generateWebviewConfig(*wvConfigsPath, fmt.Sprintf("config_%s.json", runtime.GOOS))

	webPath := fmt.Sprintf("%s:%s", config.Host, config.Port)
	handler := generateHandler(config.HtmlFilesPattern)
	server := startFileServer(webPath, handler)
	go openBrowser(webPath)
	checkStopServer(server, config)
}

func generateWebviewConfig(configsPath, configFile string) *WebviewConfig {
	configBytes, err := ioutil.ReadFile(filepath.Join(configsPath, configFile))
	LogFatal(err)

	var config WebviewConfig
	err = json.Unmarshal(configBytes, &config)
	LogFatal(err)

	if config.Host == "" || config.Port == "" {
		LogFatal(errors.New("No `host` and/or `port` configured"))
	}

	return &config
}

func checkStopServer(server *http.Server, config *WebviewConfig) {
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

func generateHandler(htmlFilesPattern string) *http.ServeMux {
	templates := template.Must(template.ParseGlob(htmlFilesPattern))

	mux := http.NewServeMux()
	mux.HandleFunc("/", defaultHtmlTemplateHandle(templates))
	mux.HandleFunc("/running", handleRunningPing)
	// TODO: add more paths for ajax calls maybe?
	return mux
}

func handleRunningPing(writer http.ResponseWriter, request *http.Request) {
	lastRunningPingTimestamp = MakeTimestamp()
}

func defaultHtmlTemplateHandle(tmpl *template.Template) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		templateName := request.URL.Path
		if templateName == "/favicon.ico" {
			return
		}

		if templateName == "/" {
			templateName = "index.html"
		}

		if strings.HasPrefix(templateName, "/") {
			templateName = templateName[1:]
		}

		LogFatal(tmpl.ExecuteTemplate(writer, templateName, nil))
	}
}

func openBrowser(webPath string) {
	if !strings.HasPrefix(webPath, "http") {
		webPath = fmt.Sprintf("http://%s", webPath)
	}
	LogFatal(browser.OpenURL(webPath))
}
