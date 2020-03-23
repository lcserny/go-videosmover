package main

import (
	"flag"
	"fmt"
	"github.com/VictoriaMetrics/fastcache"
	"github.com/lcserny/goutils"
	"io/ioutil"
	"net/http"
	"os"
	"videosmover/pkg/config"
	"videosmover/pkg/ext/json"
)

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		_, _ = fmt.Fprintln(os.Stderr, "ERROR: Please provide `config` flag")
		os.Exit(1)
	}

	cfgPath := flag.String("config", "", "path to cache server config")
	flag.Parse()

	jsonCodec := json.NewJsonCodec()
	c := config.MakeCacheServerConfig(*cfgPath, jsonCodec)
	goutils.InitFileLogger(c.LogFile)
	cache := fastcache.LoadFromFileOrNew(c.CacheDBPath, c.MaxSizeBytes)

	server := http.NewServeMux()
	server.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
	})
	server.HandleFunc("/get", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", jsonCodec.ContentType())
		if request.Method != http.MethodPost {
			writer.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		writer.WriteHeader(http.StatusOK)

		reqData, err := ioutil.ReadAll(request.Body)
		if err != nil {
			goutils.LogError(err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		var data map[string]string
		if err = jsonCodec.Decode(reqData, &data); err != nil {
			goutils.LogError(err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		fmt.Fprint(writer, string(cache.Get(nil, []byte(data["key"]))))
	})
	server.HandleFunc("/set", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost {
			writer.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		writer.WriteHeader(http.StatusOK)

		reqData, err := ioutil.ReadAll(request.Body)
		if err != nil {
			goutils.LogError(err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		var data map[string]string
		if err = jsonCodec.Decode(reqData, &data); err != nil {
			goutils.LogError(err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		cache.Set([]byte(data["key"]), []byte(data["val"]))
	})
	server.HandleFunc("/close", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
		go func() {
			cache.SaveToFile(c.CacheDBPath)
		}()
	})

	goutils.LogInfo(fmt.Sprintf("Started server on port %s...", c.Port))
	goutils.LogFatal(http.ListenAndServe(fmt.Sprintf(":%s", c.Port), server))
}
