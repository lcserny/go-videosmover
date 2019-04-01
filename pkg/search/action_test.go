package search

import (
	"github.com/lcserny/go-videosmover/pkg/convert"
	intest "github.com/lcserny/go-videosmover/pkg/testing"
	"path/filepath"
	"testing"
)

func TestAction(t *testing.T) {
	path, cleanup := intest.SetupTmpDir(t, "videosmover_search_test-")
	defer cleanup()
	video1 := intest.AddVideo(t, path, filepath.Join("Video1 Folder", "video.mp4"))
	video1Subs := intest.AddSubtitles(t, video1, []string{filepath.Join("Sub", "subtitle.srt")})
	video2 := intest.AddVideo(t, path, "video2.mp4")
	_ = intest.AddSubtitles(t, video2, []string{"subtitle.srt"})
	video3 := intest.AddVideo(t, path, filepath.Join("A-Video1", "video1.mp4"))
	_ = intest.AddFile(t, path, "nonVideo.txt", 1)
	hiddenVideo4 := intest.AddVideo(t, path, "HiddenVideo.txt")
	_ = intest.AddFile(t, path, "bigEmptyFile.nfo", 300)
	_ = intest.AddVideo(t, path, filepath.Join("Programming Stuff", "disallowedPathVideo.mp4"))

	testData := []intest.TestActionData{
		{
			Request: convert.SearchRequestData{path},
			Response: []convert.SearchResponseData{
				{video3, make([]string, 0)},
				{hiddenVideo4, make([]string, 0)},
				{video1, video1Subs},
				{video2, make([]string, 0)},
			},
		},
	}

	intest.RunActionTest(t, testData, Action)
}
