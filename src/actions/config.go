package actions

import (
	"github.com/gobuffalo/packr"
	. "github.com/lcserny/goutils"
	"log"
	"os"
)

const (
	TMDB_API_KEY   = "TMDB_API_KEY"
	BASE_CONF_FILE = "videosmover.properties"
	LOG_FILE       = "videosmover.log"
)

// if these need to be accessible in other packages, put this in a `shared` package and export vars
var (
	appProperties *ConfigProperties
	configFolder  packr.Box
	tmdbApiKey    string
)

func init() {
	initLogger()

	configFolder = packr.NewBox("../../config")
	content, err := configFolder.FindString(BASE_CONF_FILE)
	LogFatal(err)

	appProperties = ReadProperties(content)

	initTMDBApiKey()
}

func initTMDBApiKey() {
	key, exists := os.LookupEnv(TMDB_API_KEY)
	if exists {
		tmdbApiKey = key
	}
}

func initLogger() {
	openFile, err := os.OpenFile(LOG_FILE, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	LogFatal(err)
	log.SetOutput(openFile)
}
