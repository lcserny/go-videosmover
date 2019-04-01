package delete

import (
	"fmt"
	"github.com/lcserny/go-videosmover/pkg/convert"
	intest "github.com/lcserny/go-videosmover/pkg/testing"
	"path/filepath"
	"testing"
)

func TestAction(t *testing.T) {
	tmpPath, tmpCleanup := intest.SetupTmpDir(t, "videosmover_delete_test-")
	defer tmpCleanup()

	video1 := intest.AddVideo(t, tmpPath, filepath.Join("MovieFolder", "videoFile.mp4"))
	video2 := intest.AddVideo(t, tmpPath, filepath.Join("MovieFolder2", "videoFile2.mp4"))
	video3 := intest.AddVideo(t, tmpPath, filepath.Join("Downloads", "videoFile3.mp4"))

	testData := []intest.TestActionData{
		{
			Request: []convert.DeleteRequestData{
				{filepath.Dir(video1)},
				{filepath.Dir(video2)},
				{filepath.Dir(video3)},
			},
			Response: []convert.DeleteResponseData{
				{filepath.Dir(video3), []string{
					fmt.Sprintf("Dir '%s' is a restricted path", filepath.Dir(video3)),
				}},
			},
		},
	}

	intest.RunActionTest(t, testData, Action)
}
