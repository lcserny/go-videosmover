package web

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
)

type ProxyServerConfig struct {
	Host                       string `json:"host"`
	Port                       string `json:"port"`
	PathVideosMoverJava        string `json:"path.videosMover.java"`
	PathVideosMoverJavaConfigs string `json:"path.videosMover.java.configs"`
	PathVideosMoverBin         string `json:"path.videosMover.bin"`
	PathVideosMoverBinConfigs  string `json:"path.videosMover.bin.configs"`
}

func GenerateServerConfig(configsPath, configFile string) *models.ProxyServerConfig {
	configBytes, err := ioutil.ReadFile(filepath.Join(configsPath, configFile))
	LogFatal(err)

	var serverConfig models.ProxyServerConfig
	err = json.Unmarshal(configBytes, &serverConfig)
	LogFatal(err)

	if serverConfig.Host == "" || serverConfig.Port == "" {
		LogFatal(errors.New("No `host` and/or `port` configured"))
	}

	return &serverConfig
}
