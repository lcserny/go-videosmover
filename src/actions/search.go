package actions

import (
	"encoding/json"
	"github.com/h2non/filetype"
	. "github.com/lcserny/goutils"
	"os"
	"path/filepath"
	"strings"
)

const (
	EXCLUDE_LIST_FILE      = "search_exclude_paths"
	MIME_TYPES_FILE        = "allowed_mime_types"
	ALLOWED_SUBTITLES_FILE = "allowed_subtitle_exts"
	MIN_VIDEO_SIZE_KEY     = "minimum.video.size"
	MAX_WALK_DEPTH         = 4
	HEADER_BYTES_SIZE      = 261
)

var (
	excludePaths        []string
	mimeTypes           []string
	allowedSubtitleExts []string
	minFileSize         int64
)

type SearchRequestData struct {
	Path string `json:"path"`
}

type SearchResponseData struct {
	Path      string   `json:"path"`
	Subtitles []string `json:"subtitles"`
}

func init() {
	excludePathsContent, err := configFolder.FindString(EXCLUDE_LIST_FILE)
	LogError(err)
	excludePaths = GetLinesFromString(excludePathsContent)

	mimeTypesContent, err := configFolder.FindString(MIME_TYPES_FILE)
	LogError(err)
	mimeTypes = GetLinesFromString(mimeTypesContent)

	allowedSubtitleExtsContent, err := configFolder.FindString(ALLOWED_SUBTITLES_FILE)
	LogError(err)
	allowedSubtitleExts = GetLinesFromString(allowedSubtitleExtsContent)

	if appProperties.HasProperty(MIN_VIDEO_SIZE_KEY) {
		minFileSize = appProperties.GetPropertyAsInt64(MIN_VIDEO_SIZE_KEY)
	}
}

func SearchAction(jsonPayload []byte) (string, error) {
	var request SearchRequestData
	err := json.Unmarshal(jsonPayload, &request)
	LogError(err)
	if err != nil {
		return "", err
	}

	realWalkRootPath, _ := GetRealPath(request.Path)
	var resultList []SearchResponseData
	err = filepath.Walk(realWalkRootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			LogError(err)
			return nil
		}

		if !info.IsDir() && walkDepthIsAcceptable(realWalkRootPath, path, MAX_WALK_DEPTH) {
			if isVideo(path, info) {
				resultList = append(resultList, SearchResponseData{path, findSubtitles(realWalkRootPath, path, info)})
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

func findSubtitles(rootPath, path string, info os.FileInfo) []string {
	subs := make([]string, 0)
	pathDir := filepath.Dir(path)
	if rootPath == pathDir {
		return subs
	}

	err := filepath.Walk(pathDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			LogError(err)
			return nil
		}

		if !info.IsDir() && isSubtitle(path, info) {
			subs = append(subs, path)
		}

		return nil
	})
	LogError(err)

	return subs
}

func isSubtitle(path string, info os.FileInfo) bool {
	ext := filepath.Ext(path)
	for _, allowedExt := range allowedSubtitleExts {
		if ext == allowedExt {
			return true
		}
	}
	return false
}
