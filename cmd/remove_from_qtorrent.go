package main

import (
	"flag"
	"fmt"
	utils "github.com/lcserny/goutils"
	"net/http"
	"net/url"
	"os"
)

var portFlag = flag.String("port", "", "the port of qBittorrent's webUI")
var hashFlag = flag.String("hash", "", "the hash of the downloaded torrent to be removed")

func main() {
	args := os.Args[1:]
	if len(args) != 2 {
		_, _ = fmt.Fprintln(os.Stderr, "ERROR: Please provide `port` and `hash` flags")
		return
	}

	utils.InitFileLogger("vm-removeqtorrent.log")

	flag.Parse()
	if err := executeRemoveTorrentRequest(); err != nil {
		utils.LogFatal(err)
	}
}

func executeRemoveTorrentRequest() error {
	values := url.Values{"hashes": {*hashFlag}}
	_, err := http.PostForm(fmt.Sprintf("http://localhost:%s/command/delete", *portFlag), values)
	return err
}
