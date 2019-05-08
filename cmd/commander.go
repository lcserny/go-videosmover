package main

import (
	"flag"
	"fmt"
	"github.com/lcserny/goutils"
	"io/ioutil"
	"os"
	"path/filepath"
	"videosmover/pkg/action"
	"videosmover/pkg/config"
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

	codec := json.NewJsonCodec()
	trashMover := wastebasket.NewTrashMover()
	cfg := config.NewActionConfig(filepath.Join(*cmdConfig, "actions.json"), codec)

	action.Register("delete", delete.NewAction(cfg, codec, trashMover))
	action.Register("move", move.NewAction(cfg, codec, trashMover))
	action.Register("output", output.NewAction(cfg, codec)) // TODO: abstract dependencies, do this after
	action.Register("search", search.NewAction(cfg, codec)) // TODO: abstract dependencies, do this first

	b, err := ioutil.ReadFile(*cmdPayload)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		goutils.LogFatal(err)
	}

	a := action.Retrieve(*cmdAction)
	response, err := a.Execute(b)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		goutils.LogFatal(err)
	}

	fmt.Fprint(os.Stdout, response)
}
