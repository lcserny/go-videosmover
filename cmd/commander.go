package main

import (
	"flag"
	"fmt"
	"github.com/lcserny/go-videosmover/pkg/action"
	. "github.com/lcserny/goutils"
	"io/ioutil"
	"os"
)

var (
	commanderConfigsFlag     = flag.String("configs", "", "configs folder path")
	commanderActionFlag      = flag.String("action", "search", "action to execute")
	commanderPayloadFileFlag = flag.String("payloadFile", "", "path to payload file")
)

func main() {
	args := os.Args[1:]
	if len(args) != 3 {
		_, _ = fmt.Fprintln(os.Stderr, "ERROR: Please provide `configs`, `action` and `payloadFile` flags")
		return
	}

	InitFileLogger("vm-commander.log")

	flag.Parse()
	// FIXME: trigger the init()'s somehow or different approach to populate
	a := action.Retrieve(*commanderActionFlag)
	c := action.NewConfig(*commanderConfigsFlag, "actions.json")

	jsonBytes, err := ioutil.ReadFile(*commanderPayloadFileFlag)
	stopOnError(err)

	response, err := a.Execute(jsonBytes, c)
	stopOnError(err)

	_, err = fmt.Fprint(os.Stdout, response)
	stopOnError(err)
}

func stopOnError(err error) {
	if err != nil {
		_, err := fmt.Fprint(os.Stderr, err)
		LogFatal(err)
	}
}
