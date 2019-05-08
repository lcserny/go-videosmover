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
	"videosmover/pkg/json"
	"videosmover/pkg/move"
	"videosmover/pkg/output"
	"videosmover/pkg/search"
	"videosmover/pkg/wastebasket"
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

	// TODO: move action in core
	// TODO: move core actions (search, delete etc) in domain also?
	codec := json.NewJsonCodec()
	trashMover := wastebasket.NewTrashMover()
	action.Register("delete", delete.NewAction(codec, trashMover))
	action.Register("move", move.NewAction(codec, trashMover))
	action.Register("output", output.NewAction(codec)) // TODO: abstract dependencies
	action.Register("search", search.NewAction(codec)) // TODO: abstract dependencies

	a := action.Retrieve(*cmdAction)
	c := action.NewConfig(filepath.Join(*cmdConfig, "actions.json"), codec)

	b, err := ioutil.ReadFile(*cmdPayload)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		goutils.LogFatal(err)
	}

	response, err := a.Execute(b, c)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		goutils.LogFatal(err)
	}

	fmt.Fprint(os.Stdout, response)
}
