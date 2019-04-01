package search

import (
	"encoding/json"
	"github.com/h2non/filetype"
	"github.com/lcserny/go-videosmover/pkg/action"
	"github.com/lcserny/go-videosmover/pkg/convert"
	. "github.com/lcserny/goutils"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func Action(jsonPayload []byte, config *convert.ActionConfig) (string, error) {
	var request convert.SearchRequestData
	err := json.Unmarshal(jsonPayload, &request)
	LogError(err)
	if err != nil {
		return "", err
	}

	if _, err := os.Stat(request.Path); os.IsNotExist(err) {
		LogError(err)
		return "", err
	}

	realWalkRootPath, _ := filepath.EvalSymlinks(request.Path)
	var resultList []convert.SearchResponseData
	err = filepath.Walk(realWalkRootPath, func(path string, info os.FileInfo, err error) error {
		// check path
		for _, exPath := range config.SearchExcludePaths {
			if strings.Contains(path, exPath) {
				return filepath.SkipDir
			}
		}

		if err != nil {
			LogError(err)
			return nil
		}

		if !info.IsDir() && action.WalkDepthIsAcceptable(realWalkRootPath, path, config.MaxSearchWalkDepth) {
			if isVideo(path, info, config) {
				resultList = append(resultList, convert.SearchResponseData{path, findSubtitles(realWalkRootPath, path, config)})
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

	return convert.GetJSONEncodedString(resultList), nil
}

func isVideo(path string, info os.FileInfo, config *convert.ActionConfig) bool {
	file, err := os.Open(path)
	if err != nil {
		LogError(err)
		return false
	}
	defer CloseFile(file)
	head := make([]byte, config.HeaderBytesSize)
	n, _ := io.ReadFull(file, head)
	if n < config.HeaderBytesSize {
		return false
	}

	// check type
	acceptedMime := false
	for _, mType := range config.AllowedMIMETypes {
		if filetype.IsMIME(head, mType) {
			acceptedMime = true
			break
		}
	}
	if !acceptedMime && !filetype.IsVideo(head) {
		return false
	}

	// check size
	if info.Size() < config.MinimumVideoSize {
		return false
	}

	return true
}

func findSubtitles(rootPath, path string, config *convert.ActionConfig) []string {
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

		if !info.IsDir() && isSubtitle(path, config.AllowedSubtitleExtensions) {
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
