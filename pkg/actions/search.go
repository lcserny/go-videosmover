package actions

import (
	"encoding/json"
	"github.com/h2non/filetype"
	. "github.com/lcserny/go-videosmover/pkg/models"
	. "github.com/lcserny/goutils"
	"os"
	"path/filepath"
	"strings"
)

func SearchAction(jsonPayload []byte, config *ActionConfig) (string, error) {
	var request SearchRequestData
	err := json.Unmarshal(jsonPayload, &request)
	LogError(err)
	if err != nil {
		return "", err
	}

	if _, err := os.Stat(request.Path); os.IsNotExist(err) {
		LogError(err)
		return "", err
	}

	realWalkRootPath, _ := GetRealPath(request.Path)
	var resultList []SearchResponseData
	err = filepath.Walk(realWalkRootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			LogError(err)
			return nil
		}

		if !info.IsDir() && walkDepthIsAcceptable(realWalkRootPath, path, config.maxSearchWalkDepth) {
			if isVideo(path, info, config) {
				resultList = append(resultList, SearchResponseData{path, findSubtitles(realWalkRootPath, path, config)})
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

	return getJSONEncodedString(resultList), nil
}

func isVideo(path string, info os.FileInfo, config *ActionConfig) bool {
	file, err := os.Open(path)
	if err != nil {
		LogError(err)
		return false
	}
	defer CloseFile(file)
	head := make([]byte, config.headerBytesSize)
	n, _ := file.Read(head)
	if n < config.headerBytesSize {
		return false
	}

	// check path
	for _, exPath := range config.searchExcludePaths {
		if strings.Contains(path, exPath) {
			return false
		}
	}

	// check type
	acceptedMime := false
	for _, mType := range config.allowedMimeTypes {
		if filetype.IsMIME(head, mType) {
			acceptedMime = true
			break
		}
	}
	if !acceptedMime && !filetype.IsVideo(head) {
		return false
	}

	// check size
	if info.Size() < config.minVideoFileSize {
		return false
	}

	return true
}

func findSubtitles(rootPath, path string, config *ActionConfig) []string {
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

		if !info.IsDir() && isSubtitle(path, config.allowedSubtitleExts) {
			subs = append(subs, path)
		}

		return nil
	})
	LogError(err)

	return subs
}

func isSubtitle(path string, allowedSubtitleExts []string) bool {
	ext := filepath.Ext(path)
	for _, allowedExt := range allowedSubtitleExts {
		if ext == allowedExt {
			return true
		}
	}
	return false
}
