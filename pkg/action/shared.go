package action

import (
	"errors"
	"fmt"
	"github.com/lcserny/go-videosmover/pkg/convert"
	"github.com/lcserny/go-videosmover/pkg/delete"
	"github.com/lcserny/go-videosmover/pkg/move"
	"github.com/lcserny/go-videosmover/pkg/output"
	"github.com/lcserny/go-videosmover/pkg/search"
	"os"
	"regexp"
	"strings"
)

const (
	MOVIE   = "movie"
	TV      = "tv"
	UNKNOWN = "unknown"

	RESTRICTED_PATH_REASON = "Dir '%s' is a restricted path"
)

var ActionsMap = map[string]Action{
	"SEARCH": search.Action,
	"OUTPUT": output.Action,
	"MOVE":   move.Action,
	"DELETE": delete.Action,
}

type Action func(jsonPayload []byte, config *convert.ActionConfig) (string, error)

func UnknownAction(jsonPayload []byte, config *convert.ActionConfig) (string, error) {
	return "", errors.New("Unknown action given")
}

func PathRemovalIsRestricted(folder string, restrictedRemovePaths []string) bool {
	for _, restrictedFolder := range restrictedRemovePaths {
		if strings.HasSuffix(folder, restrictedFolder) {
			return true
		}
	}
	return false
}

func WalkDepthIsAcceptable(rootPath string, path string, maxWalkDepth int) bool {
	trimmed := path[len(rootPath):]
	separatorCount := 0
	for _, char := range trimmed {
		if char == os.PathSeparator {
			separatorCount++
		}
	}
	return separatorCount < maxWalkDepth
}

func GetDiskPath(videoType string, config *convert.WebviewConfig) string {
	loweredType := strings.ToLower(videoType)
	if loweredType != UNKNOWN {
		diskPath := config.MoviesPath
		if loweredType != MOVIE {
			diskPath = config.TvSeriesPath
		}
		return diskPath
	}
	return ""
}

func GetRegexList(patterns []string) (regxs []*regexp.Regexp) {
	for _, pat := range patterns {
		regxs = append(regxs, regexp.MustCompile(fmt.Sprintf("(?i)(-?%s)", pat)))
	}
	return regxs
}
