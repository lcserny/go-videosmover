package actions

import (
	. "github.com/lcserny/goutils"
	"log"
	"os"
	"path/filepath"
)

const (
	LOG_FILE    = "videosmover.log"
	CONFIG_FILE = "videosmover.properties" // TODO: embed these properties?
)

var AppProperties *ConfigProperties

func init() {
	initLogger()
	initProperties()
}

func initProperties() {
	AppProperties = ReadPropertiesFile(CONFIG_FILE)
}

func initLogger() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	LogFatal(err)

	openFile, err := os.OpenFile(filepath.Join(dir, LOG_FILE), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	LogFatal(err)

	log.SetOutput(openFile)
}
