package goutils

import (
	"github.com/pkg/errors"
	"log"
	"os"
	"strings"
)

const NO_MESSAGE = "--- no message provided ---"

func LogFatal(err error) {
	if err != nil {
		log.Fatalf("ERROR: %+v\n\n", errors.Wrap(err, NO_MESSAGE))
	}
}

func LogFatalWithMessage(message string, err error) {
	if err != nil {
		log.Fatalf("ERROR: %+v\n\n", errors.Wrap(err, message))
	}
}

func LogError(err error) {
	if err != nil {
		log.Printf("ERROR: %+v\n\n", errors.Wrap(err, NO_MESSAGE))
	}
}

func LogErrorWithMessage(message string, err error) {
	if err != nil {
		log.Printf("ERROR: %+v\n\n", errors.Wrap(err, message))
	}
}

func LogInfo(message string) {
	log.Printf("INFO: %s\n", message)
}

func LogWarning(message string) {
	log.Printf("WARN: %s\n", message)
}

// If file created needs to be where its executed (as opposed to where binary is situated) set `GOUTILS_EXEC_LOGINIT` env var to true
func InitFileLogger(logFileName string) {
	enabled, exists := os.LookupEnv("GOUTILS_EXEC_LOGINIT")
	if !exists || strings.ToLower(enabled) != "true" {
		logFileName = GetAbsCurrentPathOf(logFileName)
	}

	openFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	LogFatal(err)
	log.SetOutput(openFile)
}
