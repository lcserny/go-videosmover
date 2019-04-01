package output

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/lcserny/go-videosmover/pkg/action"
	"github.com/lcserny/go-videosmover/pkg/convert"
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
	ORIGIN_NAME       = "NAME"
	ORIGIN_DISK       = "DISK"
	ORIGIN_TMDB       = "TMDB"
	ORIGIN_TMDB_CACHE = "TMDB_CACHE"
)

type searchTMDBFunc func(normalizedName string, year int, tmdbAPI *tmdb.TMDb, maxTMDBResultCount int) ([]string, bool)

type diskResult struct {
	name        string
	coefficient int
}

var (
	outputTMDBCacheFile               = GetAbsCurrentPathOf("tmdbOutput.cache")
	outputTMDBCacheSeparator          = "###"
	outputTMDBCacheFileNamesSeparator = ";"
	specialCharsReplaceMap            = map[string]string{"&": "and"}
	specialCharsRegex                 = regexp.MustCompile(`[^a-zA-Z0-9-\s]`)
	spaceMergeRegex                   = regexp.MustCompile(`\s{2,}`)
	yearRegex                         = regexp.MustCompile(`\s\d{4}$`)
	releaseDateRegex                  = regexp.MustCompile(`\s+\(\d{4}(-\d{2}-\d{2})?\)$`)
	tmdbFuncMap                       = map[string]searchTMDBFunc{
		action.MOVIE: movieTMDBSearch,
		action.TV:    tvSeriesTMDBSearch,
	}
)

func Action(jsonPayload []byte, config *convert.ActionConfig) (string, error) {
	var request convert.OutputRequestData
	err := json.Unmarshal(jsonPayload, &request)
	if err != nil {
		LogError(err)
		return "", err
	}

	normalized, year := normalize(request.Name, config.CompiledNameTrimRegexes)
	normalizedWithYear := appendYear(normalized, year)
	if onDisk, found := findOnDisk(normalized, request.DiskPath, config.MaxOutputWalkDepth, config.SimilarityPercent); found {
		return convert.GetJSONEncodedString(convert.OutputResponseData{onDisk, ORIGIN_DISK}), nil
	}

	if !request.SkipOnlineSearch && config.TmdbAPI != nil {
		tmdbFunc, err := getTMDBFunc(request.Type)
		if err != nil {
			LogError(err)
			return "", err
		}

		cacheKey := generateTMDBOutputCacheKey(request.Type, normalizedWithYear, outputTMDBCacheSeparator)
		if !request.SkipCache {
			if _, err := os.Stat(outputTMDBCacheFile); !os.IsNotExist(err) {
				if cachedTMDBNames, exist := getFromTMDBOutputCache(cacheKey, outputTMDBCacheFile); exist {
					return convert.GetJSONEncodedString(convert.OutputResponseData{cachedTMDBNames, ORIGIN_TMDB_CACHE}), nil
				}
			}
		}

		if tmdbNames, found := tmdbFunc(normalized, year, config.TmdbAPI, config.MaxTMDBResultCount); found {
			saveInTMDBOutputCache(cacheKey, tmdbNames, outputTMDBCacheFile, config.OutTMDBCacheLimit)
			return convert.GetJSONEncodedString(convert.OutputResponseData{tmdbNames, ORIGIN_TMDB}), nil
		}
	}

	return convert.GetJSONEncodedString(convert.OutputResponseData{[]string{normalizedWithYear}, ORIGIN_NAME}), nil
}

