package config

import (
	"errors"
	"github.com/lcserny/goutils"
	"io/ioutil"
	"strings"
	"videosmover/pkg"
)

func ProxyProxyConfig(configFile string, codec core.Codec) *core.ProxyConfig {
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

func MakeWebviewConfig(configFile string, codec core.Codec) *core.WebviewConfig {
	configBytes, err := ioutil.ReadFile(configFile)
	goutils.LogFatal(err)

	// defaults
	config := &core.WebviewConfig{
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

	var ac core.ActionConfig
	err = codec.Decode(content, &ac)
	goutils.LogFatal(err)

	// validate
	if ac.MinimumVideoSize < 1000000 {
		goutils.LogFatal(errors.New("minimumVideoSize less than 1mb"))
	}
	if ac.SimilarityPercent < 1 {
		goutils.LogFatal(errors.New("similarityPercent is not a positive number"))
	}
	if ac.MaxOutputWalkDepth < 1 || ac.MaxSearchWalkDepth < 1 {
		goutils.LogFatal(errors.New("walkDepths are not a positive number"))
	}
	if ac.MaxWebSearchResultCount < 1 {
		goutils.LogFatal(errors.New("webSearchCount is not a positive number"))
	}
	if ac.OutWebSearchCacheLimit < 10 {
		goutils.LogFatal(errors.New("webSearchCacheLimit specified too low"))
	}
	if ac.HeaderBytesSize < 10 {
		goutils.LogFatal(errors.New("headerBytesSize specified too low"))
	}

	return &ac
}
