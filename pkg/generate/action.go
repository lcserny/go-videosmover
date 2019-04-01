package generate

import (
	"encoding/json"
	"github.com/lcserny/go-videosmover/pkg/action"
	"github.com/lcserny/go-videosmover/pkg/convert"
	"github.com/lcserny/goutils"
	"github.com/ryanbradynd05/go-tmdb"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func NewActionFrom(val string) action.Action {
	if a, ok := action.ActionsMap[strings.ToUpper(val)]; ok {
		return a
	}
	return action.UnknownAction
}

func NewActionConfig(configsPath, actionConfigFile string) *convert.ActionConfig {
	configBytes, err := ioutil.ReadFile(filepath.Join(configsPath, actionConfigFile))
	goutils.LogFatal(err)

	var actionConfig convert.ActionConfig
	err = json.Unmarshal(configBytes, &actionConfig)
	goutils.LogFatal(err)

	if key, exists := os.LookupEnv("TMDB_API_KEY"); exists {
		actionConfig.TmdbAPI = tmdb.Init(tmdb.Config{key, false, nil})
	}

	if actionConfig.NameTrimRegexes != nil {
		actionConfig.CompiledNameTrimRegexes = action.GetRegexList(actionConfig.NameTrimRegexes)
	}

	return &actionConfig
}

func NewTestActionConfig() *convert.ActionConfig {
	return NewActionConfig("../../cfg/commander", "actions.test.json")
}
