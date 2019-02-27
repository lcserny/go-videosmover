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
	EXCLUDE_LIST      = "video.exclude.paths"
	MIME_TYPES        = "video.mime.types"
	MIN_VIDEO_SIZE    = "minimum.video.size"
	MAX_WALK_DEPTH    = 4
	HEADER_BYTES_SIZE = 261
)

var (
	excludePaths []string
	mimeTypes    []string
	minFileSize  int64
)

func init() {
	if AppProperties.HasProperty(EXCLUDE_LIST) {
		excludePaths = strings.Split(AppProperties.GetPropertyAsString(EXCLUDE_LIST), ",")
	}
	if AppProperties.HasProperty(MIME_TYPES) {
		mimeTypes = strings.Split(AppProperties.GetPropertyAsString(MIME_TYPES), ",")
	}
	if AppProperties.HasProperty(MIN_VIDEO_SIZE) {
		minFileSize = AppProperties.GetPropertyAsInt64(MIN_VIDEO_SIZE)
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
