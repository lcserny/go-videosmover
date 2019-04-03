package action

import (
	"github.com/lcserny/go-videosmover/pkg/models"
	"os"
	"strings"
)

const (
	RESTRICTED_PATH_REASON = "Dir '%s' is a restricted path"

	MOVIE   = "movie"
	TV      = "tv"
	UNKNOWN = "unknown"
)

func pathRemovalIsRestricted(folder string, restrictedRemovePaths []string) bool {
	for _, restrictedFolder := range restrictedRemovePaths {
		if strings.HasSuffix(folder, restrictedFolder) {
			return true
		}
	}
	return false
}

func getDiskPath(videoType string, config *models.WebviewConfig) string {
	loweredType := strings.ToLower(videoType)
	if loweredType != models.UNKNOWN {
		diskPath := config.MoviesPath
		if loweredType != models.MOVIE {
			diskPath = config.TvSeriesPath
		}
		return diskPath
	}
	return ""
}

func walkDepthIsAcceptable(rootPath string, path string, maxWalkDepth int) bool {
	trimmed := path[len(rootPath):]
	separatorCount := 0
	for _, char := range trimmed {
		if char == os.PathSeparator {
			separatorCount++
		}
	}
	return separatorCount < maxWalkDepth
}
