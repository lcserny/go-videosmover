package videosmover

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
	EXCLUDE_LIST   = "video.exclude.paths"
	MIME_TYPES     = "video.mime.types"
	MIN_VIDEO_SIZE = "minimum.video.size"
	MAX_WALK_DEPTH = 3
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

	var resultList []ResponseSearchData
	err = filepath.Walk(request.Path, func(path string, info os.FileInfo, err error) error {
		LogError(err)
		if !info.IsDir() && walkDepthIsAcceptable(request.Path, path, MAX_WALK_DEPTH) {
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
	// check path
	for _, exPath := range excludePaths {
		if strings.Contains(path, exPath) {
			return false
		}
	}

	// check type
	file, _ := os.Open(path)
	head := make([]byte, 261)
	_, err := file.Read(head)
	LogError(err)
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
