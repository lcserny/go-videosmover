package core

import "videosmover/pkg/shutdown"

var AvailableCommands = map[string]func(){
	"shutdown": func() { shutdown.Shutdown("0") }, // TODO: use/impl seconds
}

type WithPort interface {
	GetPort() string
}

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

func (w WebviewConfig) GetPort() string {
	return w.Port
}

type CmdHandlerConfig struct {
	Type       string `json:"type"`
	Uri        string `json:"uri"`
	Path       string `json:"path"`
	ConfigPath string `json:"cfgPath"`
}

type ProxyConfig struct {
	ServerName            string             `json:"serverName"`
	LogFile               string             `json:"logFile"`
	Port                  string             `json:"port"`
	UDPPort               string             `json:"udp-port"`
	CloudDBUrl            string             `json:"cloudDbURL"`
	CloudDBAccountKeyfile string             `json:"cloudDbKeyfile"`
	Bin                   []CmdHandlerConfig `json:"bin"`
}

func (p ProxyConfig) GetPort() string {
	return p.Port
}

type CacheServerConfig struct {
	LogFile      string `json:"logFile"`
	Port         string `json:"port"`
	CacheDBPath  string `json:"cacheDBPath"`
	MaxSizeBytes int    `json:"maxSizeBytes"`
}

func (c CacheServerConfig) GetPort() string {
	return c.Port
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
	TmdbApiKey                string   `json:"tmdbApiKey"`
}
type QTorrentConfig struct {
	TorrentsUrl string `json:"torrentsUrl"`
}

type MongoDBConfig struct {
	Url        string `json:"url"`
	Database   string `json:"database"`
	Collection string `json:"collection"`
}

type RemoveTorrentConfig struct {
	LogFile  string         `json:"logFile"`
	QTorrent QTorrentConfig `json:"qtorrent"`
	MongoDB  MongoDBConfig  `json:"mongodb"`
}
