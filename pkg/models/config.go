package models

import (
	"github.com/ryanbradynd05/go-tmdb"
	"regexp"
)

type ActionConfig struct {
	MinimumVideoSize          int64    `json:"minimumVideoSize"`
	SimilarityPercent         int      `json:"similarityPercent"`
	MaxOutputWalkDepth        int      `json:"maxOutputWalkDepth"`
	MaxSearchWalkDepth        int      `json:"maxSearchWalkDepth"`
	MaxTMDBResultCount        int      `json:"maxTMDBResultCount"`
	OutTMDBCacheLimit         int      `json:"outTMDBCacheLimit"`
	HeaderBytesSize           int      `json:"headerBytesSize"`
	RestrictedRemovePaths     []string `json:"restrictedRemovePaths"`
	NameTrimRegexes           []string `json:"nameTrimRegexes"`
	SearchExcludePaths        []string `json:"searchExcludePaths"`
	AllowedMIMETypes          []string `json:"allowedMIMETypes"`
	AllowedSubtitleExtensions []string `json:"allowedSubtitleExtensions"`

	TmdbAPI                 *tmdb.TMDb
	CompiledNameTrimRegexes []*regexp.Regexp
}

type ProxyServerConfig struct {
	Host                       string `json:"host"`
	Port                       string `json:"port"`
	PathVideosMoverJava        string `json:"path.videosMover.java"`
	PathVideosMoverJavaConfigs string `json:"path.videosMover.java.configs"`
	PathVideosMoverBin         string `json:"path.videosMover.bin"`
	PathVideosMoverBinConfigs  string `json:"path.videosMover.bin.configs"`
}

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
