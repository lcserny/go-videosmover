package main

import (
	"errors"
	"flag"
	"fmt"
	. "github.com/lcserny/goutils"
	"log"
	"net/http"
	"net/url"
	"os"
)

var portFlag = flag.String("port", "", "the port of qBittorrent's webUI")
var hashFlag = flag.String("hash", "", "the hash of the downloaded torrent to be removed")

func main() {
	initRemoveTorrentLogger()

	args := os.Args[1:]
	if len(args) != 2 {
		LogError(errors.New("ERROR: Please provide `port` and `hash` flags\n"))
		return
	}

	flag.Parse()
	if err := executeRemoveTorrentRequest(); err != nil {
		LogFatal(err)
	}
}

func executeRemoveTorrentRequest() error {
	values := url.Values{"hashes": {*hashFlag}}
	_, err := http.PostForm(fmt.Sprintf("http://localhost:%s/command/delete", *portFlag), values)
	return err
}

func initRemoveTorrentLogger() {
	openFile, err := os.OpenFile(GetAbsCurrentPathOf("vm-removeqtorrent.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	LogFatal(err)
	log.SetOutput(openFile)
}
