package main

import (
	"flag"
	"fmt"
	"github.com/lcserny/goutils"
	"net/http"
	"net/url"
	"os"
	core "videosmover/pkg"
	"videosmover/pkg/ext/json"
)

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

	values := url.Values{"hashes": {*hash}}
	host := fmt.Sprintf("http://localhost:%s/command/delete", *port)
	if _, err := http.PostForm(host, values); err != nil {
		goutils.LogFatal(err)
	}

	var completed []*string // TODO: not string, pair of file name cu date
	httpCache.Get("downComplete", completed)

	downloadedFile := "de unde iau file downloaded data, din hash printr-un API de qTorrent"
	completed = append(completed, &downloadedFile)

	httpCache.Set("downComplete", completed)
}
