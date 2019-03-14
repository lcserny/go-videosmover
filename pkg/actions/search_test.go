package actions

import (
	. "github.com/lcserny/go-videosmover/pkg/models"
	"path/filepath"
	"testing"
)

func TestSearchAction(t *testing.T) {
	path, cleanup := setupTmpDir(t, "videosmover_search_test-")
	defer cleanup()
	video1 := addVideo(t, path, filepath.Join("Video1 Folder", "video.mp4"))
	video1Subs := addSubtitles(t, video1, []string{filepath.Join("Sub", "subtitle.srt")})
	video2 := addVideo(t, path, "video2.mp4")
	_ = addSubtitles(t, video2, []string{"subtitle.srt"})

	searches := []testActionData{
		{
			request: SearchRequestData{path},
			response: []SearchResponseData{
				{video1, video1Subs},
				{video2, make([]string, 0)},
			},
		},
	}

	runActionTest(t, searches, SearchAction)
}
