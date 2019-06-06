package main

import (
	"flag"
	"fmt"
	"github.com/lcserny/goutils"
	"io/ioutil"
	"os"
	"videosmover/pkg/action"
	"videosmover/pkg/config"
	"videosmover/pkg/delete"
	"videosmover/pkg/ext/etcd"
	"videosmover/pkg/ext/godirwalk"
	"videosmover/pkg/ext/h2non"
	"videosmover/pkg/ext/json"
	"videosmover/pkg/ext/tmdb"
	"videosmover/pkg/ext/wastebasket"
	"videosmover/pkg/move"
	"videosmover/pkg/output"
	"videosmover/pkg/search"
)

func main() {
	args := os.Args[1:]
	if len(args) != 3 {
		_, _ = fmt.Fprintln(os.Stderr, "ERROR: Please provide `config`, `action` and `payloadFile` flags")
		return
	}

	cmdConfig := flag.String("config", "", "actions config file path")
	cmdAction := flag.String("action", "search", "action to execute")
	cmdPayload := flag.String("payloadFile", "", "path to payload file")
	flag.Parse()

	codec := json.NewJsonCodec()
	cfg := config.MakeActionConfig(*cmdConfig, codec)
	goutils.InitFileLogger(cfg.LogFile)
	trashMover := wastebasket.NewTrashMover()
	mimeChecker := h2non.NewVideoChecker(cfg)
	videoPathWalker := godirwalk.NewVideoPathWalker(cfg)
	videoWebSearcher := tmdb.NewVideoWebSearcher()
	etcdCache := etcd.NewCacheStore(cfg.CacheAddress, codec)
	defer etcdCache.Close()

	actionRepo := action.NewActionRepository()
	actionRepo.Register("delete", delete.NewAction(cfg, codec, trashMover))
	actionRepo.Register("move", move.NewAction(cfg, codec, trashMover))
	actionRepo.Register("output", output.NewAction(cfg, codec, videoWebSearcher, etcdCache))
	actionRepo.Register("search", search.NewAction(cfg, codec, mimeChecker, videoPathWalker))

	b, err := ioutil.ReadFile(*cmdPayload)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		goutils.LogFatal(err)
	}

	a := actionRepo.Retrieve(*cmdAction)
	response, err := a.Execute(b)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		goutils.LogFatal(err)
	}

	fmt.Fprint(os.Stdout, response)
}
