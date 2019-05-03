package main

import (
	"flag"
	"fmt"
	"github.com/lcserny/goutils"
	"io/ioutil"
	"os"
	"path/filepath"
	"videosmover/pkg/action"
	"videosmover/pkg/delete"
	"videosmover/pkg/move"
	"videosmover/pkg/output"
	"videosmover/pkg/search"
)

func main() {
	args := os.Args[1:]
	if len(args) != 3 {
		_, _ = fmt.Fprintln(os.Stderr, "ERROR: Please provide `configs`, `action` and `payloadFile` flags")
		return
	}

	goutils.InitFileLogger("vm-commander.log")

	cmdConfig := flag.String("configs", "", "configs folder path")
	cmdAction := flag.String("action", "search", "action to execute")
	cmdPayload := flag.String("payloadFile", "", "path to payload file")
	flag.Parse()

	action.Register("delete", delete.NewAction())
	action.Register("move", move.NewAction())
	action.Register("output", output.NewAction())
	action.Register("search", search.NewAction())

	a := action.Retrieve(*cmdAction)
	c := action.NewConfig(filepath.Join(*cmdConfig, "actions.json"))

	jsonBytes, err := ioutil.ReadFile(*cmdPayload)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		goutils.LogFatal(err)
	}

	response, err := a.Execute(jsonBytes, c)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		goutils.LogFatal(err)
	}

	fmt.Fprint(os.Stdout, response)
}
