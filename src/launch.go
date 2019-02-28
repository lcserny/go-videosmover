package main

import (
	"flag"
	"fmt"
	"github.com/lcserny/go-videosmover/src/actions"
	. "github.com/lcserny/goutils"
	"io/ioutil"
	"os"
)

var actionFlag = flag.String("action", "search", "Please provide a `action` flag like: SEARCH")

func init() {
	flag.Parse()
}

func main() {
	args := os.Args[1:]
	if len(args) < 2 {
		_, err := fmt.Fprint(os.Stderr, "ERROR: Please provide `action` flag and `jsonPayloadFilePath` args")
		LogError(err)
		return
	}

	action := actions.NewActionFrom(*actionFlag)

	jsonBytes, err := ioutil.ReadFile(args[1])
	stopOnError(err)

	response, err := action(jsonBytes)
	stopOnError(err)

	_, err = fmt.Fprint(os.Stdout, response)
	LogError(err)
}

func stopOnError(err error) {
	if err != nil {
		_, err := fmt.Fprint(os.Stderr, err)
		LogFatal(err)
	}
}
