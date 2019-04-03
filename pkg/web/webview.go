package web

import (
	"encoding/json"
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

func GenerateWebviewConfig(configsPath, configFile string) *models.WebviewConfig {
	configBytes, err := ioutil.ReadFile(filepath.Join(configsPath, configFile))
	LogFatal(err)

	var config models.WebviewConfig
	err = json.Unmarshal(configBytes, &config)
	LogFatal(err)

	if config.Host == "" || config.Port == "" {
		LogFatal(errors.New("No `host` and/or `port` configured"))
	}

	return &config
}
