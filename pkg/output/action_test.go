package output

import (
	"path/filepath"
	"testing"
	"videosmover/pkg/action"
)

func TestOutputAction(t *testing.T) {
	tmpPath, cleanup := action.SetupTestTmpDir(t, "videosmover_output_test-")
	defer cleanup()
	_ = action.AddTestVideo(t, tmpPath, filepath.Join("The Big Sick (2017)", "video.mp4"))
	_ = action.AddTestVideo(t, tmpPath, filepath.Join("Extras (2005)", "video.mkv"))

	testData := []action.TestActionData{
		{
			Request:  RequestData{Name: "The Lord of the Rings: The Fellowship of <>the Ring (2001)", Type: "movie", SkipOnlineSearch: true},
			Response: ResponseData{[]string{"The Lord Of The Rings The Fellowship Of The Ring (2001)"}, ORIGIN_NAME},
		},
		{
			Request: RequestData{
				Name:             "The Big Sick 2017",
				Type:             "movie",
				SkipOnlineSearch: true,
				DiskPath:         tmpPath,
			},
			Response: ResponseData{[]string{"The Big Sick (2017)"}, ORIGIN_DISK},
		},
		{
			Request: RequestData{
				Name:             "Bodyguard-S01-Series.1--BBC-2018-720p-w.subs-x265-HEVC",
				Type:             "tv",
				SkipOnlineSearch: true,
			},
			Response: ResponseData{[]string{"Bodyguard"}, ORIGIN_NAME},
		},
		{
			Request: RequestData{
				Name:             "1922.1080p.[2017].x264",
				Type:             "movie",
				SkipOnlineSearch: true,
			},
			Response: ResponseData{[]string{"1922"}, ORIGIN_NAME},
		},
		{
			Request: RequestData{
				Name:             "Game.Of.Thrones.s0e10",
				Type:             "tv",
				SkipOnlineSearch: true,
			},
			Response: ResponseData{[]string{"Game Of Thrones"}, ORIGIN_NAME},
		},
		{
			Request: RequestData{
				Name:             "Acrimony.2018.1080p.WEB-DL.DD5.1.H264-FGT",
				Type:             "movie",
				SkipOnlineSearch: true,
			},
			Response: ResponseData{[]string{"Acrimony (2018)"}, ORIGIN_NAME},
		},
		{
			Request: RequestData{
				Name:             "Criminal.Minds.s01e01",
				Type:             "tv",
				SkipOnlineSearch: true,
			},
			Response: ResponseData{[]string{"Criminal Minds"}, ORIGIN_NAME},
		},
		{
			Request: RequestData{
				Name:             "Chicago.PD.S05E21.REPACK.HDTV.x264-KILLERS[rarbg]",
				Type:             "tv",
				SkipOnlineSearch: true,
			},
			Response: ResponseData{[]string{"Chicago PD"}, ORIGIN_NAME},
		},
		{
			Request: RequestData{
				Name:             "House.of.Cards.S06.1080p.NF.WEBRip.DD5.1.x264-NTG[rartv]",
				Type:             "tv",
				SkipOnlineSearch: true,
			},
			Response: ResponseData{[]string{"House Of Cards"}, ORIGIN_NAME},
		},
		{
			Request: RequestData{
				Name:             "House of Cards s06",
				Type:             "tv",
				SkipOnlineSearch: true,
			},
			Response: ResponseData{[]string{"House Of Cards"}, ORIGIN_NAME},
		},
		{
			Request: RequestData{
				Name:             "Extras S02e01-06",
				Type:             "tv",
				SkipOnlineSearch: true,
				DiskPath:         tmpPath,
			},
			Response: ResponseData{[]string{"Extras (2005)"}, ORIGIN_DISK},
		},
		{
			Request: RequestData{
				Name:      "Fight Club",
				Type:      "movie",
				SkipCache: true,
			},
			Response: ResponseData{[]string{
				"Fight Club (1999-10-15)",
			}, ORIGIN_TMDB},
		},
		{
			Request: RequestData{
				Name:      "Game of Thrones",
				Type:      "tv",
				SkipCache: true,
			},
			Response: ResponseData{[]string{
				"Game of Thrones (2011)",
			}, ORIGIN_TMDB},
		},
		{
			Request: RequestData{
				Name:             "Law.&.Order",
				Type:             "tv",
				SkipOnlineSearch: true,
			},
			Response: ResponseData{[]string{"Law And Order"}, ORIGIN_NAME},
		},
	}

	action.RunTestAction(t, testData, NewAction())
}
