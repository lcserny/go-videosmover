package actions

import (
	"encoding/json"
	"fmt"
	. "github.com/lcserny/go-videosmover/src/shared"
	. "github.com/lcserny/goutils"
	"github.com/pkg/errors"
	"regexp"
	"strings"
)

const NAME_TRIM_REGX_FILE = "name_trim_regx"

type normalizeFunc func(name string) (string, int)
type searchTMDBFunc func(normalizedName string, year int) ([]string, bool)

type OutputConfig struct {
	Normalize  normalizeFunc
	SearchTMDB searchTMDBFunc
}

var (
	configsMap = map[string]*OutputConfig{
		"movie": {movieNormalize, movieTMDBSearch},
		"tv":    {tvSeriesNormalize, tvSeriesTMDBSearch},
	}
	nameTrimPartsRegxs []*regexp.Regexp
)

func init() {
	nameTrimPartsContent, err := configFolder.FindString(NAME_TRIM_REGX_FILE)
	LogError(err)
	nameTrimPartsRegxs = getRegexList(GetLinesFromString(nameTrimPartsContent))
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

	if tmdbApiKey != "" {
		// TODO: get a list of 10? API resolved output names from that decent name
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
	return "", 0
}

func tvSeriesNormalize(name string) (string, int) {
	return "", 0
}

func findOnDisk(normalizedWithYear, diskPath string) ([]string, bool) {
	return nil, false
}

func movieTMDBSearch(normalizedName string, year int) ([]string, bool) {
	return nil, false
}

func tvSeriesTMDBSearch(normalizedName string, year int) ([]string, bool) {
	return nil, false
}
