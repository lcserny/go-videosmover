package actions

import (
	"fmt"
	"github.com/gobuffalo/packr"
	. "github.com/lcserny/goutils"
	"log"
	"os"
	"runtime"
)

const (
	LOG_PATH_KEY = "log.path"
)

var AppProperties *ConfigProperties

func init() {
	content, err := packr.NewBox("../config").FindString(fmt.Sprintf("videosmover_%s.properties", runtime.GOOS))
	LogFatal(err)

	AppProperties = ReadProperties(content)
	if AppProperties.HasProperty(LOG_PATH_KEY) {
		initLogger(AppProperties.GetPropertyAsString(LOG_PATH_KEY))
	}
}

func initLogger(logPath string) {
	openFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	LogFatal(err)
	log.SetOutput(openFile)
}
