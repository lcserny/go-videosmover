package web

import (
	"github.com/lcserny/goutils"
	"io/ioutil"
	"path/filepath"
	"videosmover/pkg/json"
)

type cmdHandlerConfig struct {
	Uri        string `json:"uri"`
	Path       string `json:"path"`
	ConfigPath string `json:"cfgPath"`
}

type ProxyConfig struct {
	Port string             `json:"port"`
	Bin  []cmdHandlerConfig `json:"bin"`
}

func GenerateProxyConfig(path, file string) *ProxyConfig {
	configBytes, err := ioutil.ReadFile(filepath.Join(path, file))
	goutils.LogFatal(err)

	serverConfig := &ProxyConfig{Port: "8077"}
	err = json.Decode(configBytes, serverConfig)
	goutils.LogFatal(err)

	return serverConfig
}
