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
	_ = addVideo(t, tmpPath, filepath.Join("Extras (2005)", "video.mkv"))

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
		{
			request: models.OutputRequestData{
				Name:             "1922.1080p.[2017].x264",
				Type:             "movie",
				SkipOnlineSearch: true,
			},
			response: models.OutputResponseData{[]string{"1922"}, models.ORIGIN_NAME},
		},
		{
			request: models.OutputRequestData{
				Name:             "Game.Of.Thrones.s0e10",
				Type:             "tv",
				SkipOnlineSearch: true,
			},
			response: models.OutputResponseData{[]string{"Game Of Thrones"}, models.ORIGIN_NAME},
		},
		{
			request: models.OutputRequestData{
				Name:             "Acrimony.2018.1080p.WEB-DL.DD5.1.H264-FGT",
				Type:             "movie",
				SkipOnlineSearch: true,
			},
			response: models.OutputResponseData{[]string{"Acrimony (2018)"}, models.ORIGIN_NAME},
		},
		{
			request: models.OutputRequestData{
				Name:             "Criminal.Minds.s01e01",
				Type:             "tv",
				SkipOnlineSearch: true,
			},
			response: models.OutputResponseData{[]string{"Criminal Minds"}, models.ORIGIN_NAME},
		},
		{
			request: models.OutputRequestData{
				Name:             "Chicago.PD.S05E21.REPACK.HDTV.x264-KILLERS[rarbg]",
				Type:             "tv",
				SkipOnlineSearch: true,
			},
			response: models.OutputResponseData{[]string{"Chicago PD"}, models.ORIGIN_NAME},
		},
		{
			request: models.OutputRequestData{
				Name:             "House.of.Cards.S06.1080p.NF.WEBRip.DD5.1.x264-NTG[rartv]",
				Type:             "tv",
				SkipOnlineSearch: true,
			},
			response: models.OutputResponseData{[]string{"House Of Cards"}, models.ORIGIN_NAME},
		},
		{
			request: models.OutputRequestData{
				Name:             "House of Cards s06",
				Type:             "tv",
				SkipOnlineSearch: true,
			},
			response: models.OutputResponseData{[]string{"House Of Cards"}, models.ORIGIN_NAME},
		},
		{
			request: models.OutputRequestData{
				Name:             "Extras S02e01-06",
				Type:             "tv",
				SkipOnlineSearch: true,
				DiskPath:         tmpPath,
			},
			response: models.OutputResponseData{[]string{"Extras (2005)"}, models.ORIGIN_DISK},
		},
		{
			request: models.OutputRequestData{
				Name:      "Fight Club",
				Type:      "movie",
				SkipCache: true,
			},
			response: models.OutputResponseData{[]string{
				"Fight Club (1999)",
				"Female Fight Club (2017)",
				"Fight Club Members Only (2006)",
				"Clubbed (2008)",
				"Zombie Fight Club (2014)",
				"Superhero Fight Club (2015)",
				"Jurassic Fight Club (2008)",
				"Lure Teen Fight Club (2010)",
				"Fight club camp kusse (2005)",
				"Florence Fight Club (2015)",
			}, models.ORIGIN_TMDB},
		},
	}

	runActionTest(t, testData, OutputAction)
}
