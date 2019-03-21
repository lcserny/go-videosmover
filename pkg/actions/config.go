package actions

import (
	"encoding/json"
	"github.com/lcserny/go-videosmover/pkg/models"
	. "github.com/lcserny/goutils"
	"github.com/pkg/errors"
	"github.com/ryanbradynd05/go-tmdb"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var actionsMap = map[string]Action{
	"SEARCH": SearchAction,
	"OUTPUT": OutputAction,
	"MOVE":   MoveAction,
	"DELETE": DeleteAction,
}

type Action func(jsonPayload []byte, config *models.ActionConfig) (string, error)

func GenerateActionConfig(configsPath, actionConfigFile string) *models.ActionConfig {
	configBytes, err := ioutil.ReadFile(filepath.Join(configsPath, actionConfigFile))
	LogFatal(err)

	var actionConfig models.ActionConfig
	err = json.Unmarshal(configBytes, &actionConfig)
	LogFatal(err)

	if key, exists := os.LookupEnv("TMDB_API_KEY"); exists {
		actionConfig.TmdbAPI = tmdb.Init(tmdb.Config{key, false, nil})
	}

	if actionConfig.NameTrimRegexes != nil {
		actionConfig.CompiledNameTrimRegexes = getRegexList(actionConfig.NameTrimRegexes)
	}

	return &actionConfig
}

func UnknownAction(jsonPayload []byte, config *models.ActionConfig) (string, error) {
	return "", errors.New("Unknown action given")
}

func NewActionFrom(val string) Action {
	if action, ok := actionsMap[strings.ToUpper(val)]; ok {
		return action
	}
	return UnknownAction
}
