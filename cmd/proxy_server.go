package main

import (
	"fmt"
	"github.com/lcserny/go-videosmover/pkg/handlers"
	. "github.com/lcserny/goutils"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
)

func main() {
	initServerLogger()
	serverConfig := handlers.GenerateServerConfig("../../cfg/server", fmt.Sprintf("config_%s.json", runtime.GOOS))

	mux := http.NewServeMux()
	mux.Handle("/exec-java/videos-mover", handlers.NewJavaJsonExecuteHandler(serverConfig))
	mux.Handle("/exec-bin/videos-mover", handlers.NewBinJsonExecuteHandler(serverConfig))

	LogInfo(fmt.Sprintf("Started server on %s port %s...", serverConfig.Host, serverConfig.Port))
	LogFatal(http.ListenAndServe(fmt.Sprintf("%s:%s", serverConfig.Host, serverConfig.Port), mux))
}

func initServerLogger() {
	openFile, err := os.OpenFile(GetAbsCurrentPathOf("vm-proxyserver.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	LogFatal(err)
	writer := io.MultiWriter(os.Stdout, openFile)
	log.SetOutput(writer)
}
