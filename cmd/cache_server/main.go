package main

import (
	"flag"
	"fmt"
	"github.com/VictoriaMetrics/fastcache"
	"github.com/lcserny/goutils"
	"net/http"
	"os"
	"videosmover/pkg/config"
	"videosmover/pkg/ext/json"
)

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		_, _ = fmt.Fprintln(os.Stderr, "ERROR: Please provide `config` flag")
		return
	}

	cfgPath := flag.String("config", "", "path to cache server config")
	flag.Parse()

	jsonCodec := json.NewJsonCodec()
	c := config.MakeCacheServerConfig(*cfgPath, jsonCodec)
	goutils.InitFileLogger(c.LogFile)
	cache := fastcache.LoadFromFileOrNew(c.CacheDBPath, c.MaxSizeBytes)

	server := http.NewServeMux()
	server.HandleFunc("/get", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		if request.Method != http.MethodPost {
			writer.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		writer.WriteHeader(http.StatusOK)
		fmt.Fprint(writer, string(cache.Get(nil, []byte(request.PostFormValue("key")))))
	})
	server.HandleFunc("/set", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		if request.Method != http.MethodPost {
			writer.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		writer.WriteHeader(http.StatusOK)
		cache.Set([]byte(request.PostFormValue("key")), []byte(request.PostFormValue("val")))
	})
	server.HandleFunc("/close", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		cache.SaveToFile(c.CacheDBPath)
	})

	goutils.LogInfo(fmt.Sprintf("Started server on port %s...", c.Port))
	goutils.LogFatal(http.ListenAndServe(fmt.Sprintf(":%s", c.Port), server))
}
