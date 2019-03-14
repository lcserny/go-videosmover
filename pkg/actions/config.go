package actions

import (
	"encoding/json"
	"github.com/gobuffalo/packr"
	. "github.com/lcserny/goutils"
	"github.com/pkg/errors"
	"github.com/ryanbradynd05/go-tmdb"
	"os"
	"regexp"
	"strings"
)

var actionsMap = map[string]Action{
	"SEARCH":         SearchAction,
	"OUTPUT":         OutputAction,
	"MOVE":           MoveAction,
	"DELETE":         DeleteAction,
	"REMOVE_TORRENT": RemoveTorrentAction,
}

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

	tmdbAPI                 *tmdb.TMDb
	compiledNameTrimRegexes []*regexp.Regexp
}

type Action func(jsonPayload []byte, config *ActionConfig) (string, error)

func GenerateActionConfig(path, file string) *ActionConfig {
	configFolder := packr.NewBox(path)
	content, err := configFolder.Find(file)
	LogFatal(err)

	var actionConfig ActionConfig
	err = json.Unmarshal(content, &actionConfig)
	LogFatal(err)

	if key, exists := os.LookupEnv("TMDB_API_KEY"); exists {
		actionConfig.tmdbAPI = tmdb.Init(tmdb.Config{key, false, nil})
	}

	if actionConfig.NameTrimRegexes != nil {
		actionConfig.compiledNameTrimRegexes = getRegexList(actionConfig.NameTrimRegexes)
	}

	return &actionConfig
}

func UnknownAction(jsonPayload []byte, config *ActionConfig) (string, error) {
	return "", errors.New("Unknown action given")
}

func NewActionFrom(val string) Action {
	if action, ok := actionsMap[strings.ToUpper(val)]; ok {
		return action
	}
	return UnknownAction
}
