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

// TODO: use packr to pack config file into binary or?
func main() {
	args := os.Args[1:]
	if len(args) < 3 {
		_, err := fmt.Fprint(os.Stderr, "ERROR: Please provide `action` flag, `configPath` and `jsonPayloadFilePath` args")
		LogError(err)
		return
	}

	actions.InitConfig(args[1])

	action := actions.NewActionFrom(*actionFlag)
	response, err := action.Execute(args[2])
	if err != nil {
		_, err := fmt.Fprint(os.Stderr, err)
		LogError(err)
		return
	}
	_, err = fmt.Fprint(os.Stdout, response)
	LogError(err)
}
