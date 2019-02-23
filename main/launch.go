package main

import (
	"errors"
	"flag"
	"github.com/lcserny/go-videosmover"
	. "github.com/lcserny/goutils"
	"os"
	"strings"
)

type Command int

const (
	UNKNOWN Command = -1
	SEARCH  Command = 0
)

var commandFlag = flag.String("command", "search", "Please provide a `command` flag like: UNUSED")

func init() {
	flag.Parse()
}

func main() {
	var response string
	args := os.Args[2:]
	command := newCommandFrom(*commandFlag)
	switch command {
	case SEARCH:
		if len(args) != 1 {
			response = "ERROR: no args passed"
			break
		}
		response = go_videosmover.Search(args[0])
		break
	case UNKNOWN:
		LogFatal(errors.New("Unknown command given" + *commandFlag))
	}
	// TODO: this writes to StdErr why?
	print(response)
}

func newCommandFrom(val string) Command {
	switch strings.ToUpper(val) {
	case "SEARCH":
		return SEARCH
	}
	return UNKNOWN
}
