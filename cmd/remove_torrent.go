package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gobuffalo/packr"
	. "github.com/lcserny/goutils"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

type RemoveTorrentsConfig struct {
	QBitTorrentWebUIPort string `json:"qBittorrentWebUIPort"`
}

var hashFlag = flag.String("hash", "", "the hash of the downloaded torrent to be removed")

func main() {
	initRemoveTorrentLogger("vm-removetorrent.log")

	args := os.Args[1:]
	if len(args) < 1 {
		_, err := fmt.Fprint(os.Stderr, "ERROR: Please provide `hash` flag arg\n")
		LogError(err)
		return
	}

	flag.Parse()
	if err := executeRemoveTorrentRequest(); err != nil {
		LogFatal(err)
	}
}

func executeRemoveTorrentRequest() error {
	values := url.Values{"hashes": {*hashFlag}}
	LogInfo(fmt.Sprintf("Sending request with values: %v", values))
	_, err := http.PostForm(getRemoveTorrentUrl(), values)
	return err
}

func getRemoveTorrentUrl() string {
	configFolder := packr.NewBox("../cfg/remove_torrents")
	content, err := configFolder.Find("config.json")
	LogFatal(err)

	var config RemoveTorrentsConfig
	err = json.Unmarshal(content, &config)
	LogFatal(err)

	return fmt.Sprintf("http://localhost:%s/command/delete", config.QBitTorrentWebUIPort)
}

func initRemoveTorrentLogger(logFile string) {
	openFile, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	LogFatal(err)
	writer := io.MultiWriter(os.Stdout, openFile)
	log.SetOutput(writer)
}
