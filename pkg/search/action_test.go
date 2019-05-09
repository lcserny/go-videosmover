package search

import (
	"path/filepath"
	"testing"
	"videosmover/pkg"
	"videosmover/pkg/action"
	"videosmover/pkg/ext/godirwalk"
	"videosmover/pkg/ext/h2non"
	"videosmover/pkg/ext/json"
)

func TestSearchAction(t *testing.T) {
	path, cleanup := action.SetupTestTmpDir(t, "videosmover_search_test-")
	defer cleanup()
	video1 := action.AddTestVideo(t, path, filepath.Join("Video1 Folder", "video.mp4"))
	video1Subs := action.AddTestSubtitles(t, video1, []string{filepath.Join("Sub", "subtitle.srt")})
	video2 := action.AddTestVideo(t, path, "video2.mp4")
	_ = action.AddTestSubtitles(t, video2, []string{"subtitle.srt"})
	video3 := action.AddTestVideo(t, path, filepath.Join("A-Video1", "video1.mp4"))
	_ = action.AddTestFile(t, path, "nonVideo.txt", 1)
	hiddenVideo4 := action.AddTestVideo(t, path, "HiddenVideo.txt")
	_ = action.AddTestFile(t, path, "bigEmptyFile.nfo", 300)
	_ = action.AddTestVideo(t, path, filepath.Join("Programming Stuff", "disallowedPathVideo.mp4"))

	testData := []action.TestActionData{
		{
			Request: RequestData{path},
			Response: []core.VideoSearchResult{
				{video3, make([]string, 0)},
				{hiddenVideo4, make([]string, 0)},
				{video1, video1Subs},
				{video2, make([]string, 0)},
			},
		},
	}

	cfg := action.GetTestActionConfig()
	jsonCodec := json.NewJsonCodec()
	mimeChecker := h2non.NewVideoChecker(cfg)
	videoPathWalker := godirwalk.NewVideoPathWalker(cfg)
	action.RunTestAction(t, testData, NewAction(cfg, jsonCodec, mimeChecker, videoPathWalker), jsonCodec)
}
