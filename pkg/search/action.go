package search

import (
	"encoding/json"
	"github.com/h2non/filetype"
	"github.com/karrick/godirwalk"
	"github.com/lcserny/goutils"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"videosmover/pkg/action"
	"videosmover/pkg/convert"
)

func NewAction() action.Action {
	return &searchAction{}
}

type searchAction struct {
}

func (sa *searchAction) Execute(jsonPayload []byte, config *action.Config) (string, error) {
	var request RequestData
	err := json.Unmarshal(jsonPayload, &request)
	goutils.LogError(err)
	if err != nil {
		return "", err
	}

	if _, err := os.Stat(request.Path); os.IsNotExist(err) {
		goutils.LogError(err)
		return "", err
	}

	realWalkRootPath, _ := filepath.EvalSymlinks(request.Path)
	var resultList []ResponseData
	err = godirwalk.Walk(realWalkRootPath, &godirwalk.Options{
		Unsorted:            true,
		FollowSymbolicLinks: false,
		Callback: func(path string, info *godirwalk.Dirent) error {
			// check path
			for _, exPath := range config.SearchExcludePaths {
				if strings.Contains(path, exPath) {
					return filepath.SkipDir
				}
			}

			if err != nil {
				goutils.LogError(err)
				return nil
			}

			if !info.IsDir() && action.WalkDepthIsAcceptable(realWalkRootPath, path, config.MaxSearchWalkDepth) {
				if isVideo(path, config) {
					resultList = append(resultList, ResponseData{path, findSubtitles(realWalkRootPath, path, config)})
				}
			}
			return nil
		},
	})
	goutils.LogError(err)
	if err != nil {
		return "", err
	}

	if len(resultList) < 1 {
		return "", nil
	}

	sort.Slice(resultList, func(i, j int) bool {
		return resultList[i].Path < resultList[j].Path
	})

	return convert.GetJSONEncodedString(resultList), nil
}

func isVideo(path string, config *action.Config) bool {
	file, err := os.Open(path)
	if err != nil {
		goutils.LogError(err)
		return false
	}
	defer goutils.CloseFile(file)
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
	info, err := os.Stat(path)
	if err != nil {
		goutils.LogError(err)
		return false
	}
	if info.Size() < config.MinimumVideoSize {
		return false
	}

	return true
}

func findSubtitles(rootPath, path string, config *action.Config) []string {
	subs := make([]string, 0)
	pathDir := filepath.Dir(path)
	if rootPath == pathDir {
		return subs
	}

	err := filepath.Walk(pathDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			goutils.LogError(err)
			return nil
		}

		if !info.IsDir() && isSubtitle(path, config.AllowedSubtitleExtensions) {
			subs = append(subs, path)
		}

		return nil
	})
	goutils.LogError(err)

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
