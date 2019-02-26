package actions

import (
	. "github.com/lcserny/goutils"
	"log"
	"os"
	"path/filepath"
)

const (
	CONFIG_FILE  = "videosmover.properties"
	LOG_PATH_KEY = "log.path"
)

var AppProperties *ConfigProperties

func InitConfig(configDir string) {
	AppProperties = ReadPropertiesFile(filepath.Join(configDir, CONFIG_FILE))
	if AppProperties.HasProperty(LOG_PATH_KEY) {
		initLogger(AppProperties.GetPropertyAsString(LOG_PATH_KEY))
	}
}

func initLogger(logPath string) {
	openFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	LogFatal(err)
	log.SetOutput(openFile)
}
