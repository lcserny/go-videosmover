package output

import (
	"path/filepath"
	"testing"
	"videosmover/pkg"
	"videosmover/pkg/action"
	"videosmover/pkg/ext/json"
	"videosmover/pkg/ext/tmdb"
)

func TestOutputAction(t *testing.T) {
	tmpPath, cleanup := action.SetupTestTmpDir(t, "videosmover_output_test-")
	defer cleanup()
	_ = action.AddTestVideo(t, tmpPath, filepath.Join("The Big Sick (2017)", "video.mp4"))
	_ = action.AddTestVideo(t, tmpPath, filepath.Join("Extras (2005)", "video.mkv"))

	testData := []action.TestActionData{
		{
			Request:  RequestData{Name: "The Lord of the Rings: The Fellowship of <>the Ring (2001)", Type: "movie", SkipOnlineSearch: true},
			Response: ResponseData{[]*core.VideoWebResult{{Title: "The Lord Of The Rings The Fellowship Of The Ring (2001)"}}, ORIGIN_NAME},
		},
		{
			Request: RequestData{
				Name:             "The Big Sick 2017",
				Type:             "movie",
				SkipOnlineSearch: true,
				DiskPath:         tmpPath,
			},
			Response: ResponseData{[]*core.VideoWebResult{{Title: "The Big Sick (2017)"}}, ORIGIN_DISK},
		},
		{
			Request: RequestData{
				Name:             "Bodyguard-S01-Series.1--BBC-2018-720p-w.subs-x265-HEVC",
				Type:             "tv",
				SkipOnlineSearch: true,
			},
			Response: ResponseData{[]*core.VideoWebResult{{Title: "Bodyguard"}}, ORIGIN_NAME},
		},
		{
			Request: RequestData{
				Name:             "1922.1080p.[2017].x264",
				Type:             "movie",
				SkipOnlineSearch: true,
			},
			Response: ResponseData{[]*core.VideoWebResult{{Title: "1922"}}, ORIGIN_NAME},
		},
		{
			Request: RequestData{
				Name:             "Game.Of.Thrones.s0e10",
				Type:             "tv",
				SkipOnlineSearch: true,
			},
			Response: ResponseData{[]*core.VideoWebResult{{Title: "Game Of Thrones"}}, ORIGIN_NAME},
		},
		{
			Request: RequestData{
				Name:             "Acrimony.2018.1080p.WEB-DL.DD5.1.H264-FGT",
				Type:             "movie",
				SkipOnlineSearch: true,
			},
			Response: ResponseData{[]*core.VideoWebResult{{Title: "Acrimony (2018)"}}, ORIGIN_NAME},
		},
		{
			Request: RequestData{
				Name:             "Criminal.Minds.s01e01",
				Type:             "tv",
				SkipOnlineSearch: true,
			},
			Response: ResponseData{[]*core.VideoWebResult{{Title: "Criminal Minds"}}, ORIGIN_NAME},
		},
		{
			Request: RequestData{
				Name:             "Chicago.PD.S05E21.REPACK.HDTV.x264-KILLERS[rarbg]",
				Type:             "tv",
				SkipOnlineSearch: true,
			},
			Response: ResponseData{[]*core.VideoWebResult{{Title: "Chicago PD"}}, ORIGIN_NAME},
		},
		{
			Request: RequestData{
				Name:             "House.of.Cards.S06.1080p.NF.WEBRip.DD5.1.x264-NTG[rartv]",
				Type:             "tv",
				SkipOnlineSearch: true,
			},
			Response: ResponseData{[]*core.VideoWebResult{{Title: "House Of Cards"}}, ORIGIN_NAME},
		},
		{
			Request: RequestData{
				Name:             "House of Cards s06",
				Type:             "tv",
				SkipOnlineSearch: true,
			},
			Response: ResponseData{[]*core.VideoWebResult{{Title: "House Of Cards"}}, ORIGIN_NAME},
		},
		{
			Request: RequestData{
				Name:             "Extras S02e01-06",
				Type:             "tv",
				SkipOnlineSearch: true,
				DiskPath:         tmpPath,
			},
			Response: ResponseData{[]*core.VideoWebResult{{Title: "Extras (2005)"}}, ORIGIN_DISK},
		},
		{
			Request: RequestData{
				Name:             "Law.&.Order",
				Type:             "tv",
				SkipOnlineSearch: true,
			},
			Response: ResponseData{[]*core.VideoWebResult{{Title: "Law And Order"}}, ORIGIN_NAME},
		},
	}

	cfg := action.GetTestActionConfig()
	jsonCodec := json.NewJsonCodec()
	videoWebSearcher := tmdb.NewVideoWebSearcher()
	if videoWebSearcher.CanSearch() {
		testData = append(testData, action.TestActionData{
			Request: RequestData{
				Name:      "Game of Thrones",
				Type:      "tv",
				SkipCache: true,
			},
			Response: ResponseData{[]*core.VideoWebResult{
				{
					Title:       "Game of Thrones (2011)",
					Description: "Seven noble families fight for control of the mythical land of Westeros. Friction between the houses leads to full-scale war. All while a very ancient evil awakens in the farthest north. Amidst the war, a neglected military order of misfits, the Night's Watch, is all that stands between the realms of men and icy horrors beyond.",
					PosterURL:   "http://image.tmdb.org/t/p/w92/u3bZgnGQ9T01sWNhyveQz0wH0Hl.jpg",
					Cast:        []string{"Liam Cunningham", "Joe Dempsie", "Emilia Clarke", "John Bradley", "Peter Dinklage"},
				},
			}, ORIGIN_TMDB},
		}, action.TestActionData{
			Request: RequestData{
				Name:      "Fight Club",
				Type:      "movie",
				SkipCache: true,
			},
			Response: ResponseData{[]*core.VideoWebResult{
				{
					Title:       "Fight Club (1999-10-15)",
					Description: "A ticking-time-bomb insomniac and a slippery soap salesman channel primal male aggression into a shocking new form of therapy. Their concept catches on, with underground \"fight clubs\" forming in every town, until an eccentric gets in the way and ignites an out-of-control spiral toward oblivion.",
					PosterURL:   "http://image.tmdb.org/t/p/w92/9dlcOgehxDK4QaC2QDfqwQFbk5C.jpg",
					Cast:        []string{"Edward Norton", "Brad Pitt", "Helena Bonham Carter", "Meat Loaf", "Jared Leto"},
				},
			}, ORIGIN_TMDB},
		})
	}

	action.RunTestAction(t, testData, NewAction(cfg, jsonCodec, videoWebSearcher, &core.MockCacheStore{}), jsonCodec)
}
