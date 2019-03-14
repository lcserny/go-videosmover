package actions

import (
	"github.com/lcserny/go-videosmover/pkg/models"
	"path/filepath"
	"testing"
)

func TestMoveAction(t *testing.T) {
	fromPath, fromCleanup := setupTmpDir(t, "videosmover_move_test-FROM-")
	defer fromCleanup()
	tvToPath, tvToCleanup := setupTmpDir(t, "videosmover_move_test-TV_TO-")
	defer tvToCleanup()
	movieToPath, movieToCleanup := setupTmpDir(t, "videosmover_move_test-MOVIE_TO-")
	defer movieToCleanup()

	video1 := addVideo(t, fromPath, filepath.Join("Some Movie1", "someMovie.mp4"))
	video2 := addVideo(t, fromPath, filepath.Join("Some Show2", "someShow.mp4"))

	testData := []testActionData{
		{
			request: []models.MoveRequestData{
				{
					Video:    video1,
					Type:     "movie",
					DiskPath: movieToPath,
					OutName:  "Some Movie1",
				},
			},
			response: []interface{}{},
		},
		{
			request: []models.MoveRequestData{
				{
					Video:    video2,
					Type:     "tv",
					DiskPath: tvToPath,
					OutName:  "Some Show2",
				},
			},
			response: []interface{}{},
		},
	}

	runActionTest(t, testData, MoveAction)
}
