package actions

import (
	"encoding/json"
	"github.com/h2non/filetype"
	. "github.com/lcserny/goutils"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	EXCLUDE_LIST_FILE  = "exclude_paths"
	MIME_TYPES_FILE    = "mime_types"
	MIN_VIDEO_SIZE_KEY = "minimum.video.size"
	MAX_WALK_DEPTH     = 4
	HEADER_BYTES_SIZE  = 261
)

var (
	excludePaths []string
	mimeTypes    []string
	minFileSize  int64
)

func init() {
	excludePathsContent, err := configFolder.FindString(EXCLUDE_LIST_FILE)
	LogError(err)
	excludePaths = strings.Split(excludePathsContent, "\n")

	mimeTypesContent, err := configFolder.FindString(MIME_TYPES_FILE)
	LogError(err)
	mimeTypes = strings.Split(mimeTypesContent, "\n")

	if appProperties.HasProperty(MIN_VIDEO_SIZE_KEY) {
		minFileSize = appProperties.GetPropertyAsInt64(MIN_VIDEO_SIZE_KEY)
	}
}

type SearchAction struct {
}

// TODO: put these in a shared go project `go-videosmover-shared`
type RequestSearchData struct {
	Path string `json:"path"`
}

type ResponseSearchData struct {
	Path string `json:"path"`
}

func (a *SearchAction) Execute(jsonFile string) (string, error) {
	jsonRequestBytes, err := ioutil.ReadFile(jsonFile)
	LogError(err)
	if err != nil {
		return "", err
	}

	var request RequestSearchData
	err = json.Unmarshal(jsonRequestBytes, &request)
	LogError(err)
	if err != nil {
		return "", err
	}

	realWalkRootPath, _ := GetRealPath(request.Path)
	var resultList []ResponseSearchData
	err = filepath.Walk(realWalkRootPath, func(path string, info os.FileInfo, err error) error {
		LogError(err)
		if info != nil && !info.IsDir() && walkDepthIsAcceptable(realWalkRootPath, path, MAX_WALK_DEPTH) {
			if isVideo(path, info) {
				resultList = append(resultList, ResponseSearchData{path})
			}
		}
		return nil
	})
	LogError(err)
	if err != nil {
		return "", err
	}

	if len(resultList) < 1 {
		return "", nil
	}

	resultBytes, err := json.Marshal(resultList)
	LogError(err)
	if err != nil {
		return "", err
	}

	return string(resultBytes), nil
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

func isVideo(path string, info os.FileInfo) bool {
	file, err := os.Open(path)
	if err != nil {
		LogError(err)
		return false
	}
	defer CloseFile(file)
	head := make([]byte, HEADER_BYTES_SIZE)
	n, _ := file.Read(head)
	if n < HEADER_BYTES_SIZE {
		return false
	}

	// check path
	for _, exPath := range excludePaths {
		if strings.Contains(path, exPath) {
			return false
		}
	}

	// check type
	acceptedMime := false
	for _, mType := range mimeTypes {
		if filetype.IsMIME(head, mType) {
			acceptedMime = true
			break
		}
	}
	if !acceptedMime && !filetype.IsVideo(head) {
		return false
	}

	// check size
	if info.Size() < minFileSize {
		return false
	}

	return true
}
