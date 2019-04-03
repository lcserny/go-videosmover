package delete

import (
	"fmt"
	"github.com/lcserny/go-videosmover/pkg/models"
	"path/filepath"
	"testing"
)

func TestDeleteAction(t *testing.T) {
	tmpPath, tmpCleanup := setupTmpDir(t, "videosmover_delete_test-")
	defer tmpCleanup()

	video1 := addVideo(t, tmpPath, filepath.Join("MovieFolder", "videoFile.mp4"))
	video2 := addVideo(t, tmpPath, filepath.Join("MovieFolder2", "videoFile2.mp4"))
	video3 := addVideo(t, tmpPath, filepath.Join("Downloads", "videoFile3.mp4"))

	testData := []testActionData{
		{
			request: []models.DeleteRequestData{
				{filepath.Dir(video1)},
				{filepath.Dir(video2)},
				{filepath.Dir(video3)},
			},
			response: []models.DeleteResponseData{
				{filepath.Dir(video3), []string{
					fmt.Sprintf("Dir '%s' is a restricted path", filepath.Dir(video3)),
				}},
			},
		},
	}

	runActionTest(t, testData, DeleteAction)
}
