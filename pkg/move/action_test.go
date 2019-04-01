package move

import (
	"github.com/lcserny/go-videosmover/pkg/convert"
	intest "github.com/lcserny/go-videosmover/pkg/testing"
	"path/filepath"
	"testing"
)

func TestAction(t *testing.T) {
	fromPath, fromCleanup := intest.SetupTmpDir(t, "videosmover_move_test-FROM-")
	defer fromCleanup()
	tvToPath, tvToCleanup := intest.SetupTmpDir(t, "videosmover_move_test-TV_TO-")
	defer tvToCleanup()
	movieToPath, movieToCleanup := intest.SetupTmpDir(t, "videosmover_move_test-MOVIE_TO-")
	defer movieToCleanup()

	video1 := intest.AddVideo(t, fromPath, filepath.Join("Some Movie1", "someMovie.mp4"))
	video2 := intest.AddVideo(t, fromPath, filepath.Join("Some Show2", "someShow.mp4"))

	videoSeries1 := intest.AddVideo(t, fromPath, filepath.Join("six.feet.under.720p", "six.feet.under.s01e01.720p.mp4"))
	videoSeries2 := intest.AddVideo(t, fromPath, filepath.Join("six.feet.under.720p", "six.feet.under.s01e02.720p.mp4"))
	videoSeries3 := intest.AddVideo(t, fromPath, filepath.Join("six.feet.under.720p", "six.feet.under.s01e03.720p.mp4"))

	video3 := intest.AddVideo(t, fromPath, filepath.Join("another.movie", "anotherMovie.avi"))
	video3Subs := intest.AddSubtitles(t, video3, []string{"subtit.srt", filepath.Join("Sub", "anotherMovie.nfo")})

	video4 := intest.AddVideo(t, fromPath, filepath.Join("hello.movie", "hello.kmv"))
	intest.AddFile(t, fromPath, filepath.Join("hello.movie", "hello.junk"), 5)

	testData := []intest.TestActionData{
		{
			Request: []convert.MoveRequestData{
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
			Request: []convert.MoveRequestData{
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
			Request: []convert.MoveRequestData{
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
			Request: []convert.MoveRequestData{
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
			Request: []convert.MoveRequestData{
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

	intest.RunActionTest(t, testData, Action)
}
