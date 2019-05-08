package action

import (
	"os"
	"strings"
	"videosmover/pkg"
)

const (
	RESTRICTED_PATH_REASON = "Dir '%s' is a restricted path"

	MOVIE   = "movie"
	TV      = "tv"
	UNKNOWN = "unknown"
)

func PathRemovalIsRestricted(folder string, restrictedRemovePaths []string) bool {
	for _, restrictedFolder := range restrictedRemovePaths {
		if strings.HasSuffix(folder, restrictedFolder) {
			return true
		}
	}
	return false
}

func GetDiskPath(videoType string, config *core.WebviewConfig) string {
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
