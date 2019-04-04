package action

import (
	"encoding/json"
	"fmt"
	"github.com/lcserny/goutils"
	"github.com/ryanbradynd05/go-tmdb"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

type Config struct {
	MinimumVideoSize          int64    `json:"minimumVideoSize"`
	SimilarityPercent         int      `json:"similarityPercent"`
	MaxOutputWalkDepth        int      `json:"maxOutputWalkDepth"`
	MaxSearchWalkDepth        int      `json:"maxSearchWalkDepth"`
	MaxTMDBResultCount        int      `json:"maxTMDBResultCount"`
	OutTMDBCacheLimit         int      `json:"outTMDBCacheLimit"`
	HeaderBytesSize           int      `json:"headerBytesSize"`
	RestrictedRemovePaths     []string `json:"restrictedRemovePaths"`
	NameTrimRegexes           []string `json:"nameTrimRegexes"`
	SearchExcludePaths        []string `json:"searchExcludePaths"`
	AllowedMIMETypes          []string `json:"allowedMIMETypes"`
	AllowedSubtitleExtensions []string `json:"allowedSubtitleExtensions"`

	TmdbAPI                 *tmdb.TMDb
	CompiledNameTrimRegexes []*regexp.Regexp
}

func NewConfig(path, file string) *Config {
	content, err := ioutil.ReadFile(filepath.Join(path, file))
	goutils.LogFatal(err)

	var ac Config
	err = json.Unmarshal(content, &ac)
	goutils.LogFatal(err)

	if key, exists := os.LookupEnv("TMDB_API_KEY"); exists {
		ac.TmdbAPI = tmdb.Init(tmdb.Config{key, false, nil})
	}

	if ac.NameTrimRegexes != nil {
		ac.CompiledNameTrimRegexes = getRegexList(ac.NameTrimRegexes)
	}

	return &ac
}

func getRegexList(patterns []string) (regxs []*regexp.Regexp) {
	for _, pat := range patterns {
		regxs = append(regxs, regexp.MustCompile(fmt.Sprintf("(?i)(-?%s)", pat)))
	}
	return regxs
}
