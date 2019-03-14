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

	videoSeries1 := addVideo(t, fromPath, filepath.Join("six.feet.under.720p", "six.feet.under.s01e01.720p.mp4"))
	videoSeries2 := addVideo(t, fromPath, filepath.Join("six.feet.under.720p", "six.feet.under.s01e02.720p.mp4"))
	videoSeries3 := addVideo(t, fromPath, filepath.Join("six.feet.under.720p", "six.feet.under.s01e03.720p.mp4"))

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
		{
			request: []models.MoveRequestData{
				{
					Video:    videoSeries1,
					Type:     "tv",
					DiskPath: tvToPath,
					OutName:  "Six Feet Under",
				},
				{
					Video:    videoSeries2,
					Type:     "tv",
					DiskPath: tvToPath,
					OutName:  "Six Feet Under",
				},
				{
					Video:    videoSeries3,
					Type:     "tv",
					DiskPath: tvToPath,
					OutName:  "Six Feet Under",
				},
			},
			response: []interface{}{},
		},
	}

	runActionTest(t, testData, MoveAction)
}
