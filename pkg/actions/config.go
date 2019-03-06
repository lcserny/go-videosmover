package actions

import (
	"github.com/gobuffalo/packr"
	. "github.com/lcserny/goutils"
	"github.com/pkg/errors"
	"github.com/ryanbradynd05/go-tmdb"
	"os"
	"regexp"
	"strings"
)

const (
	RESTRICTED_REMOVE_PATHS_FILE = "restricted_remove_paths"
	NAME_TRIM_REGX_FILE          = "name_trim_regx"
	SEARCH_EXCLUDE_PATHS_FILE    = "search_exclude_paths"
	ALLOWED_MIME_TYPES_FILE      = "allowed_mime_types"
	ALLOWED_SUBTITLE_EXTS_FILE   = "allowed_subtitle_exts"

	SIM_PERCENT_KEY    = "similarity.percent"
	MIN_VIDEO_SIZE_KEY = "minimum.video.size"
	TMDB_API_KEY       = "TMDB_API_KEY"
)

// TODO: add action to use from qBittorrent when done downloading to add to a db or something,
//  then in Android app on startup it can maybe show you last finished downloading torrents
var actionsMap = map[string]Action{
	"SEARCH": SearchAction,
	"OUTPUT": OutputAction,
	"MOVE":   MoveAction,
	"DELETE": DeleteAction,
}

type ActionConfig struct {
	restrictedRemovePaths []string
	searchExcludePaths    []string
	allowedMimeTypes      []string
	allowedSubtitleExts   []string
	maxOutputWalkDepth    int
	maxSearchWalkDepth    int
	maxTMDBResultCount    int
	similarityPercent     int
	headerBytesSize       int
	minVideoFileSize      int64
	tmdbAPI               *tmdb.TMDb
	nameTrimPartRegexs    []*regexp.Regexp
}

type Action func(jsonPayload []byte, config *ActionConfig) (string, error)

// TODO: don't use constants for util files like RESTRICTED_REMOVE_PATHS_FILE, get the file name from properties
func GenerateActionConfig(propertiesFile string) *ActionConfig {
	configFolder := packr.NewBox("../../config")
	content, err := configFolder.FindString(propertiesFile)
	LogFatal(err)

	appProperties := ReadProperties(content)

	restrictedRemovePathsContent, err := configFolder.FindString(RESTRICTED_REMOVE_PATHS_FILE)
	LogError(err)
	restrictedRemovePaths := GetLinesFromString(restrictedRemovePathsContent)

	nameTrimPartsContent, err := configFolder.FindString(NAME_TRIM_REGX_FILE)
	LogError(err)
	nameTrimPartsRegxs := getRegexList(GetLinesFromString(nameTrimPartsContent))

	excludePathsContent, err := configFolder.FindString(SEARCH_EXCLUDE_PATHS_FILE)
	LogError(err)
	excludePaths := GetLinesFromString(excludePathsContent)

	mimeTypesContent, err := configFolder.FindString(ALLOWED_MIME_TYPES_FILE)
	LogError(err)
	mimeTypes := GetLinesFromString(mimeTypesContent)

	allowedSubtitleExtsContent, err := configFolder.FindString(ALLOWED_SUBTITLE_EXTS_FILE)
	LogError(err)
	allowedSubtitleExts := GetLinesFromString(allowedSubtitleExtsContent)

	config := ActionConfig{
		restrictedRemovePaths: restrictedRemovePaths,
		searchExcludePaths:    excludePaths,
		allowedMimeTypes:      mimeTypes,
		allowedSubtitleExts:   allowedSubtitleExts,
		maxOutputWalkDepth:    2,
		maxSearchWalkDepth:    4,
		maxTMDBResultCount:    10,
		headerBytesSize:       261,
		nameTrimPartRegexs:    nameTrimPartsRegxs,
	}

	if appProperties.HasProperty(SIM_PERCENT_KEY) {
		config.similarityPercent = appProperties.GetPropertyAsInt(SIM_PERCENT_KEY)
	}

	if key, exists := os.LookupEnv(TMDB_API_KEY); exists {
		config.tmdbAPI = tmdb.Init(tmdb.Config{key, false, nil})
	}

	if appProperties.HasProperty(MIN_VIDEO_SIZE_KEY) {
		config.minVideoFileSize = appProperties.GetPropertyAsInt64(MIN_VIDEO_SIZE_KEY)
	}

	return &config
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
