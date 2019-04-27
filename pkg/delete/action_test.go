package delete

import (
	"fmt"
	"github.com/lcserny/go-videosmover/pkg/action"
	"path/filepath"
	"testing"
)

func TestDeleteAction(t *testing.T) {
	tmpPath, tmpCleanup := action.SetupTestTmpDir(t, "videosmover_delete_test-")
	defer tmpCleanup()

	video1 := action.AddTestVideo(t, tmpPath, filepath.Join("MovieFolder", "videoFile.mp4"))
	video2 := action.AddTestVideo(t, tmpPath, filepath.Join("MovieFolder2", "videoFile2.mp4"))
	video3 := action.AddTestVideo(t, tmpPath, filepath.Join("Downloads", "videoFile3.mp4"))

	testData := []action.TestActionData{
		{
			Request: []RequestData{
				{filepath.Dir(video1)},
				{filepath.Dir(video2)},
				{filepath.Dir(video3)},
			},
			Response: []ResponseData{
				{filepath.Dir(video3), []string{
					fmt.Sprintf("Dir '%s' is a restricted path", filepath.Dir(video3)),
				}},
			},
		},
	}

	action.RunTestAction(t, testData, NewAction())
}
