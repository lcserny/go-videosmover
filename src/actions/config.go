package actions

import (
	"fmt"
	. "github.com/lcserny/goutils"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

const (
	PROPERTIES_FILE = "videosmover_%s.properties"
	LOG_PATH_KEY    = "log.path"
)

var AppProperties *ConfigProperties

func init() {
	AppProperties = ReadPropertiesFile(filepath.Join("config", fmt.Sprintf(PROPERTIES_FILE, runtime.GOOS)))
	if AppProperties.HasProperty(LOG_PATH_KEY) {
		initLogger(AppProperties.GetPropertyAsString(LOG_PATH_KEY))
	}
}

func initLogger(logPath string) {
	openFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	LogFatal(err)
	log.SetOutput(openFile)
}
