package actions

import (
	"encoding/json"
	"fmt"
	. "github.com/lcserny/go-videosmover/src/shared"
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

const (
	MAX_OUTPUT_WALK_DEPTH = 2
	MAX_TMDB_RES_COUNT    = 10
	NAME_TRIM_REGX_FILE   = "name_trim_regx"
	SIM_PERCENT_KEY       = "similarity.percent"
)

type normalizeFunc func(name string) (string, int)
type searchTMDBFunc func(normalizedName string, year int) ([]string, bool)

type OutputConfig struct {
	Normalize  normalizeFunc
	SearchTMDB searchTMDBFunc
}

type diskResult struct {
	name        string
	coefficient int
}

var (
	specialCharsRegex = regexp.MustCompile(`([\[._\]])`)
	yearPatternRegex  = regexp.MustCompile(`\s\d{4}$`)
	configsMap        = map[string]*OutputConfig{
		"movie": {movieNormalize, movieTMDBSearch},
		"tv":    {tvSeriesNormalize, tvSeriesTMDBSearch},
	}
	nameTrimPartsRegxs []*regexp.Regexp
	acceptedSimPercent int
	tmdbAPI            *tmdb.TMDb
)

func init() {
	nameTrimPartsContent, err := configFolder.FindString(NAME_TRIM_REGX_FILE)
	LogError(err)
	nameTrimPartsRegxs = getRegexList(GetLinesFromString(nameTrimPartsContent))

	if appProperties.HasProperty(SIM_PERCENT_KEY) {
		acceptedSimPercent = appProperties.GetPropertyAsInt(SIM_PERCENT_KEY)
	}

	if tmdbApiKey != "" {
		tmdbAPI = tmdb.Init(tmdb.Config{tmdbApiKey, false, nil})
	}
}

func OutputAction(jsonPayload []byte) (string, error) {
	var request OutputRequestData
	err := json.Unmarshal(jsonPayload, &request)
	LogError(err)
	if err != nil {
		return "", err
	}

	config, err := getConfig(request.Type)
	if err != nil {
		LogError(err)
		return "", err
	}

	normalized, year := config.Normalize(request.Name)
	normalizedWithYear := normalizeWithYear(normalized, year)
	if onDisk, found := findOnDisk(normalizedWithYear, request.DiskPath); found {
		return getJSONEncodedString(onDisk), nil
	}

	if tmdbAPI != nil {
		if tmdbNames, found := config.SearchTMDB(normalized, year); found {
			return getJSONEncodedString(tmdbNames), nil
		}
	}

	return getJSONEncodedString([]string{normalizedWithYear}), nil
}

func getRegexList(patterns []string) (regxs []*regexp.Regexp) {
	for _, pat := range patterns {
		regxs = append(regxs, regexp.MustCompile(fmt.Sprintf("(?i)(-?%s)", pat)))
	}
	return regxs
}

func normalizeWithYear(normalizedName string, year int) string {
	if year > 0 {
		return fmt.Sprintf("%s (%d)", normalizedName, year)
	}
	return normalizedName
}

func getConfig(itemType string) (*OutputConfig, error) {
	if config, found := configsMap[strings.ToLower(itemType)]; found {
		return config, nil
	}
	return nil, errors.New("No OutputConfig found for type:" + itemType)
}

func movieNormalize(name string) (string, int) {
	normalized, _ := tvSeriesNormalize(name)

	// resolve year
	yearLoc := yearPatternRegex.FindStringIndex(normalized)
	if yearLoc != nil {
		year, err := strconv.ParseInt(normalized[yearLoc[0]+1:], 0, 32)
		LogError(err)
		return normalized[0:yearLoc[0]], int(year)
	}

	return normalized, 0
}

func tvSeriesNormalize(name string) (string, int) {
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

	return name, 0
}

func findOnDisk(normalizedWithYear, diskPath string) (results []string, found bool) {
	var tmpList []diskResult
	err := filepath.Walk(diskPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			LogError(err)
			return nil
		}

		if info.IsDir() && diskPath != path && walkDepthIsAcceptable(diskPath, path, MAX_OUTPUT_WALK_DEPTH) {
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

func movieTMDBSearch(normalizedName string, year int) (searchedList []string, found bool) {
	results, err := tmdbAPI.SearchMovie(normalizedName, map[string]string{"year": string(year), "page": "1", "language": "en"})
	if err != nil {
		LogError(err)
		return searchedList, false
	}

	for i := 0; i < MAX_TMDB_RES_COUNT; i++ {
		movie := results.Results[i]
		searchedList = append(searchedList, fmt.Sprintf("%s (%s)", movie.Title, movie.ReleaseDate))
	}

	return searchedList, true
}

func tvSeriesTMDBSearch(normalizedName string, year int) (searchedList []string, found bool) {
	results, err := tmdbAPI.SearchTv(normalizedName, map[string]string{"first_air_date_year": string(year), "page": "1", "language": "en"})
	if err != nil {
		LogError(err)
		return searchedList, false
	}

	for i := 0; i < MAX_TMDB_RES_COUNT; i++ {
		tvShow := results.Results[i]
		searchedList = append(searchedList, fmt.Sprintf("%s (%s)", tvShow.Name, tvShow.FirstAirDate[0:4]))
	}

	return searchedList, true
}
