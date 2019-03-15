package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/lcserny/go-videosmover/pkg/actions"
	. "github.com/lcserny/goutils"
	"io/ioutil"
	"log"
	"os"
)

var commanderActionFlag = flag.String("action", "search", "Please provide a `action` flag like: SEARCH")

func main() {
	initCommanderLogger()

	args := os.Args[1:]
	if len(args) != 2 {
		LogError(errors.New("ERROR: Please provide `action` flag and `jsonPayloadFilePath` args\n"))
		return
	}

	flag.Parse()
	action := actions.NewActionFrom(*commanderActionFlag)
	config := actions.GenerateActionConfig("../../cfg/commander", "actions.json")

	jsonBytes, err := ioutil.ReadFile(args[1])
	stopOnError(err)

	response, err := action(jsonBytes, config)
	stopOnError(err)

	_, err = fmt.Fprint(os.Stdout, response)
	stopOnError(err)
}

func initCommanderLogger() {
	openFile, err := os.OpenFile(GetAbsCurrentPathOf("vm-commander.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	LogFatal(err)
	log.SetOutput(openFile)
}

func stopOnError(err error) {
	if err != nil {
		_, err := fmt.Fprint(os.Stderr, err)
		LogFatal(err)
	}
}
