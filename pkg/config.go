package core

type WebviewConfig struct {
	AutoDarkModeEnable  bool   `json:"autoDarkModeEnable"`
	AutoDarkModeHourMax int    `json:"autoDarkModeHourMax"`
	AutoDarkModeHourMin int    `json:"autoDarkModeHourMin"`
	Port                string `json:"port"`
	HtmlFilesPath       string `json:"htmlFilesPath"`
	ServerPingTimeoutMs int64  `json:"serverPingTimeoutMs"`
	VideosMoverAPI      string `json:"videosMoverAPI"`
	DownloadsPath       string `json:"downloadsPath"`
	MoviesPath          string `json:"moviesPath"`
	TvSeriesPath        string `json:"tvSeriesPath"`
}

type CmdHandlerConfig struct {
	Uri        string `json:"uri"`
	Path       string `json:"path"`
	ConfigPath string `json:"cfgPath"`
}

type ProxyConfig struct {
	Port string             `json:"port"`
	Bin  []CmdHandlerConfig `json:"bin"`
}

type CacheConfig struct {
	Port                  string `json:"port"`
	CachePath             string `json:"cachePath"`
	MaxSizeBytes          int    `json:"maxSizeBytes"`
	PersistenceIntervalMs int64  `json:"persistenceIntervalMs"`
}

type ActionConfig struct {
	MinimumVideoSize          int64    `json:"minimumVideoSize"`
	SimilarityPercent         int      `json:"similarityPercent"`
	MaxOutputWalkDepth        int      `json:"maxOutputWalkDepth"`
	MaxSearchWalkDepth        int      `json:"maxSearchWalkDepth"`
	MaxWebSearchResultCount   int      `json:"maxWebSearchResultCount"`
	HeaderBytesSize           int      `json:"headerBytesSize"`
	RestrictedRemovePaths     []string `json:"restrictedRemovePaths"`
	NameTrimRegexes           []string `json:"nameTrimRegexes"`
	SearchExcludePaths        []string `json:"searchExcludePaths"`
	AllowedMIMETypes          []string `json:"allowedMIMETypes"`
	AllowedSubtitleExtensions []string `json:"allowedSubtitleExtensions"`
	CacheApiURL               string   `json:"cacheApiUrl"`
}
