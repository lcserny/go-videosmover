package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"github.com/lcserny/goutils"
	"os"
	"videosmover/pkg"
	"videosmover/pkg/action"
	"videosmover/pkg/config"
	"videosmover/pkg/delete"
	"videosmover/pkg/ext/godirwalk"
	"videosmover/pkg/ext/h2non"
	"videosmover/pkg/ext/json"
	"videosmover/pkg/ext/tmdb"
	"videosmover/pkg/move"
	"videosmover/pkg/output"
	"videosmover/pkg/search"
)


/* TODO: fix subtitles
If TV Show
	var countNrOfEpisodes
	If Subs|subs folder inside current folder
		Go inside folder
		var countOfSubfolders
		If there are same nr of subfolders here as nr of episodes
			Loop over each folder
				var nameOfFolder
				Go inside folder
				Find English subtitle file
				Rename it to nameOfFolder
				Move it up one level
				Go back one level
				Delete folder
			Move subtitle files one level up, to root of TV Show
			Go back one level
			Delete Subs|subs folder
*/


func main() {
	args := os.Args[1:]
	if len(args) != 3 {
		_, _ = fmt.Fprintln(os.Stderr, "ERROR: Please provide `config`, `action` and `payloadFile` flags")
		os.Exit(1)
	}

	cmdConfig := flag.String("config", "", "actions config file path")
	cmdAction := flag.String("action", "search", "action to execute")
	cmdPayload := flag.String("payload", "", "base64 encoded json payload")
	flag.Parse()

	codec := json.NewJsonCodec()
	cfg := config.MakeActionConfig(*cmdConfig, codec)
	goutils.InitFileLogger(cfg.LogFile)
	trashMover := core.NewTrashMover()
	mimeChecker := h2non.NewVideoChecker(cfg)
	videoPathWalker := godirwalk.NewVideoPathWalker(cfg)
	videoWebSearcher := tmdb.NewVideoWebSearcher(cfg)
	httpCache := core.NewHttpCacheStore(cfg.CacheAddress, "/get", "/set", "/close", codec)
	defer httpCache.Close()

	actionRepo := action.NewActionRepository()
	actionRepo.Register("delete", delete.NewAction(cfg, codec, trashMover))
	actionRepo.Register("move", move.NewAction(cfg, codec, trashMover))
	actionRepo.Register("output", output.NewAction(cfg, codec, videoWebSearcher, httpCache))
	actionRepo.Register("search", search.NewAction(cfg, codec, mimeChecker, videoPathWalker))

	b, err := base64.StdEncoding.DecodeString(*cmdPayload)
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
