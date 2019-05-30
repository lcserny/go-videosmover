package tmdb

import (
	"fmt"
	"github.com/lcserny/goutils"
	"github.com/ryanbradynd05/go-tmdb"
	"os"
	"regexp"
	"strings"
	"videosmover/pkg"
)

type videoWebSearcher struct {
	tmdbAPI       *tmdb.TMDb
	posterPattern string
}

func NewVideoWebSearcher() core.VideoWebSearcher {
	vws := videoWebSearcher{posterPattern: "http://image.tmdb.org/t/p/w92%s"}
	if key, exists := os.LookupEnv("TMDB_API_KEY"); exists {
		vws.tmdbAPI = tmdb.Init(tmdb.Config{key, false, nil})
	}
	return &vws
}

func (vws videoWebSearcher) CanSearch() bool {
	return vws.tmdbAPI != nil
}

func (vws videoWebSearcher) SearchMovies(name string, year, maxResCount int, specialCharsRegex *regexp.Regexp) ([]*core.VideoWebResult, bool) {
	options := make(map[string]string)
	if year > 0 {
		options["year"] = string(year)
	}

	results, err := vws.tmdbAPI.SearchMovie(name, options)
	if err != nil {
		goutils.LogError(err)
		return nil, false
	}

	searchedList := make([]*core.VideoWebResult, 0)
	for i := 0; i < len(results.Results); i++ {
		movie, err := vws.tmdbAPI.GetMovieInfo(results.Results[i].ID, nil)
		if err != nil {
			goutils.LogError(err)
			continue
		}

		title := strings.ReplaceAll(movie.Title, "&", "and")
		title = specialCharsRegex.ReplaceAllString(title, "")
		if movie.ReleaseDate != "" {
			title += " (" + movie.ReleaseDate + ")"
		}
		searchedList = append(searchedList, &core.VideoWebResult{
			Title:       title,
			Description: movie.Overview,
			PosterURL:   fmt.Sprintf(vws.posterPattern, movie.PosterPath),
			Cast:        vws.generateMovieCastNames(movie.Credits.Cast),
		})
	}

	if len(searchedList) < 1 {
		return nil, false
	}

	return searchedList, true
}

func (vws videoWebSearcher) SearchTVSeries(name string, year, maxResCount int, specialCharsRegex *regexp.Regexp) ([]*core.VideoWebResult, bool) {
	options := make(map[string]string)
	if year > 0 {
		options["first_air_date_year"] = string(year)
	}

	results, err := vws.tmdbAPI.SearchTv(name, options)
	if err != nil {
		goutils.LogError(err)
		return nil, false
	}

	searchedList := make([]*core.VideoWebResult, 0)
	for i := 0; i < len(results.Results); i++ {
		tvShow, err := vws.tmdbAPI.GetTvInfo(results.Results[i].ID, nil)
		if err != nil {
			goutils.LogError(err)
			continue
		}

		title := strings.ReplaceAll(tvShow.Name, "&", "and")
		title = specialCharsRegex.ReplaceAllString(title, "")
		if tvShow.FirstAirDate != "" {
			title += " (" + tvShow.FirstAirDate[0:4] + ")"
		}
		searchedList = append(searchedList, &core.VideoWebResult{
			Title:       title,
			Description: tvShow.Overview,
			PosterURL:   fmt.Sprintf(vws.posterPattern, tvShow.PosterPath),
			Cast:        vws.generateTvShowCastNames(tvShow.Credits.Cast),
		})
	}

	if len(searchedList) < 1 {
		return nil, false
	}

	return searchedList, true
}

func (vws videoWebSearcher) generateMovieCastNames(cast []struct {
	CastID      int `json:"cast_id"`
	Character   string
	CreditID    string `json:"credit_id"`
	ID          int
	Name        string
	Order       int
	ProfilePath string `json:"profile_path"`
}) []string {
	list := make([]string, len(cast))
	for _, c := range cast {
		list = append(list, c.Name)
	}
	return list
}

func (vws videoWebSearcher) generateTvShowCastNames(cast []struct {
	Character   string
	CreditID    string `json:"credit_id"`
	ID          int
	Name        string
	Order       int
	ProfilePath string `json:"profile_path"`
}) []string {
	list := make([]string, len(cast))
	for _, c := range cast {
		list = append(list, c.Name)
	}
	return list
}
