package actions

import (
	"encoding/json"
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
		request       SearchRequestData
		response      []SearchResponseData
		expectedError error
	}{
		{
			request: SearchRequestData{path},
			response: []SearchResponseData{
				{video1, video1Subs},
			},
			expectedError: nil,
		},
	}

	// TODO: abstract this?
	for _, search := range searches {
		reqBytes, err := json.Marshal(search.request)
		if err != nil {
			t.Fatalf("Couldn't decode request: %+v", err)
		}
		resBytes, err := json.Marshal(search.response)
		if err != nil {
			t.Fatalf("Couldn't decode response: %+v", err)
		}

		result, err := SearchAction(reqBytes)
		if err != nil {
			t.Fatalf("Error occurred: %+v", err)
		}
		if result != string(resBytes) {
			t.Errorf("Result mismatch, wanted %s, got: %s", string(resBytes), result)
		}
	}
}
