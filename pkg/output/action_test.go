package output

import (
	"github.com/lcserny/go-videosmover/pkg/convert"
	intest "github.com/lcserny/go-videosmover/pkg/testing"
	"path/filepath"
	"testing"
)

func TestAction(t *testing.T) {
	tmpPath, cleanup := intest.SetupTmpDir(t, "videosmover_output_test-")
	defer cleanup()
	_ = intest.AddVideo(t, tmpPath, filepath.Join("The Big Sick (2017)", "video.mp4"))
	_ = intest.AddVideo(t, tmpPath, filepath.Join("Extras (2005)", "video.mkv"))

	testData := []intest.TestActionData{
		{
			Request:  convert.OutputRequestData{Name: "The Lord of the Rings: The Fellowship of <>the Ring (2001)", Type: "movie", SkipOnlineSearch: true},
			Response: convert.OutputResponseData{[]string{"The Lord Of The Rings The Fellowship Of The Ring (2001)"}, ORIGIN_NAME},
		},
		{
			Request: convert.OutputRequestData{
				Name:             "The Big Sick 2017",
				Type:             "movie",
				SkipOnlineSearch: true,
				DiskPath:         tmpPath,
			},
			Response: convert.OutputResponseData{[]string{"The Big Sick (2017)"}, ORIGIN_DISK},
		},
		{
			Request: convert.OutputRequestData{
				Name:             "Bodyguard-S01-Series.1--BBC-2018-720p-w.subs-x265-HEVC",
				Type:             "tv",
				SkipOnlineSearch: true,
			},
			Response: convert.OutputResponseData{[]string{"Bodyguard"}, ORIGIN_NAME},
		},
		{
			Request: convert.OutputRequestData{
				Name:             "1922.1080p.[2017].x264",
				Type:             "movie",
				SkipOnlineSearch: true,
			},
			Response: convert.OutputResponseData{[]string{"1922"}, ORIGIN_NAME},
		},
		{
			Request: convert.OutputRequestData{
				Name:             "Game.Of.Thrones.s0e10",
				Type:             "tv",
				SkipOnlineSearch: true,
			},
			Response: convert.OutputResponseData{[]string{"Game Of Thrones"}, ORIGIN_NAME},
		},
		{
			Request: convert.OutputRequestData{
				Name:             "Acrimony.2018.1080p.WEB-DL.DD5.1.H264-FGT",
				Type:             "movie",
				SkipOnlineSearch: true,
			},
			Response: convert.OutputResponseData{[]string{"Acrimony (2018)"}, ORIGIN_NAME},
		},
		{
			Request: convert.OutputRequestData{
				Name:             "Criminal.Minds.s01e01",
				Type:             "tv",
				SkipOnlineSearch: true,
			},
			Response: convert.OutputResponseData{[]string{"Criminal Minds"}, ORIGIN_NAME},
		},
		{
			Request: convert.OutputRequestData{
				Name:             "Chicago.PD.S05E21.REPACK.HDTV.x264-KILLERS[rarbg]",
				Type:             "tv",
				SkipOnlineSearch: true,
			},
			Response: convert.OutputResponseData{[]string{"Chicago PD"}, ORIGIN_NAME},
		},
		{
			Request: convert.OutputRequestData{
				Name:             "House.of.Cards.S06.1080p.NF.WEBRip.DD5.1.x264-NTG[rartv]",
				Type:             "tv",
				SkipOnlineSearch: true,
			},
			Response: convert.OutputResponseData{[]string{"House Of Cards"}, ORIGIN_NAME},
		},
		{
			Request: convert.OutputRequestData{
				Name:             "House of Cards s06",
				Type:             "tv",
				SkipOnlineSearch: true,
			},
			Response: convert.OutputResponseData{[]string{"House Of Cards"}, ORIGIN_NAME},
		},
		{
			Request: convert.OutputRequestData{
				Name:             "Extras S02e01-06",
				Type:             "tv",
				SkipOnlineSearch: true,
				DiskPath:         tmpPath,
			},
			Response: convert.OutputResponseData{[]string{"Extras (2005)"}, ORIGIN_DISK},
		},
		{
			Request: convert.OutputRequestData{
				Name:      "Fight Club",
				Type:      "movie",
				SkipCache: true,
			},
			Response: convert.OutputResponseData{[]string{
				"Fight Club (1999-10-15)",
			}, ORIGIN_TMDB},
		},
		{
			Request: convert.OutputRequestData{
				Name:      "Game of Thrones",
				Type:      "tv",
				SkipCache: true,
			},
			Response: convert.OutputResponseData{[]string{
				"Game of Thrones (2011)",
			}, ORIGIN_TMDB},
		},
		{
			Request: convert.OutputRequestData{
				Name:             "Law.&.Order",
				Type:             "tv",
				SkipOnlineSearch: true,
			},
			Response: convert.OutputResponseData{[]string{"Law And Order"}, ORIGIN_NAME},
		},
	}

	intest.RunActionTest(t, testData, Action)
}
