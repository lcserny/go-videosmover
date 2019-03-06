package actions

import (
	. "github.com/lcserny/go-videosmover/src/shared"
	"path/filepath"
	"testing"
)

func init() {
	minFileSize = 256
}

func TestSearchAction(t *testing.T) {
	path, cleanup := setupTmpDir(t, "videosmover_search_test-")
	defer cleanup()
	video1 := addVideo(t, path, filepath.Join("Video1 Folder", "video.mp4"))
	video1Subs := addSubtitles(t, video1, []string{filepath.Join("Sub", "subtitle.srt")})

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
