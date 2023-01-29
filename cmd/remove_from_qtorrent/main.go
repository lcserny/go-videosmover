package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
	core "videosmover/pkg"
	"videosmover/pkg/config"
	"videosmover/pkg/ext/json"

	"github.com/lcserny/goutils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	args := os.Args[1:]
	if len(args) != 2 {
		_, _ = fmt.Fprintln(os.Stderr, "ERROR: Please provide `cfgFile` and `hash` flags")
		os.Exit(1)
	}

	cfgFile := flag.String("cfgFile", "", "config file path")
	hash := flag.String("hash", "", "the hash of the downloaded torrent to be removed")
	flag.Parse()

	codec := json.NewJsonCodec()
	cfg := config.MakeRemoveTorrentConfig(*cfgFile, codec)
	goutils.InitFileLogger(cfg.LogFile)

	updateCache(*hash, cfg, codec)
	removeTorrent(*hash, cfg)
}

func removeTorrent(hash string, config *core.RemoveTorrentConfig) {
	deleteUrl := fmt.Sprintf("%s/delete", config.QTorrent.TorrentsUrl)
	data := url.Values{
		"hashes":      {hash},
		"deleteFiles": {"false"},
	}
	if _, err := http.PostForm(deleteUrl, data); err != nil {
		goutils.LogFatal(err)
	}
}

type TorrentFile struct {
	Name           string `json:"name" bson:"file_name"`
	Size           int64  `json:"size" bson:"file_size"`
	DateDownloaded string `json:"dateDownloaded" bson:"date_downloaded"`
}

func updateCache(hash string, config *core.RemoveTorrentConfig, codec core.Codec) {
	getFilesUrl := fmt.Sprintf("%s/files", config.QTorrent.TorrentsUrl)
	data := url.Values{
		"hash": {hash},
	}
	getFilesResp, err := http.PostForm(getFilesUrl, data)
	if err != nil {
		goutils.LogFatal(err)
	}
	getFilesBytes, err := ioutil.ReadAll(getFilesResp.Body)
	if err != nil {
		goutils.LogFatal(err)
	}
	var torrentFilesList []TorrentFile
	if err = codec.Decode(getFilesBytes, &torrentFilesList); err != nil {
		goutils.LogFatal(err)
	}

	dateDownloaded := time.Now().Format(time.RFC1123Z)
	bsons := make([]interface{}, len(torrentFilesList))
	for i, t := range torrentFilesList {
		t.DateDownloaded = dateDownloaded
		bsons[i] = t
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoDB.Url))
	if err != nil {
		goutils.LogFatal(err)
	}
	defer cancel()
	collection := client.Database(config.MongoDB.Database).Collection(config.MongoDB.Collection)
	if _, err = collection.InsertMany(context.Background(), bsons); err != nil {
		goutils.LogFatal(err)
	}
}
