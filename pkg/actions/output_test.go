package actions

import (
	"github.com/lcserny/go-videosmover/pkg/models"
	"path/filepath"
	"testing"
)

func TestOutputAction(t *testing.T) {
	tmpPath, cleanup := setupTmpDir(t, "videosmover_output_test-")
	defer cleanup()
	_ = addVideo(t, tmpPath, filepath.Join("The Big Sick (2017)", "video.mp4"))

	testData := []testActionData{
		{
			request:  models.OutputRequestData{Name: "The Lord of the Rings: The Fellowship of <>the Ring (2001)", Type: "movie", SkipOnlineSearch: true},
			response: models.OutputResponseData{[]string{"The Lord Of The Rings The Fellowship Of The Ring (2001)"}, models.ORIGIN_NAME},
		},
		{
			request: models.OutputRequestData{
				Name:             "The Big Sick 2017",
				Type:             "movie",
				SkipOnlineSearch: true,
				DiskPath:         tmpPath,
			},
			response: models.OutputResponseData{[]string{"The Big Sick (2017)"}, models.ORIGIN_DISK},
		},
		{
			request: models.OutputRequestData{
				Name:             "Bodyguard-S01-Series.1--BBC-2018-720p-w.subs-x265-HEVC",
				Type:             "tv",
				SkipOnlineSearch: true,
			},
			response: models.OutputResponseData{[]string{"Bodyguard"}, models.ORIGIN_NAME},
		},
	}

	runActionTest(t, testData, OutputAction)
}
