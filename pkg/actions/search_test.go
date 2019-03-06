package actions

import (
	. "github.com/lcserny/go-videosmover/pkg/models"
	"path/filepath"
	"testing"
)

var testSearchActionConfig *ActionConfig

func TestSearchAction(t *testing.T) {
	path, cleanup := setupTmpDir(t, "videosmover_search_test-")
	defer cleanup()
	video1 := addVideo(t, path, filepath.Join("Video1 Folder", "video.mp4"))
	video1Subs := addSubtitles(t, video1, []string{filepath.Join("Sub", "subtitle.srt")})
	video2 := addVideo(t, path, "video2.mp4")
	_ = addSubtitles(t, video2, []string{"subtitle.srt"})

	searches := []struct {
		request  SearchRequestData
		response []SearchResponseData
	}{
		{
			request: SearchRequestData{path},
			response: []SearchResponseData{
				{video1, video1Subs},
				{video2, make([]string, 0)},
			},
		},
	}

	config := getTestActionConfig()

	for _, search := range searches {
		reqBytes := getJSONBytesForTest(t, search.request)
		resString := getJSONStringForTest(t, search.response)

		result, err := SearchAction(reqBytes, config)
		if err != nil {
			t.Fatalf("Error occurred: %+v", err)
		}

		if result != resString {
			t.Errorf("Result mismatch:\nwanted %s\ngot: %s", resString, result)
		}
	}
}

func getTestActionConfig() *ActionConfig {
	if testSearchActionConfig == nil {
		testSearchActionConfig = GenerateActionConfig("commander_test.properties")
	}
	return testSearchActionConfig
}
