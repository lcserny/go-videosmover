package web

import (
	"encoding/json"
	"errors"
	"github.com/lcserny/goutils"
	"io/ioutil"
	"path/filepath"
)

type ProxyConfig struct {
	Host                       string `json:"host"`
	Port                       string `json:"port"`
	PathVideosMoverJava        string `json:"path.videosMover.java"`
	PathVideosMoverJavaConfigs string `json:"path.videosMover.java.configs"`
	PathVideosMoverBin         string `json:"path.videosMover.bin"`
	PathVideosMoverBinConfigs  string `json:"path.videosMover.bin.configs"`
}

func GenerateProxyConfig(path, file string) *ProxyConfig {
	configBytes, err := ioutil.ReadFile(filepath.Join(path, file))
	goutils.LogFatal(err)

	var serverConfig ProxyConfig
	err = json.Unmarshal(configBytes, &serverConfig)
	goutils.LogFatal(err)

	if serverConfig.Host == "" || serverConfig.Port == "" {
		goutils.LogFatal(errors.New("No `host` and/or `port` configured"))
	}

	return &serverConfig
}
