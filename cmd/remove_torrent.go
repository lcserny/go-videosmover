package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gobuffalo/packr"
	. "github.com/lcserny/goutils"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

type RemoveTorrentsConfig struct {
	QBitTorrentWebUIPort string `json:"qBittorrentWebUIPort"`
}

var hashFlag = flag.String("hash", "", "the hash of the downloaded torrent to be removed")

func main() {
	initRemoveTorrentLogger()

	args := os.Args[1:]
	if len(args) < 1 {
		LogError(errors.New("ERROR: Please provide `hash` flag\n"))
		return
	}

	flag.Parse()
	if err := executeRemoveTorrentRequest(); err != nil {
		LogFatal(err)
	}
}

func executeRemoveTorrentRequest() error {
	values := url.Values{"hashes": {*hashFlag}}
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

func initRemoveTorrentLogger() {
	logfile := filepath.Join(filepath.Dir(os.Args[0]), "vm-removetorrent.log")
	openFile, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	LogFatal(err)
	log.SetOutput(openFile)
}
