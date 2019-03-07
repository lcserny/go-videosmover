package main

import (
	"fmt"
	"github.com/gobuffalo/packr"
	"github.com/lcserny/go-videosmover/pkg/handlers"
	. "github.com/lcserny/goutils"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
)

const (
	SERVER_HOST_KEY = "server.host"
	SERVER_PORT_KEY = "server.port"
	LOG_PATH_KEY    = "log.path"
)

var serverProperties *ConfigProperties

func init() {
	content, err := packr.NewBox("../cfg").FindString(fmt.Sprintf("server_%s.properties", runtime.GOOS))
	LogFatal(err)

	serverProperties = ReadProperties(content)
	if serverProperties.HasProperty(LOG_PATH_KEY) {
		initServerLogger(serverProperties.GetPropertyAsString(LOG_PATH_KEY))
	}
}

func initServerLogger(logPath string) {
	openFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	LogFatal(err)
	writer := io.MultiWriter(os.Stdout, openFile)
	log.SetOutput(writer)
}

func main() {
	if !serverProperties.HasProperty(SERVER_PORT_KEY) || !serverProperties.HasProperty(SERVER_HOST_KEY) {
		LogFatalWithMessage(fmt.Sprintf("Please provide %s and %s in config", SERVER_HOST_KEY, SERVER_PORT_KEY), nil)
	}

	port := serverProperties.GetPropertyAsInt(SERVER_PORT_KEY)
	host := serverProperties.GetPropertyAsString(SERVER_HOST_KEY)

	mux := http.NewServeMux()
	mux.Handle("/exec-java/videos-mover", handlers.NewJavaJsonExecuteHandler(serverProperties))
	mux.Handle("/exec-bin/videos-mover", handlers.NewBinJsonExecuteHandler(serverProperties))

	LogInfo(fmt.Sprintf("Started server on %s port %d...", host, port))
	LogFatal(http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), mux))
}
