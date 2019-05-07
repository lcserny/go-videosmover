package web

import (
	"github.com/lcserny/goutils"
	"io/ioutil"
	"path/filepath"
	"videosmover/pkg/json"
)

type WebviewConfig struct {
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

	config := &WebviewConfig{
		Port:                "8079",
		ServerPingTimeoutMs: int64(10000),
		VideosMoverAPI:      "http://localhost:8077/exec-bin/videos-mover",
	}

	err = json.Decode(configBytes, config)
	goutils.LogFatal(err)

	return config
}
