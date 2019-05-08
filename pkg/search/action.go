package search

import (
	"github.com/lcserny/goutils"
	"io"
	"os"
	"path/filepath"
	"sort"
	"videosmover/pkg"
	"videosmover/pkg/action"
)

func NewAction(cfg *core.ActionConfig, c core.Codec, mc core.VideoChecker, pw core.VideoPathWalker) action.Action {
	return &searchAction{config: cfg, codec: c, mimeChecker: mc, walker: pw}
}

type searchAction struct {
	config      *core.ActionConfig
	codec       core.Codec
	mimeChecker core.VideoChecker
	walker      core.VideoPathWalker
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
	resultList, err := sa.walker.Walk(realWalkRootPath, sa.isVideo, sa.findSubtitles)

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
	if ok := sa.mimeChecker.IsVideo(head); !ok {
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

		if !info.IsDir() && sa.mimeChecker.IsSubtitle(path) {
			subs = append(subs, path)
		}

		return nil
	})
	goutils.LogError(err)

	return subs
}
