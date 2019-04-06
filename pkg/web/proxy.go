package web

import (
	"encoding/json"
	"errors"
	"github.com/lcserny/goutils"
	"io/ioutil"
	"path/filepath"
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

	var serverConfig ProxyConfig
	err = json.Unmarshal(configBytes, &serverConfig)
	goutils.LogFatal(err)

	if serverConfig.Port == "" {
		goutils.LogFatal(errors.New("no port configured"))
	}

	return &serverConfig
}
