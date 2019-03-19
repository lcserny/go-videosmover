package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/lcserny/go-videosmover/pkg/handlers"
	. "github.com/lcserny/goutils"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

type ProxyServerConfig struct {
	Host                       string `json:"host"`
	Port                       string `json:"port"`
	PathVideosMoverJava        string `json:"path.videosMover.java"`
	PathVideosMoverJavaConfigs string `json:"path.videosMover.java.configs"`
	PathVideosMoverBin         string `json:"path.videosMover.bin"`
	PathVideosMoverBinConfigs  string `json:"path.videosMover.bin.configs"`
}

var proxyServerConfigsPathFlag = flag.String("configPath", "", "path to proxy server configs folder")

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		_, _ = fmt.Fprintln(os.Stderr, "ERROR: Please provide `configPath` flag")
		return
	}

	InitCurrentPathFileLogger("vm-proxyserver.log")

	flag.Parse()
	serverConfig := generateServerConfig(*proxyServerConfigsPathFlag, fmt.Sprintf("config_%s.json", runtime.GOOS))

	mux := http.NewServeMux()
	mux.Handle("/exec-java/videos-mover", handlers.NewJavaJsonExecuteHandler(serverConfig))
	mux.Handle("/exec-bin/videos-mover", handlers.NewBinJsonExecuteHandler(serverConfig))

	LogInfo(fmt.Sprintf("Started server on %s port %s...", serverConfig.Host, serverConfig.Port))
	LogFatal(http.ListenAndServe(fmt.Sprintf("%s:%s", serverConfig.Host, serverConfig.Port), mux))
}

func generateServerConfig(configsPath, configFile string) *ProxyServerConfig {
	configBytes, err := ioutil.ReadFile(filepath.Join(configsPath, configFile))
	LogFatal(err)

	var serverConfig ProxyServerConfig
	err = json.Unmarshal(configBytes, &serverConfig)
	LogFatal(err)

	if serverConfig.Host == "" || serverConfig.Port == "" {
		LogFatal(errors.New("No `host` and/or `port` configured"))
	}

	return &serverConfig
}
