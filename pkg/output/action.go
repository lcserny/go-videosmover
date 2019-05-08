package output

import (
	"bufio"
	"fmt"
	"github.com/lcserny/goutils"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"videosmover/pkg"
	"videosmover/pkg/action"
)

// TODO: improve this

func NewAction(cfg *core.ActionConfig, c core.Codec, ws core.VideoWebSearcher) action.Action {
	oa := outputAction{config: cfg, codec: c, webSearcher: ws}
	if cfg.NameTrimRegexes != nil {
		for _, pat := range cfg.NameTrimRegexes {
			oa.namePatterns = append(oa.namePatterns, regexp.MustCompile(fmt.Sprintf("(?i)(-?%s)", pat)))
		}
	}
	return &oa
}

type outputAction struct {
	config       *core.ActionConfig
	codec        core.Codec
	webSearcher  core.VideoWebSearcher
	namePatterns []*regexp.Regexp
}

type videoSearchFunc func(normalizedName string, year int, webSearcher core.VideoWebSearcher, maxAccepted int) ([]string, bool)

type diskResult struct {
	name        string
	coefficient int
}

var (
	outputTMDBCacheFile               = goutils.GetAbsCurrentPathOf("tmdbOutput.cache")
	outputTMDBCacheSeparator          = "###"
	outputTMDBCacheFileNamesSeparator = ";"
	specialCharsReplaceMap            = map[string]string{"&": "and"}
	specialCharsRegex                 = regexp.MustCompile(`[^a-zA-Z0-9-\s]`)
	spaceMergeRegex                   = regexp.MustCompile(`\s{2,}`)
	yearRegex                         = regexp.MustCompile(`\s\d{4}$`)
	releaseDateRegex                  = regexp.MustCompile(`\s+\(\d{4}(-\d{2}-\d{2})?\)$`)
	tmdbFuncMap                       = map[string]videoSearchFunc{
		action.MOVIE: movieTMDBSearch,
		action.TV:    tvSeriesTMDBSearch,
	}
)

func (oa outputAction) Execute(jsonPayload []byte) (string, error) {
	var request RequestData
	if err := oa.codec.Decode(jsonPayload, &request); err != nil {
		goutils.LogError(err)
		return "", err
	}

	normalized, year := normalize(request.Name, oa.namePatterns)
	normalizedWithYear := appendYear(normalized, year)
	if onDisk, found := findOnDisk(normalized, request.DiskPath, oa.config.MaxOutputWalkDepth, oa.config.SimilarityPercent); found {
		return oa.codec.EncodeString(ResponseData{onDisk, ORIGIN_DISK})
	}

	if !request.SkipOnlineSearch && oa.webSearcher != nil {
		tmdbFunc, err := getTMDBFunc(request.Type)
		if err != nil {
			goutils.LogError(err)
			return "", err
		}

		cacheKey := generateTMDBOutputCacheKey(request.Type, normalizedWithYear, outputTMDBCacheSeparator)
		if !request.SkipCache {
			if _, err := os.Stat(outputTMDBCacheFile); !os.IsNotExist(err) {
				if cachedTMDBNames, exist := getFromTMDBOutputCache(cacheKey, outputTMDBCacheFile); exist {
					return oa.codec.EncodeString(ResponseData{cachedTMDBNames, ORIGIN_TMDB_CACHE})
				}
			}
		}

		if tmdbNames, found := tmdbFunc(normalized, year, oa.webSearcher, oa.config.MaxWebSearchResultCount); found {
			saveInTMDBOutputCache(cacheKey, tmdbNames, outputTMDBCacheFile, oa.config.OutWebSearchCacheLimit)
			return oa.codec.EncodeString(ResponseData{tmdbNames, ORIGIN_TMDB})
		}
	}

	return oa.codec.EncodeString(ResponseData{[]string{normalizedWithYear}, ORIGIN_NAME})
}

func saveInTMDBOutputCache(cacheKey string, tmdbNames []string, cacheFile string, cacheLimit int) {
	oldFile, err := os.OpenFile(cacheFile, os.O_RDONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		goutils.LogError(err)
		return
	}

	tmpName := cacheFile + "_tmp"
	newFile, err := os.OpenFile(tmpName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		goutils.LogError(err)
		return
	}

	cleanupAndFail := func(err error) {
		goutils.LogError(err)
		goutils.CloseFile(oldFile)
		goutils.CloseFile(newFile)
		err = os.Remove(tmpName)
		goutils.LogError(err)
	}
	moveAndSucceed := func() {
		goutils.CloseFile(oldFile)
		goutils.CloseFile(newFile)
		err = os.Rename(tmpName, cacheFile)
		goutils.LogError(err)
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
		goutils.LogError(err)
		return nil, false
	}
	defer goutils.CloseFile(openFile)

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

func getTMDBFunc(itemType string) (videoSearchFunc, error) {
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
		goutils.LogError(err)
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
		goutils.LogWarning(fmt.Sprintf("Diskpath provided not found: %s", diskPath))
		return results, false
	}

	var tmpList []diskResult
	err := filepath.Walk(diskPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			goutils.LogError(err)
			return nil
		}

		if info.IsDir() && diskPath != path && action.WalkDepthIsAcceptable(diskPath, path, maxOutputWalkDepth) {
			nameWithoutDate := trimReleaseDate(info.Name())
			distance := goutils.LevenshteinDistance(nameWithoutDate, normalized)
			bigger := goutils.MaxInt(len(normalized), len(nameWithoutDate))
			simPercent := int(float32(bigger-distance) / float32(bigger) * 100)
			if simPercent > acceptedSimPercent {
				tmpList = append(tmpList, diskResult{info.Name(), simPercent})
			}
		}

		return nil
	})
	if err != nil {
		goutils.LogError(err)
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

func movieTMDBSearch(normalizedName string, year int, webSearcher core.VideoWebSearcher, maxAccepted int) (searchedList []string, found bool) {
	results, err := webSearcher.SearchMovies(normalizedName, year)
	if err != nil {
		goutils.LogError(err)
		return searchedList, false
	}

	if len(results) < 1 {
		return searchedList, false
	}

	for i := 0; i < goutils.MinInt(maxAccepted, len(results)); i++ {
		video := results[i]
		outName := replaceSpecialChars(video.Title)
		outName = specialCharsRegex.ReplaceAllString(outName, "")
		if video.ReleaseDate != "" {
			outName += " (" + video.ReleaseDate + ")"
		}
		searchedList = append(searchedList, outName)
	}

	return searchedList, true
}

func tvSeriesTMDBSearch(normalizedName string, year int, webSearcher core.VideoWebSearcher, maxAccepted int) (searchedList []string, found bool) {
	results, err := webSearcher.SearchTVSeries(normalizedName, year)
	if err != nil {
		goutils.LogError(err)
		return searchedList, false
	}

	if len(results) < 1 {
		return searchedList, false
	}

	for i := 0; i < goutils.MinInt(maxAccepted, len(results)); i++ {
		video := results[i]
		outName := replaceSpecialChars(video.Title)
		outName = specialCharsRegex.ReplaceAllString(outName, "")
		if video.ReleaseDate != "" {
			outName += " (" + video.ReleaseDate[0:4] + ")"
		}
		searchedList = append(searchedList, outName)
	}

	return searchedList, true
}
