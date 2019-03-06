package actions

import (
	"encoding/json"
	. "github.com/lcserny/goutils"
	"os"
	"strings"
)

const (
	RESTRICTED_REMOVE_PATHS_FILE_KEY = "restricted_remove_paths"
	RESTRICTED_PATH_REASON           = "Dir '%s' is a restricted path"
)

var restrictedRemovePaths []string

func init() {
	restrictedRemovePathsContent, err := configFolder.FindString(RESTRICTED_REMOVE_PATHS_FILE_KEY)
	LogError(err)
	restrictedRemovePaths = GetLinesFromString(restrictedRemovePathsContent)
}

func pathRemovalIsRestricted(folder string) bool {
	for _, restrictedFolder := range restrictedRemovePaths {
		if strings.HasSuffix(folder, restrictedFolder) {
			return true
		}
	}
	return false
}

func getJSONEncodedString(data interface{}) string {
	resultBytes, err := json.Marshal(data)
	LogError(err)
	return string(resultBytes)
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
