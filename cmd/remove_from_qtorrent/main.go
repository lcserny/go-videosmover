package main

import (
	"flag"
	"fmt"
	"github.com/lcserny/goutils"
	"io/ioutil"
	"net/http"
	"os"
	"time"
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

	updateCache(port, hash, codec, httpCache)
	removeTorrent(hash, port)
}

func removeTorrent(hash *string, port *string) {
	deleteUrl := fmt.Sprintf("http://localhost:%s/api/v2/torrents/delete?hashes=%s&deleteFiles=false", *port, *hash)
	if _, err := http.Get(deleteUrl); err != nil {
		goutils.LogFatal(err)
	}
}

func updateCache(port *string, hash *string, codec core.Codec, httpCache core.CacheStore) {
	torrentPathGetUrl := fmt.Sprintf("http://localhost:%s/api/v2/torrents/properties?hash=%s", *port, *hash)
	pathResp, err := http.Get(torrentPathGetUrl)
	if err != nil {
		goutils.LogFatal(err)
	}
	torrentNameGetUrl := fmt.Sprintf("http://localhost:%s/api/v2/torrents/files?hash=%s", *port, *hash)
	nameResp, err := http.Get(torrentNameGetUrl)
	if err != nil {
		goutils.LogFatal(err)
	}
	pathRespBytes, err := ioutil.ReadAll(pathResp.Body)
	if err != nil {
		goutils.LogFatal(err)
	}
	var pathRespData map[string]interface{}
	if err = codec.Decode(pathRespBytes, &pathRespData); err != nil {
		goutils.LogFatal(err)
	}
	savePath := pathRespData["save_path"].(string)
	nameRespBytes, err := ioutil.ReadAll(nameResp.Body)
	if err != nil {
		goutils.LogFatal(err)
	}
	var nameRespData []map[string]interface{}
	if err = codec.Decode(nameRespBytes, &nameRespData); err != nil {
		goutils.LogFatal(err)
	}
	now := time.Now().Format(core.CacheKeyDatePattern)
	key := core.CacheKeyPrefix + now
	var completed []*core.TorrentData
	if err = httpCache.Get(key, &completed); err != nil {
		goutils.LogFatal(err)
	}
	hostname, _ := os.Hostname()
	for i := 0; i < len(nameRespData); i++ {
		name := nameRespData[i]["name"].(string)
		completed = append(completed, &core.TorrentData{
			Host:     hostname,
			SavePath: savePath + name,
			Date:     now,
		})
	}
	if err = httpCache.Set(key, completed); err != nil {
		goutils.LogFatal(err)
	}
}
