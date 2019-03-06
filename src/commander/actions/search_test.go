package actions

import (
	. "github.com/lcserny/go-videosmover/src/shared"
	"path/filepath"
	"testing"
)

func TestSearchAction(t *testing.T) {
	minFileSize = 256

	var (
		prefix         = "videosmover_search_test-"
		videoName1     = filepath.Join("Video1 Folder", "video.mp4")
		video1SubNames = []string{filepath.Join("Sub", "subtitle.srt")}
	)

	path, cleanup := setupTmpDir(t, prefix)
	defer cleanup()
	video1 := addVideo(t, path, videoName1)
	video1Subs := addSubtitles(t, video1, video1SubNames)

	searches := []struct {
		request  SearchRequestData
		response []SearchResponseData
	}{
		{
			request: SearchRequestData{path},
			response: []SearchResponseData{
				{video1, video1Subs},
			},
		},
	}

	for _, search := range searches {
		reqBytes := getJSONBytes(t, search.request)
		resString := getJSONString(t, search.response)

		result, err := SearchAction(reqBytes)
		if err != nil {
			t.Fatalf("Error occurred: %+v", err)
		}

		if result != resString {
			t.Errorf("Result mismatch:\nwanted %s\ngot: %s", resString, result)
		}
	}
}
