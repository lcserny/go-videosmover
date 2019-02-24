package main

import (
	"flag"
	"fmt"
	. "github.com/lcserny/go-videosmover"
	. "github.com/lcserny/goutils"
	"log"
	"os"
	"path/filepath"
)

var actionFlag = flag.String("action", "search", "Please provide a `action` flag like: SEARCH")

func init() {
	flag.Parse()
	initLogger()
}

func initLogger() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	LogFatal(err)

	openFile, err := os.OpenFile(filepath.Join(dir, "videosmover.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	LogFatal(err)

	log.SetOutput(openFile)
}

func main() {
	args := os.Args[1:]
	argsLength := len(args)
	if argsLength < 2 {
		_, err := fmt.Fprint(os.Stderr, "ERROR: Please provide `action` flag and `jsonPayloadFilePath` arg")
		LogError(err)
		return
	}

	action := NewActionFrom(*actionFlag)
	response, err := action.Execute(args[1])
	if err != nil {
		_, err := fmt.Fprint(os.Stderr, err)
		LogError(err)
		return
	}
	_, err = fmt.Fprint(os.Stdout, response)
	LogError(err)
}
