package tmdb

import (
	"github.com/lcserny/goutils"
	"github.com/ryanbradynd05/go-tmdb"
	"os"
	"videosmover/pkg"
)

type videoWebSearcher struct {
	tmdbAPI *tmdb.TMDb
}

func NewVideoWebSearcher() core.VideoWebSearcher {
	vws := videoWebSearcher{}
	if key, exists := os.LookupEnv("TMDB_API_KEY"); exists {
		vws.tmdbAPI = tmdb.Init(tmdb.Config{key, false, nil})
	}
	return &vws
}

func (vws videoWebSearcher) CanSearch() bool {
	return vws.tmdbAPI != nil
}

func (vws videoWebSearcher) SearchMovies(name string, year int) ([]*core.VideoWebResult, error) {
	options := map[string]string{"page": "1", "language": "en"}
	if year > 0 {
		options["year"] = string(year)
	}

	results, err := vws.tmdbAPI.SearchMovie(name, options)
	if err != nil {
		goutils.LogError(err)
		return nil, err
	}

	searchedList := make([]*core.VideoWebResult, 0)
	for i := 0; i < len(results.Results); i++ {
		movie := results.Results[i]
		searchedList = append(searchedList, &core.VideoWebResult{Title: movie.Title, ReleaseDate: movie.ReleaseDate})
	}

	return searchedList, nil
}

func (vws videoWebSearcher) SearchTVSeries(name string, year int) ([]*core.VideoWebResult, error) {
	options := map[string]string{"page": "1", "language": "en"}
	if year > 0 {
		options["first_air_date_year"] = string(year)
	}

	results, err := vws.tmdbAPI.SearchTv(name, options)
	if err != nil {
		goutils.LogError(err)
		return nil, err
	}

	searchedList := make([]*core.VideoWebResult, 0)
	for i := 0; i < len(results.Results); i++ {
		tvShow := results.Results[i]
		searchedList = append(searchedList, &core.VideoWebResult{Title: tvShow.Name, ReleaseDate: tvShow.FirstAirDate})
	}

	return searchedList, nil
}
