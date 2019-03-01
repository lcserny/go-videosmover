package actions

import (
	"encoding/json"
	. "github.com/lcserny/goutils"
	"os"
)

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
