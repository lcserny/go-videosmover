package actions

import (
	"github.com/gobuffalo/packr"
	. "github.com/lcserny/goutils"
	"log"
	"os"
)

const (
	BASE_CONF_FILE = "commander.properties"
	LOG_FILE       = "videosmover-commander.log"
)

// if these need to be accessible in other packages, put this in a `shared` package and export vars
var (
	appProperties *ConfigProperties
	configFolder  packr.Box
)

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
