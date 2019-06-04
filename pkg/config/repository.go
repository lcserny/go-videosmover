package config

import (
	"errors"
	"github.com/lcserny/goutils"
	"io/ioutil"
	"strings"
	"videosmover/pkg"
)

func MakeProxyConfig(configFile string, codec core.Codec) *core.ProxyConfig {
	configBytes, err := ioutil.ReadFile(configFile)
	goutils.LogFatal(err)

	// defaults
	serverConfig := &core.ProxyConfig{Port: "8077"}

	err = codec.Decode(configBytes, serverConfig)
	goutils.LogFatal(err)

	// validate
	if len(serverConfig.Port) < 4 {
		goutils.LogFatal(errors.New("port not valid"))
	}
	if len(serverConfig.Bin) > 0 {
		for _, e := range serverConfig.Bin {
			if strings.HasPrefix(e.Uri, "/") || strings.HasSuffix(e.Uri, "/") {
				goutils.LogFatal(errors.New("serverConfig URI should not start or end with slash"))
			}
			if len(e.Path) < 3 || len(e.ConfigPath) < 3 {
				goutils.LogFatal(errors.New("serverConfig path and cfgPath should not be empty"))
			}
		}
	}

	return serverConfig
}

func MakeCacheConfig(configFile string, codec core.Codec) *core.CacheConfig {
	configBytes, err := ioutil.ReadFile(configFile)
	goutils.LogFatal(err)

	// defaults
	serverConfig := &core.CacheConfig{
		Port:                  "8075",
		MaxSizeBytes:          10000000,
		PersistenceIntervalMs: 1800000,
	}

	err = codec.Decode(configBytes, serverConfig)
	goutils.LogFatal(err)

	// validate
	if len(serverConfig.Port) < 4 {
		goutils.LogFatal(errors.New("port not valid"))
	}
	if serverConfig.MaxSizeBytes < 512 {
		goutils.LogFatal(errors.New("maxSizeBytes too small"))
	}
	if serverConfig.PersistenceIntervalMs < 5999 {
		goutils.LogFatal(errors.New("persistenceIntervalMs too small"))
	}

	return serverConfig
}

func MakeWebviewConfig(configFile string, codec core.Codec) *core.WebviewConfig {
	configBytes, err := ioutil.ReadFile(configFile)
	goutils.LogFatal(err)

	// defaults
	config := &core.WebviewConfig{
		AutoDarkModeEnable:  true,
		AutoDarkModeHourMax: 18,
		AutoDarkModeHourMin: 6,
		Port:                "8079",
		ServerPingTimeoutMs: int64(10000),
		VideosMoverAPI:      "http://localhost:8077/exec-bin/videos-mover",
	}

	err = codec.Decode(configBytes, config)
	goutils.LogFatal(err)

	// validate
	if len(config.Port) < 4 {
		goutils.LogFatal(errors.New("port not valid"))
	}
	if config.ServerPingTimeoutMs < 100 {
		goutils.LogFatal(errors.New("server ping timeout is too low"))
	}
	if !strings.HasPrefix(config.VideosMoverAPI, "http") {
		goutils.LogFatal(errors.New("invalid videosMover API url provided"))
	}
	if len(config.HtmlFilesPath) < 3 || len(config.DownloadsPath) < 3 || len(config.MoviesPath) < 3 || len(config.TvSeriesPath) < 3 {
		goutils.LogFatal(errors.New("invalid paths specified for one of [html, downloads, movies, tvShows]"))
	}

	return config
}

func MakeActionConfig(cfgPath string, codec core.Codec) *core.ActionConfig {
	content, err := ioutil.ReadFile(cfgPath)
	goutils.LogFatal(err)

	// defaults
	config := &core.ActionConfig{
		CacheApiURL: "http://localhost:8075/cache",
	}
	err = codec.Decode(content, config)
	goutils.LogFatal(err)

	// validate
	if config.MinimumVideoSize < 1000000 {
		goutils.LogFatal(errors.New("minimumVideoSize less than 1mb"))
	}
	if config.SimilarityPercent < 1 {
		goutils.LogFatal(errors.New("similarityPercent is not a positive number"))
	}
	if config.MaxOutputWalkDepth < 1 || config.MaxSearchWalkDepth < 1 {
		goutils.LogFatal(errors.New("walkDepths are not a positive number"))
	}
	if config.MaxWebSearchResultCount < 1 {
		goutils.LogFatal(errors.New("webSearchCount is not a positive number"))
	}
	if config.HeaderBytesSize < 10 {
		goutils.LogFatal(errors.New("headerBytesSize specified too low"))
	}

	return config
}
