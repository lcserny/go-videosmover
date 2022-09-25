package godirwalk

import (
	"github.com/karrick/godirwalk"
	"path/filepath"
	"strings"
	"videosmover/pkg"
	"videosmover/pkg/action"
)

type videoPathWalker struct {
	config *core.ActionConfig
}

func NewVideoPathWalker(cfg *core.ActionConfig) core.VideoPathWalker {
	return &videoPathWalker{config: cfg}
}

func (vpw videoPathWalker) Walk(root string, vidFn core.VideoCheckingFunc, subFn core.SubtitleFindingFunc) ([]*core.VideoSearchResult, error) {
	resultList := make([]*core.VideoSearchResult, 0)
	err := godirwalk.Walk(root, &godirwalk.Options{
		Unsorted:            true,
		FollowSymbolicLinks: false,
		Callback: func(path string, info *godirwalk.Dirent) error {
			// check path
			for _, exPath := range vpw.config.SearchExcludePaths {
				if strings.Contains(path, exPath) {
					return filepath.SkipDir
				}
			}

			if !info.IsDir() && action.WalkDepthIsAcceptable(root, path, vpw.config.MaxSearchWalkDepth) {
				if vidFn(path) {
					resultList = append(resultList, &core.VideoSearchResult{path, subFn(root, path)})
				}
			}
			return nil
		},
	})
	return resultList, err
}
