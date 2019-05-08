package search

import (
	"github.com/h2non/filetype"
	"github.com/karrick/godirwalk"
	"github.com/lcserny/goutils"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"videosmover/pkg"
	"videosmover/pkg/action"
)

func NewAction(cfg *core.ActionConfig, c core.Codec) action.Action {
	return &searchAction{config: cfg, codec: c}
}

type searchAction struct {
	config *core.ActionConfig
	codec  core.Codec
}

func (sa searchAction) Execute(jsonPayload []byte) (string, error) {
	var request RequestData
	if err := sa.codec.Decode(jsonPayload, &request); err != nil {
		goutils.LogError(err)
		return "", err
	}

	if _, err := os.Stat(request.Path); os.IsNotExist(err) {
		goutils.LogError(err)
		return "", err
	}

	realWalkRootPath, _ := filepath.EvalSymlinks(request.Path)
	var resultList []ResponseData
	err := godirwalk.Walk(realWalkRootPath, &godirwalk.Options{
		Unsorted:            true,
		FollowSymbolicLinks: false,
		Callback: func(path string, info *godirwalk.Dirent) error {
			// check path
			for _, exPath := range sa.config.SearchExcludePaths {
				if strings.Contains(path, exPath) {
					return filepath.SkipDir
				}
			}

			if !info.IsDir() && action.WalkDepthIsAcceptable(realWalkRootPath, path, sa.config.MaxSearchWalkDepth) {
				if sa.isVideo(path) {
					resultList = append(resultList, ResponseData{path, sa.findSubtitles(realWalkRootPath, path)})
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

	return sa.codec.EncodeString(resultList)
}

func (sa searchAction) isVideo(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		goutils.LogError(err)
		return false
	}
	defer goutils.CloseFile(file)
	head := make([]byte, sa.config.HeaderBytesSize)
	n, _ := io.ReadFull(file, head)
	if n < sa.config.HeaderBytesSize {
		return false
	}

	// check type
	acceptedMime := false
	for _, mType := range sa.config.AllowedMIMETypes {
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
	if info.Size() < sa.config.MinimumVideoSize {
		return false
	}

	return true
}

func (sa searchAction) findSubtitles(rootPath, path string) []string {
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

		if !info.IsDir() && sa.isSubtitle(path, sa.config.AllowedSubtitleExtensions) {
			subs = append(subs, path)
		}

		return nil
	})
	goutils.LogError(err)

	return subs
}

func (sa searchAction) isSubtitle(path string, allowedSubtitleExts []string) bool {
	ext := filepath.Ext(path)
	for _, allowedExt := range allowedSubtitleExts {
		if ext == allowedExt {
			return true
		}
	}
	return false
}
