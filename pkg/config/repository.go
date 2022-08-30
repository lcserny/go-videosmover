package config

import (
	"errors"
	"io/ioutil"
	"strings"
	core "videosmover/pkg"

	"github.com/lcserny/goutils"
)

func MakeProxyConfig(configFile string, codec core.Codec) *core.ProxyConfig {
	configBytes, err := ioutil.ReadFile(configFile)
	goutils.LogFatal(err)

	// defaults
	serverConfig := &core.ProxyConfig{
		Port:    "8077",
		UDPPort: "41234",
		LogFile: "vm-proxyserver.log",
	}

	err = codec.Decode(configBytes, serverConfig)
	goutils.LogFatal(err)

	// validate
	if len(serverConfig.Port) < 4 {
		goutils.LogFatal(errors.New("port not valid"))
	}
	if len(serverConfig.ServerName) < 1 {
		goutils.LogFatal(errors.New("server name not provided"))
	}
	if len(serverConfig.CloudDBUrl) < 1 {
		goutils.LogFatal(errors.New("cloud database url not provided"))
	}
	if len(serverConfig.CloudDBAccountKeyfile) < 1 {
		goutils.LogFatal(errors.New("cloud database account keyfile not provided"))
	}
	if len(serverConfig.LogFile) < 1 {
		goutils.LogFatal(errors.New("log file path not provided"))
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

func MakeCacheServerConfig(configFile string, codec core.Codec) *core.CacheServerConfig {
	configBytes, err := ioutil.ReadFile(configFile)
	goutils.LogFatal(err)

	// defaults
	serverConfig := &core.CacheServerConfig{
		Port:         "8076",
		LogFile:      "vm-cacheserver.log",
		MaxSizeBytes: 1024 * 1024 * 10,
	}

	err = codec.Decode(configBytes, serverConfig)
	goutils.LogFatal(err)

	// validate
	if len(serverConfig.Port) < 4 {
		goutils.LogFatal(errors.New("port not valid"))
	}
	if len(serverConfig.LogFile) < 1 {
		goutils.LogFatal(errors.New("log file path not provided"))
	}
	if serverConfig.MaxSizeBytes < 128 {
		goutils.LogFatal(errors.New("max size bytes needs to be bigger than 128"))
	}
	if len(serverConfig.CacheDBPath) < 1 {
		goutils.LogFatal(errors.New("cache DB file path not provided"))
	}

	return serverConfig
}

func MakeWebviewConfig(configFile string, codec core.Codec) *core.WebviewConfig {
	configBytes, err := ioutil.ReadFile(configFile)
	goutils.LogFatal(err)

	// defaults
	config := &core.WebviewConfig{
		LogFile:               "vm-webview.log",
		AutoDarkModeEnable:    true,
		AutoDarkModeHourStart: 18,
		AutoDarkModeHourEnd:   6,
		Port:                  "8079",
		ServerPingTimeoutMs:   int64(10000),
		VideosMoverAPI:        "http://localhost:8077/exec-bin/videos-mover",
	}

	err = codec.Decode(configBytes, config)
	goutils.LogFatal(err)

	// validate
	if len(config.LogFile) < 1 {
		goutils.LogFatal(errors.New("log file path not provided"))
	}
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
	ac := core.ActionConfig{
		LogFile:       "vm-commander.log",
		CacheAddress:  "http://localhost:8076",
		CachePoolSize: 10,
	}

	err = codec.Decode(content, &ac)
	goutils.LogFatal(err)

	// validate
	if len(ac.LogFile) < 1 {
		goutils.LogFatal(errors.New("log file path not provided"))
	}
	if len(ac.CacheAddress) < 1 {
		goutils.LogFatal(errors.New("cacheAddress not set correctly"))
	}
	if ac.CachePoolSize < 1 {
		goutils.LogFatal(errors.New("cachePoolSize must be bigger than 0"))
	}
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
	if ac.HeaderBytesSize < 10 {
		goutils.LogFatal(errors.New("headerBytesSize specified too low"))
	}

	return &ac
}

func MakeRemoveTorrentConfig(configFile string, codec core.Codec) *core.RemoveTorrentConfig {
	content, err := ioutil.ReadFile(configFile)
	goutils.LogFatal(err)

	// defaults
	rtc := core.RemoveTorrentConfig{
		LogFile: "vm-remove-torrents.log",
		MongoDB: core.MongoDBConfig{
			Database:   "videosmover",
			Collection: "download_cache",
		},
	}

	err = codec.Decode(content, &rtc)
	goutils.LogFatal(err)	

	// validate
	if len(rtc.LogFile) < 1 {
		goutils.LogFatal(errors.New("log file path not provided"))
	}
	if len(rtc.QTorrent.TorrentsUrl) < 1 {
		goutils.LogFatal(errors.New("qTorrent url not provided"))
	}
	if len(rtc.MongoDB.Url) < 1 {
		goutils.LogFatal(errors.New("MongoDB connection url not provided"))
	}

	return &rtc
}
