package move

import (
	"path/filepath"
	"testing"
	"videosmover/pkg"
	"videosmover/pkg/action"
	"videosmover/pkg/ext/json"
)

func TestMoveAction(t *testing.T) {
	fromPath, fromCleanup := action.SetupTestTmpDir(t, "videosmover_move_test-FROM-")
	defer fromCleanup()
	tvToPath, tvToCleanup := action.SetupTestTmpDir(t, "videosmover_move_test-TV_TO-")
	defer tvToCleanup()
	movieToPath, movieToCleanup := action.SetupTestTmpDir(t, "videosmover_move_test-MOVIE_TO-")
	defer movieToCleanup()

	video1 := action.AddTestVideo(t, fromPath, filepath.Join("Some Movie1", "someMovie.mp4"))
	video2 := action.AddTestVideo(t, fromPath, filepath.Join("Some Show2", "someShow.mp4"))

	videoSeries1 := action.AddTestVideo(t, fromPath, filepath.Join("six.feet.under.720p", "six.feet.under.s01e01.720p.mp4"))
	videoSeries2 := action.AddTestVideo(t, fromPath, filepath.Join("six.feet.under.720p", "six.feet.under.s01e02.720p.mp4"))
	videoSeries3 := action.AddTestVideo(t, fromPath, filepath.Join("six.feet.under.720p", "six.feet.under.s01e03.720p.mp4"))

	video3 := action.AddTestVideo(t, fromPath, filepath.Join("another.movie", "anotherMovie.avi"))
	video3Subs := action.AddTestSubtitles(t, video3, []string{"subtit.srt", filepath.Join("Sub", "anotherMovie.nfo")})

	video4 := action.AddTestVideo(t, fromPath, filepath.Join("hello.movie", "hello.kmv"))
	action.AddTestFile(t, fromPath, filepath.Join("hello.movie", "hello.junk"), 5)

	testData := []action.TestActionData{
		{
			Request: []RequestData{
				{
					Video:    video1,
					Type:     "movie",
					DiskPath: movieToPath,
					OutName:  "Some Movie1",
				},
			},
			Response: []interface{}{},
		},
		{
			Request: []RequestData{
				{
					Video:    video2,
					Type:     "tv",
					DiskPath: tvToPath,
					OutName:  "Some Show2",
				},
			},
			Response: []interface{}{},
		},
		{
			Request: []RequestData{
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
			Response: []interface{}{},
		},
		{
			Request: []RequestData{
				{
					Video:    video3,
					Type:     "movie",
					Subs:     video3Subs,
					DiskPath: movieToPath,
					OutName:  "Another Movie",
				},
			},
			Response: []interface{}{},
		},
		{
			Request: []RequestData{
				{
					Video:    video4,
					Type:     "movie",
					DiskPath: movieToPath,
					OutName:  "Movie With Junk Files",
				},
			},
			Response: []interface{}{},
		},
	}

	jsonCodec := json.NewJsonCodec()
	action.RunTestAction(t, testData, NewAction(action.GetTestActionConfig(), jsonCodec, core.NewTrashMover()), jsonCodec)
}
