package actions

import (
	"encoding/json"
	"fmt"
	. "github.com/lcserny/go-videosmover/pkg/models"
	. "github.com/lcserny/goutils"
	"github.com/pkg/errors"
	"github.com/ryanbradynd05/go-tmdb"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type searchTMDBFunc func(normalizedName string, year int, tmdbAPI *tmdb.TMDb, maxTMDBResultCount int) ([]string, bool)

type diskResult struct {
	name        string
	coefficient int
}

var (
	specialCharsRegex = regexp.MustCompile(`([\[._\]])`)
	yearPatternRegex  = regexp.MustCompile(`\s\d{4}$`)
	tmdbFuncMap       = map[string]searchTMDBFunc{
		MOVIE: movieTMDBSearch,
		TV:    tvSeriesTMDBSearch,
	}
)

func OutputAction(jsonPayload []byte, config *ActionConfig) (string, error) {
	var request OutputRequestData
	err := json.Unmarshal(jsonPayload, &request)
	LogError(err)
	if err != nil {
		return "", err
	}

	normalized, year := normalize(request.Name, config.compiledNameTrimRegexes)
	normalizedWithYear := appendYear(normalized, year)
	if onDisk, found := findOnDisk(normalizedWithYear, request.DiskPath, config.MaxOutputWalkDepth, config.SimilarityPercent); found {
		return getJSONEncodedString(OutputResponseData{onDisk, ORIGIN_DISK}), nil
	}

	if config.tmdbAPI != nil {
		tmdbFunc, err := getTMDBFunc(request.Type)
		if err != nil {
			LogError(err)
			return "", err
		}

		// TODO: implement local persisted simple file caching (benchmark IO, goroutine load cache?)
		// get from cache if available (if cacheFile exists)
		// generate cache key (using normalized, year, type)
		// open cache file
		// go through line by line, if starts with cacheKey, get tmdbNames

		if tmdbNames, found := tmdbFunc(normalized, year, config.tmdbAPI, config.MaxTMDBResultCount); found {
			// save in cache (create cacheFile if necessary)
			// generate cache key (using normalized, year, type)
			// save in cache file to line 0 always as <cacheKey>###<commaSeparatedNames>
			// if nr. of lines is > X, trim those lines from file

			return getJSONEncodedString(OutputResponseData{tmdbNames, ORIGIN_TMDB}), nil
		}
	}

	return getJSONEncodedString(OutputResponseData{[]string{normalizedWithYear}, ORIGIN_NAME}), nil
}

func getRegexList(patterns []string) (regxs []*regexp.Regexp) {
	for _, pat := range patterns {
		regxs = append(regxs, regexp.MustCompile(fmt.Sprintf("(?i)(-?%s)", pat)))
	}
	return regxs
}

func appendYear(normalizedName string, year int) string {
	if year > 0 {
		return fmt.Sprintf("%s (%d)", normalizedName, year)
	}
	return normalizedName
}

func getTMDBFunc(itemType string) (searchTMDBFunc, error) {
	if tmdbFunc, found := tmdbFuncMap[strings.ToLower(itemType)]; found {
		return tmdbFunc, nil
	}
	return nil, errors.New("No TMDB function found for type:" + itemType)
}

func normalize(name string, nameTrimPartsRegxs []*regexp.Regexp) (string, int) {
	// trim
	for _, pat := range nameTrimPartsRegxs {
		if loc := pat.FindStringIndex(name); loc != nil {
			name = name[0:loc[0]]
		}
	}

	// strip special chars
	name = specialCharsRegex.ReplaceAllString(name, " ")
	name = strings.Trim(name, " ")

	// title case
	name = strings.Title(name)

	// resolve year
	yearLoc := yearPatternRegex.FindStringIndex(name)
	if yearLoc != nil {
		year, err := strconv.ParseInt(name[yearLoc[0]+1:], 0, 32)
		LogError(err)
		return name[0:yearLoc[0]], int(year)
	}

	return name, 0
}

func findOnDisk(normalizedWithYear, diskPath string, maxOutputWalkDepth, acceptedSimPercent int) (results []string, found bool) {
	var tmpList []diskResult
	err := filepath.Walk(diskPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			LogError(err)
			return nil
		}

		if info.IsDir() && diskPath != path && walkDepthIsAcceptable(diskPath, path, maxOutputWalkDepth) {
			distance := LevenshteinDistance(info.Name(), normalizedWithYear)
			bigger := MaxInt(len(normalizedWithYear), len(info.Name()))
			simPercent := int(float32(bigger-distance) / float32(bigger) * 100)
			if simPercent > acceptedSimPercent {
				tmpList = append(tmpList, diskResult{info.Name(), simPercent})
			}
		}

		return nil
	})
	LogError(err)
	if err != nil {
		return results, false
	}

	if len(tmpList) < 1 {
		return results, false
	}

	// sort by highest coefficient
	sort.Slice(tmpList, func(i, j int) bool {
		return tmpList[i].coefficient > tmpList[j].coefficient
	})
	for _, ele := range tmpList {
		results = append(results, ele.name)
	}

	return results, true
}

func movieTMDBSearch(normalizedName string, year int, tmdbAPI *tmdb.TMDb, maxTMDBResultCount int) (searchedList []string, found bool) {
	options := map[string]string{"page": "1", "language": "en"}
	if year > 0 {
		options["year"] = string(year)
	}

	results, err := tmdbAPI.SearchMovie(normalizedName, options)
	if err != nil {
		LogError(err)
		return searchedList, false
	}

	if len(results.Results) < 1 {
		return searchedList, false
	}

	for i := 0; i < MinInt(maxTMDBResultCount, len(results.Results)); i++ {
		movie := results.Results[i]
		searchedList = append(searchedList, fmt.Sprintf("%s (%s)", movie.Title, movie.ReleaseDate))
	}

	return searchedList, true
}

func tvSeriesTMDBSearch(normalizedName string, year int, tmdbAPI *tmdb.TMDb, maxTMDBResultCount int) (searchedList []string, found bool) {
	options := map[string]string{"page": "1", "language": "en"}
	if year > 0 {
		options["first_air_date_year"] = string(year)
	}

	results, err := tmdbAPI.SearchTv(normalizedName, options)
	if err != nil {
		LogError(err)
		return searchedList, false
	}

	if len(results.Results) < 1 {
		return searchedList, false
	}

	for i := 0; i < MinInt(maxTMDBResultCount, len(results.Results)); i++ {
		tvShow := results.Results[i]
		searchedList = append(searchedList, fmt.Sprintf("%s (%s)", tvShow.Name, tvShow.FirstAirDate[0:4]))
	}

	return searchedList, true
}
