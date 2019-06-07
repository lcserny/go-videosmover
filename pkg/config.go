package core

type WebviewConfig struct {
	LogFile               string `json:"logFile"`
	AutoDarkModeEnable    bool   `json:"autoDarkModeEnable"`
	AutoDarkModeHourStart int    `json:"autoDarkModeHourStart"`
	AutoDarkModeHourEnd   int    `json:"autoDarkModeHourEnd"`
	Port                  string `json:"port"`
	HtmlFilesPath         string `json:"htmlFilesPath"`
	ServerPingTimeoutMs   int64  `json:"serverPingTimeoutMs"`
	VideosMoverAPI        string `json:"videosMoverAPI"`
	DownloadsPath         string `json:"downloadsPath"`
	MoviesPath            string `json:"moviesPath"`
	TvSeriesPath          string `json:"tvSeriesPath"`
}

type CmdHandlerConfig struct {
	Uri        string `json:"uri"`
	Path       string `json:"path"`
	ConfigPath string `json:"cfgPath"`
}

type ProxyConfig struct {
	LogFile string             `json:"logFile"`
	Port    string             `json:"port"`
	Bin     []CmdHandlerConfig `json:"bin"`
}

type ActionConfig struct {
	LogFile                   string   `json:"logFile"`
	MinimumVideoSize          int64    `json:"minimumVideoSize"`
	SimilarityPercent         int      `json:"similarityPercent"`
	MaxOutputWalkDepth        int      `json:"maxOutputWalkDepth"`
	MaxSearchWalkDepth        int      `json:"maxSearchWalkDepth"`
	MaxWebSearchResultCount   int      `json:"maxWebSearchResultCount"`
	CacheAddress              string   `json:"cacheAddress"`
	CachePoolSize             int      `json:"cachePoolSize"`
	HeaderBytesSize           int      `json:"headerBytesSize"`
	RestrictedRemovePaths     []string `json:"restrictedRemovePaths"`
	NameTrimRegexes           []string `json:"nameTrimRegexes"`
	SearchExcludePaths        []string `json:"searchExcludePaths"`
	AllowedMIMETypes          []string `json:"allowedMIMETypes"`
	AllowedSubtitleExtensions []string `json:"allowedSubtitleExtensions"`
}
