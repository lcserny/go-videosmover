package main

import (
	"flag"
	"fmt"
	"github.com/lcserny/goutils"
	"net/http"
	"net/url"
	"os"
)

func main() {
	args := os.Args[1:]
	if len(args) != 2 {
		_, _ = fmt.Fprintln(os.Stderr, "ERROR: Please provide `port` and `hash` flags")
		return
	}

	goutils.InitFileLogger("vm-removeqtorrent.log")

	port := flag.String("port", "", "the port of qtorrent's webUI")
	hash := flag.String("hash", "", "the hash of the downloaded torrent to be removed")
	flag.Parse()

	values := url.Values{"hashes": {*hash}}
	host := fmt.Sprintf("http://localhost:%s/command/delete", *port)
	if _, err := http.PostForm(host, values); err != nil {
		goutils.LogFatal(err)
	}
}
