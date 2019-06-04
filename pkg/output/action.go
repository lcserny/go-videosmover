package output

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/lcserny/goutils"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"videosmover/pkg"
	"videosmover/pkg/action"
)

func NewAction(cfg *core.ActionConfig, c core.Codec, ws core.VideoWebSearcher) core.Action {
	oa := outputAction{
		config:                 cfg,
		codec:                  c,
		webSearcher:            ws,
		preNormalizedNameRegex: regexp.MustCompile(`^\s*(?P<name>[a-zA-Z0-9-\s]+)\s\((?P<year>\d{4})(-\d{1,2}-\d{1,2})?\)$`),
		specialCharsRegex:      regexp.MustCompile(`[^a-zA-Z0-9-\s]`),
		spaceMergeRegex:        regexp.MustCompile(`\s{2,}`),
		yearRegex:              regexp.MustCompile(`\s\d{4}$`),
		releaseDateRegex:       regexp.MustCompile(`\s+\(\d{4}(-\d{2}-\d{2})?\)$`),
	}
	if cfg.NameTrimRegexes != nil {
		for _, pat := range cfg.NameTrimRegexes {
			oa.namePatterns = append(oa.namePatterns, regexp.MustCompile(fmt.Sprintf("(?i)(-?%s)", pat)))
		}
	}
	return &oa
}

type outputAction struct {
	config                 *core.ActionConfig
	codec                  core.Codec
	webSearcher            core.VideoWebSearcher
	namePatterns           []*regexp.Regexp
	preNormalizedNameRegex *regexp.Regexp
	specialCharsRegex      *regexp.Regexp
	spaceMergeRegex        *regexp.Regexp
	yearRegex              *regexp.Regexp
	releaseDateRegex       *regexp.Regexp
}

type videoSearchFunc func(name string, year, maxResCount int, specialCharsRegex *regexp.Regexp) ([]*core.VideoWebResult, bool)

type diskResult struct {
	name        string
	coefficient int
}

func (oa outputAction) Execute(jsonPayload []byte) (string, error) {
	var request RequestData
	if err := oa.codec.Decode(jsonPayload, &request); err != nil {
		goutils.LogError(err)
		return "", err
	}

	normalized, year := oa.normalize(request.Name)
	normalizedWithYear := oa.appendYear(normalized, year)
	if onDisk, found := oa.findOnDisk(normalized, request.DiskPath); found {
		return oa.codec.EncodeString(ResponseData{oa.generateVideoWebResultsFromStrings(onDisk), ORIGIN_DISK})
	}

	if !request.SkipOnlineSearch && oa.webSearcher.CanSearch() {
		webSearcherFunc, err := oa.getWebSearcherFunc(request.Type)
		if err != nil {
			goutils.LogError(err)
			return "", err
		}

		cacheKey := oa.generateOutputCacheKey(request.Type, normalizedWithYear)
		if !request.SkipCache {
			// TODO: does this work
			reqData, err := oa.codec.EncodeBytes(map[string]string{"key": cacheKey})
			if err != nil {
				goutils.LogError(err)
				return "", err
			}
			resp, err := http.Post(fmt.Sprintf("%s/get", oa.config.CacheApiURL), "application/json", bytes.NewBuffer(reqData))
			if err != nil {
				goutils.LogError(err)
				return "", err
			}
			defer resp.Body.Close()

			var cachedResults []*core.VideoWebResult
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				goutils.LogError(err)
				return "", err
			}

			if len(b) > 0 {
				err = oa.codec.Decode(b, cachedResults)
				if err != nil {
					goutils.LogError(err)
					return "", err
				}
				return oa.codec.EncodeString(ResponseData{cachedResults, ORIGIN_TMDB_CACHE})
			}
		}

		if webResults, found := webSearcherFunc(normalized, year, oa.config.MaxWebSearchResultCount, oa.specialCharsRegex); found {
			// TODO: does this work
			reqData, err := oa.codec.EncodeBytes(map[string]interface{}{"key": cacheKey, "val": webResults})
			if err != nil {
				goutils.LogError(err)
				return "", err
			}
			resp, err := http.Post(fmt.Sprintf("%s/set", oa.config.CacheApiURL), "application/json", bytes.NewBuffer(reqData))
			if err != nil {
				goutils.LogError(err)
				return "", err
			}
			defer resp.Body.Close()
			return oa.codec.EncodeString(ResponseData{webResults, ORIGIN_TMDB})
		}
	}

	return oa.codec.EncodeString(ResponseData{oa.generateVideoWebResultsFromStrings([]string{normalizedWithYear}), ORIGIN_NAME})
}

func (oa outputAction) generateVideoWebResultsFromStrings(strs []string) []*core.VideoWebResult {
	res := make([]*core.VideoWebResult, 0)
	for _, s := range strs {
		res = append(res, &core.VideoWebResult{Title: s})
	}
	return res
}

func (oa outputAction) getWebSearcherFunc(itemType string) (videoSearchFunc, error) {
	if itemType == action.MOVIE {
		return oa.webSearcher.SearchMovies, nil
	} else if itemType == action.TV {
		return oa.webSearcher.SearchTVSeries, nil
	}
	return nil, errors.New("No TMDB function found for type:" + itemType)
}

func (oa outputAction) generateOutputCacheKey(videoType, normalizedWithYear string) string {
	return fmt.Sprintf("%s__%s", strings.ToUpper(videoType), normalizedWithYear)
}

func (oa outputAction) appendYear(normalizedName string, year int) string {
	if year > 0 {
		return fmt.Sprintf("%s (%d)", normalizedName, year)
	}
	return normalizedName
}

func (oa outputAction) trimReleaseDate(nameWithReleaseDate string) string {
	return oa.releaseDateRegex.ReplaceAllString(nameWithReleaseDate, "")
}

func (oa outputAction) normalize(name string) (string, int) {
	// handle already normalized text
	if oa.preNormalizedNameRegex.MatchString(name) {
		resMap := goutils.GetRegexSubgroups(oa.preNormalizedNameRegex, name)
		n := strings.Trim(resMap["name"], " ")
		y, _ := strconv.Atoi(resMap["year"])
		return n, y
	}

	// trim
	for _, pat := range oa.namePatterns {
		if loc := pat.FindStringIndex(name); loc != nil {
			name = name[0:loc[0]]
		}
	}

	// strip special chars
	name = strings.ReplaceAll(name, "&", "and")
	name = oa.specialCharsRegex.ReplaceAllString(name, " ")
	name = oa.spaceMergeRegex.ReplaceAllString(name, " ")
	name = strings.Trim(name, " ")

	// title case
	name = strings.Title(name)

	// resolve year
	yearLoc := oa.yearRegex.FindStringIndex(name)
	if yearLoc != nil {
		year, err := strconv.ParseInt(name[yearLoc[0]+1:], 0, 32)
		goutils.LogError(err)
		return name[0:yearLoc[0]], int(year)
	}

	return name, 0
}

func (oa outputAction) findOnDisk(normalized, diskPath string) (results []string, found bool) {
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

		if info.IsDir() && diskPath != path && action.WalkDepthIsAcceptable(diskPath, path, oa.config.MaxOutputWalkDepth) {
			nameWithoutDate := oa.trimReleaseDate(info.Name())
			distance := goutils.LevenshteinDistance(nameWithoutDate, normalized)
			bigger := goutils.MaxInt(len(normalized), len(nameWithoutDate))
			simPercent := int(float32(bigger-distance) / float32(bigger) * 100)
			if simPercent > oa.config.SimilarityPercent {
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
