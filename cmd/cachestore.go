package main

import (
	"flag"
	"fmt"
	"github.com/lcserny/goutils"
	"net/http"
	"os"
	"time"
	"videosmover/pkg/config"
	"videosmover/pkg/ext/fastcache"
	"videosmover/pkg/ext/json"
	"videosmover/pkg/web"
)

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		_, _ = fmt.Fprintln(os.Stderr, "ERROR: Please provide `config` flag")
		return
	}

	goutils.InitFileLogger("vm-cachestore.log")

	cfgPath := flag.String("config", "", "path to cache store config")
	flag.Parse()

	jsonCodec := json.NewJsonCodec()
	c := config.MakeCacheConfig(*cfgPath, jsonCodec)
	cache := fastcache.NewCacheStore(c.CachePath, c.MaxSizeBytes)
	handler := web.NewCacheHandler(cache, c, jsonCodec)

	mux := http.NewServeMux()
	mux.HandleFunc("/cache/get", handler.Get)
	mux.HandleFunc("/cache/set", handler.Set)

	go func() {
		for range time.NewTicker(time.Duration(c.PersistenceIntervalMs) * time.Millisecond).C {
			goutils.LogError(cache.Persist())
		}
	}()

	goutils.LogInfo(fmt.Sprintf("Started server on port %s...", c.Port))
	goutils.LogFatal(http.ListenAndServe(fmt.Sprintf(":%s", c.Port), mux))
}
