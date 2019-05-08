package config

import (
	"github.com/lcserny/goutils"
	"io/ioutil"
	"path/filepath"
	"videosmover/pkg"
)

func GenerateProxyConfig(path, file string, codec core.Codec) *core.ProxyConfig {
	configBytes, err := ioutil.ReadFile(filepath.Join(path, file))
	goutils.LogFatal(err)

	// defaults
	serverConfig := &core.ProxyConfig{Port: "8077"}

	err = codec.Decode(configBytes, serverConfig)
	goutils.LogFatal(err)

	return serverConfig
}

func GenerateWebviewConfig(configsPath, configFile string, codec core.Codec) *core.WebviewConfig {
	configBytes, err := ioutil.ReadFile(filepath.Join(configsPath, configFile))
	goutils.LogFatal(err)

	// defaults
	config := &core.WebviewConfig{
		Port:                "8079",
		ServerPingTimeoutMs: int64(10000),
		VideosMoverAPI:      "http://localhost:8077/exec-bin/videos-mover",
	}

	err = codec.Decode(configBytes, config)
	goutils.LogFatal(err)

	return config
}

func NewActionConfig(cfgPath string, codec core.Codec) *core.ActionConfig {
	content, err := ioutil.ReadFile(cfgPath)
	goutils.LogFatal(err)

	var ac core.ActionConfig
	err = codec.Decode(content, &ac)
	goutils.LogFatal(err)

	return &ac
}
