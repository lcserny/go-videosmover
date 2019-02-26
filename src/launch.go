package main

import (
	"flag"
	"fmt"
	"github.com/lcserny/go-videosmover/src/actions"
	. "github.com/lcserny/goutils"
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
	response, err := action.Execute(args[1])
	if err != nil {
		_, err := fmt.Fprint(os.Stderr, err)
		LogError(err)
		return
	}
	_, err = fmt.Fprint(os.Stdout, response)
	LogError(err)
}
