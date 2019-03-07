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
	MAX_OUTPUT_WALK_DEPTH_KEY   = "max.output.walk.depth"
	MAX_SEARCH_WALK_DEPTH_KEY   = "max.search.walk.depth"
	MAX_TMDB_RESULT_COUNT_KEY   = "max.tmdb.result.count"
	HEADER_BYTES_SIZE_KEY       = "header.bytes.size"
	RESTRICTED_REMOVE_PATHS_KEY = "restricted.remove.paths.file"
	NAME_TRIM_REGX_KEY          = "name.trim.regx.file"
	SEARCH_EXCLUDE_PATHS_KEY    = "search.exclude.paths.file"
	ALLOWED_MIME_TYPES_KEY      = "allowed.mime.types.file"
	ALLOWED_SUBTITLE_EXTS_KEY   = "allowed.subtitle.exts.file"
	SIM_PERCENT_KEY             = "similarity.percent"
	MIN_VIDEO_SIZE_KEY          = "minimum.video.size"
	TMDB_API_KEY                = "TMDB_API_KEY"
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

func GenerateActionConfig(propertiesFile string) *ActionConfig {
	configFolder := packr.NewBox("../../cfg")
	content, err := configFolder.FindString(propertiesFile)
	LogFatal(err)

	appProperties := ReadProperties(content)

	config := ActionConfig{
		maxOutputWalkDepth: appProperties.GetPropertyAsInt(MAX_OUTPUT_WALK_DEPTH_KEY),
		maxSearchWalkDepth: appProperties.GetPropertyAsInt(MAX_SEARCH_WALK_DEPTH_KEY),
		similarityPercent:  appProperties.GetPropertyAsInt(SIM_PERCENT_KEY),
		minVideoFileSize:   appProperties.GetPropertyAsInt64(MIN_VIDEO_SIZE_KEY),
		maxTMDBResultCount: appProperties.GetPropertyAsInt(MAX_TMDB_RESULT_COUNT_KEY),
		headerBytesSize:    appProperties.GetPropertyAsInt(HEADER_BYTES_SIZE_KEY),
	}

	if key, exists := os.LookupEnv(TMDB_API_KEY); exists {
		config.tmdbAPI = tmdb.Init(tmdb.Config{key, false, nil})
	}

	restrictedRemovePathsContent, err := configFolder.FindString(appProperties.GetPropertyAsString(RESTRICTED_REMOVE_PATHS_KEY))
	LogError(err)
	config.restrictedRemovePaths = GetLinesFromString(restrictedRemovePathsContent)

	nameTrimPartsContent, err := configFolder.FindString(appProperties.GetPropertyAsString(NAME_TRIM_REGX_KEY))
	LogError(err)
	config.nameTrimPartRegexs = getRegexList(GetLinesFromString(nameTrimPartsContent))

	excludePathsContent, err := configFolder.FindString(appProperties.GetPropertyAsString(SEARCH_EXCLUDE_PATHS_KEY))
	LogError(err)
	config.searchExcludePaths = GetLinesFromString(excludePathsContent)

	mimeTypesContent, err := configFolder.FindString(appProperties.GetPropertyAsString(ALLOWED_MIME_TYPES_KEY))
	LogError(err)
	config.allowedMimeTypes = GetLinesFromString(mimeTypesContent)

	allowedSubtitleExtsContent, err := configFolder.FindString(appProperties.GetPropertyAsString(ALLOWED_SUBTITLE_EXTS_KEY))
	LogError(err)
	config.allowedSubtitleExts = GetLinesFromString(allowedSubtitleExtsContent)

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
