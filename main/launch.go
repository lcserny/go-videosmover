package main

import (
	"flag"
	"fmt"
	. "github.com/lcserny/go-videosmover"
	"os"
)

var actionFlag = flag.String("action", "search", "Please provide a `action` flag like: SEARCH")

func init() {
	flag.Parse()
}

func main() {
	args := os.Args[1:]
	argsLength := len(args)
	if argsLength < 1 {
		print("ERROR: No args provided, please provide `jsonPayloadFilePath`")
		return
	}

	jsonFile := args[0]
	if argsLength == 2 {
		jsonFile = args[1]
	}

	action := NewActionFrom(*actionFlag)
	response, err := action.Execute(jsonFile)
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, err)
		return
	}
	_, _ = fmt.Fprint(os.Stdout, response)
}
