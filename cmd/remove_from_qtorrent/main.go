package main

import (
	"flag"
	"fmt"
	"github.com/lcserny/goutils"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
	core "videosmover/pkg"
	"videosmover/pkg/ext/json"
)

type TorrentData struct {
	savePath string
	date     time.Time
}

func main() {
	args := os.Args[1:]
	if len(args) != 3 {
		_, _ = fmt.Fprintln(os.Stderr, "ERROR: Please provide `logFile`, `port` and `hash` flags")
		os.Exit(1)
	}

	cacheAddress := "http://localhost:8076"
	codec := json.NewJsonCodec()
	httpCache := core.NewHttpCacheStore(cacheAddress, "/get", "/set", "/close", codec)
	defer httpCache.Close()

	logFile := flag.String("logFile", "", "app log file path")
	port := flag.String("port", "", "the port of qtorrent's webUI")
	hash := flag.String("hash", "", "the hash of the downloaded torrent to be removed")
	flag.Parse()

	goutils.InitFileLogger(*logFile)

	torrentDataUrl := fmt.Sprintf("http://localhost:%s/query/propertiesGeneral/%s", *port, *hash)
	resp, err := http.Get(torrentDataUrl)
	if err != nil {
		goutils.LogFatal(err)
	}

	values := url.Values{"hashes": {*hash}}
	deleteUrl := fmt.Sprintf("http://localhost:%s/command/delete", *port)
	if _, err := http.PostForm(deleteUrl, values); err != nil {
		goutils.LogFatal(err)
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		goutils.LogFatal(err)
	}
	var respData map[string]interface{}
	if err = codec.Decode(respBytes, &respData); err != nil {
		goutils.LogFatal(err)
	}
	savePath := respData["save_path"].(string)

	var completed []TorrentData
	if err = httpCache.Get("downComplete", &completed); err != nil {
		goutils.LogFatal(err)
	}

	// TODO: is this valid JSON converted?
	completed = append(completed, TorrentData{
		savePath: savePath,
		date:     time.Time{},
	})

	if err = httpCache.Set("downComplete", completed); err != nil {
		goutils.LogFatal(err)
	}
}
