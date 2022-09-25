package core

type VideoCheckingFunc func(videoPath string) bool
type SubtitleFindingFunc func(rootPath, path string) []string

type VideoSearchResult struct {
	Path      string   `json:"path"`
	Subtitles []string `json:"subtitles"`
}

type VideoPathWalker interface {
	Walk(root string, vidFn VideoCheckingFunc, subFn SubtitleFindingFunc) ([]*VideoSearchResult, error)
}