func saveInTMDBOutputCache(cacheKey string, tmdbNames []string, cacheFile string, cacheLimit int) {
	oldFile, err := os.OpenFile(cacheFile, os.O_RDONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		LogError(err)
		return
	}

	tmpName := cacheFile + "_tmp"
	newFile, err := os.OpenFile(tmpName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		LogError(err)
		return
	}

	cleanupAndFail := func(err error) {
		LogError(err)
		CloseFile(oldFile)
		CloseFile(newFile)
		err = os.Remove(tmpName)
		LogError(err)
	}
	moveAndSucceed := func() {
		CloseFile(oldFile)
		CloseFile(newFile)
		err = os.Rename(tmpName, cacheFile)
		LogError(err)
	}

	firstLine := cacheKey + strings.Join(tmdbNames, outputTMDBCacheFileNamesSeparator)
	if _, err = fmt.Fprintln(newFile, firstLine); err != nil {
		cleanupAndFail(err)
		return
	}

	lineCounter := 1
	scanner := bufio.NewScanner(oldFile)
	for scanner.Scan() {
		if lineCounter >= cacheLimit {
			break
		}

		if _, err = fmt.Fprintln(newFile, scanner.Text()); err != nil {
			cleanupAndFail(err)
			return
		}

		lineCounter++
	}

	moveAndSucceed()
}

func getFromTMDBOutputCache(cacheKey, cacheFile string) ([]string, bool) {
	openFile, err := os.OpenFile(cacheFile, os.O_RDONLY, os.ModePerm)
	if err != nil {
		LogError(err)
		return nil, false
	}
	defer CloseFile(openFile)

	scanner := bufio.NewScanner(openFile)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, cacheKey) {
			return strings.Split(line[len(cacheKey):], outputTMDBCacheFileNamesSeparator), true
		}
	}
	return nil, false
}

func generateTMDBOutputCacheKey(videoType, normalizedWithYear, separator string) string {
	return fmt.Sprintf("%s__%s%s", strings.ToUpper(videoType), normalizedWithYear, separator)
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
	name = replaceSpecialChars(name)
	name = specialCharsRegex.ReplaceAllString(name, " ")
	name = spaceMergeRegex.ReplaceAllString(name, " ")
	name = strings.Trim(name, " ")

	// title case
	name = strings.Title(name)

	// resolve year
	yearLoc := yearRegex.FindStringIndex(name)
	if yearLoc != nil {
		year, err := strconv.ParseInt(name[yearLoc[0]+1:], 0, 32)
		LogError(err)
		return name[0:yearLoc[0]], int(year)
	}

	return name, 0
}

func replaceSpecialChars(text string) string {
	for k, v := range specialCharsReplaceMap {
		text = strings.ReplaceAll(text, k, v)
	}
	return text
}

func findOnDisk(normalized, diskPath string, maxOutputWalkDepth, acceptedSimPercent int) (results []string, found bool) {
	if _, err := os.Stat(diskPath); os.IsNotExist(err) {
		LogWarning(fmt.Sprintf("Diskpath provided not found: %s", diskPath))
		return results, false
	}

	var tmpList []diskResult
	err := filepath.Walk(diskPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			LogError(err)
			return nil
		}

		if info.IsDir() && diskPath != path && action.WalkDepthIsAcceptable(diskPath, path, maxOutputWalkDepth) {
			nameWithoutDate := trimReleaseDate(info.Name())
			distance := LevenshteinDistance(nameWithoutDate, normalized)
			bigger := MaxInt(len(normalized), len(nameWithoutDate))
			simPercent := int(float32(bigger-distance) / float32(bigger) * 100)
			if simPercent > acceptedSimPercent {
				tmpList = append(tmpList, diskResult{info.Name(), simPercent})
			}
		}

		return nil
	})
	if err != nil {
		LogError(err)
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

func trimReleaseDate(nameWithReleaseDate string) string {
	return releaseDateRegex.ReplaceAllString(nameWithReleaseDate, "")
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
		outName := replaceSpecialChars(movie.Title)
		outName = specialCharsRegex.ReplaceAllString(outName, "")
		if movie.ReleaseDate != "" {
			outName += " (" + movie.ReleaseDate + ")"
		}
		searchedList = append(searchedList, outName)
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
		outName := replaceSpecialChars(tvShow.Name)
		outName = specialCharsRegex.ReplaceAllString(outName, "")
		if tvShow.FirstAirDate != "" {
			outName += " (" + tvShow.FirstAirDate[0:4] + ")"
		}
		searchedList = append(searchedList, outName)
	}

	return searchedList, true
}
