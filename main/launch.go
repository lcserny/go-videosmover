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
	if argsLength < 2 {
		_, _ = fmt.Fprint(os.Stderr, "ERROR: Please provide `action` flag and `jsonPayloadFilePath` arg")
		return
	}

	action := NewActionFrom(*actionFlag)
	response, err := action.Execute(args[1])
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, err)
		return
	}
	_, _ = fmt.Fprint(os.Stdout, response)
}
