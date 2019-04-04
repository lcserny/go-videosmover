package main

import (
	"flag"
	"fmt"
	. "github.com/lcserny/goutils"
	"io/ioutil"
	"os"

	"github.com/lcserny/go-videosmover/pkg/action"
	_ "github.com/lcserny/go-videosmover/pkg/delete"
	_ "github.com/lcserny/go-videosmover/pkg/move"
	_ "github.com/lcserny/go-videosmover/pkg/output"
	_ "github.com/lcserny/go-videosmover/pkg/search"
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
