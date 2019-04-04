package web

import (
	"encoding/json"
	"errors"
	"github.com/lcserny/goutils"
	"io/ioutil"
	"path/filepath"
)

type WebviewConfig struct {
	Host                string `json:"host"`
	Port                string `json:"port"`
	HtmlFilesPath       string `json:"htmlFilesPath"`
	ServerPingTimeoutMs int64  `json:"serverPingTimeoutMs"`
	VideosMoverAPI      string `json:"videosMoverAPI"`
	DownloadsPath       string `json:"downloadsPath"`
	MoviesPath          string `json:"moviesPath"`
	TvSeriesPath        string `json:"tvSeriesPath"`
}

func GenerateWebviewConfig(configsPath, configFile string) *WebviewConfig {
	configBytes, err := ioutil.ReadFile(filepath.Join(configsPath, configFile))
	goutils.LogFatal(err)

	var config WebviewConfig
	err = json.Unmarshal(configBytes, &config)
	goutils.LogFatal(err)

	if config.Host == "" || config.Port == "" {
		goutils.LogFatal(errors.New("No `host` and/or `port` configured"))
	}

	return &config
}
