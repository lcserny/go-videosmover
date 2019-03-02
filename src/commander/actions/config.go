package actions

import (
	"github.com/gobuffalo/packr"
	. "github.com/lcserny/goutils"
	"github.com/pkg/errors"
	"log"
	"os"
	"strings"
)

const (
	BASE_CONF_FILE = "commander.properties"
	LOG_FILE       = "videosmover-commander.log"
)

// if these need to be accessible in other packages, put this in a `shared` package and export vars
var (
	appProperties *ConfigProperties
	configFolder  packr.Box
	// TODO: add action to use from qBittorrent when done downloading to add to a db or something,
	//  then in Android app on startup it can maybe show you last finished downloading torrents
	actionsMap = map[string]Action{
		"SEARCH": SearchAction,
		"OUTPUT": OutputAction,
	}
)

type Action func(jsonPayload []byte) (string, error)

func init() {
	initLogger()

	configFolder = packr.NewBox("../../../config")
	content, err := configFolder.FindString(BASE_CONF_FILE)
	LogFatal(err)

	appProperties = ReadProperties(content)
}

func initLogger() {
	openFile, err := os.OpenFile(LOG_FILE, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	LogFatal(err)
	log.SetOutput(openFile)
}

func UnknownAction(jsonPayload []byte) (string, error) {
	return "", errors.New("Unknown action given")
}

func NewActionFrom(val string) Action {
	if action, ok := actionsMap[strings.ToUpper(val)]; ok {
		return action
	}
	return UnknownAction
}
