package delete

import (
	"fmt"
	"path/filepath"
	"testing"
	"videosmover/pkg/action"
	"videosmover/pkg/json"
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

	jsonCodec := json.NewJsonCodec()
	action.RunTestAction(t, testData, NewAction(jsonCodec), jsonCodec)
}
