package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"
	core "videosmover/pkg"
	"videosmover/pkg/config"
	"videosmover/pkg/ext/firebase"
	"videosmover/pkg/ext/json"
	"videosmover/pkg/shutdown"
	"videosmover/pkg/web"

	"github.com/lcserny/goutils"
)

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		_, _ = fmt.Fprintln(os.Stderr, "ERROR: Please provide `config` flag")
		os.Exit(1)
	}

	cfgPath := flag.String("config", "", "path to proxy server config")
	flag.Parse()

	jsonCodec := json.NewJsonCodec()
	apiRequester := web.NewApiRequester(jsonCodec)
	c := config.MakeProxyConfig(*cfgPath, jsonCodec)
	cloudDatabase := firebase.NewCloudDatabase(c)
	cacheAddress := "http://localhost:8076"
	httpCache := core.NewHttpCacheStore(cacheAddress, "/get", "/set", "/close", jsonCodec)
	goutils.InitFileLogger(c.LogFile)

	mux := http.NewServeMux()
	for _, binCmd := range c.Bin {
		if len(binCmd.Type) <= 0 || strings.ToLower(binCmd.Type) == "bin" {
			mux.Handle(fmt.Sprintf("/exec-bin/%s", binCmd.Uri), web.NewBinExecutor(&binCmd, jsonCodec, apiRequester))
		}
		if strings.ToLower(binCmd.Type) == "java" {
			mux.Handle(fmt.Sprintf("/exec-java/%s", binCmd.Uri), web.NewJavaExecutor(&binCmd, jsonCodec))
		}
	}
	addInternalHandlers(mux, httpCache, jsonCodec)
	// startUDPListener(c)

	core.StartKeepWarmPing(c)

	cloudDatabase.Init()

	goutils.LogInfo(fmt.Sprintf("Started server on port %s...", c.Port))
	goutils.LogFatal(http.ListenAndServe(fmt.Sprintf(":%s", c.Port), mux))
}

func startUDPListener(proxyConfig *core.ProxyConfig) {
	go func() {
		udpAddr, err := net.ResolveUDPAddr("udp4", ":"+proxyConfig.UDPPort)
		if err != nil {
			goutils.LogFatal(err)
		}
		udpConn, err := net.ListenUDP("udp4", udpAddr)
		if err != nil {
			goutils.LogFatal(err)
		}

		for {
			buffer := make([]byte, 1024)
			_, addr, err := udpConn.ReadFromUDP(buffer)
			if err != nil {
				goutils.LogError(err)
				return
			}
			_, err = udpConn.WriteToUDP([]byte("connected"), addr)
			if err != nil {
				goutils.LogError(err)
				return
			}
		}
	}()
}

func addInternalHandlers(mux *http.ServeMux, cache core.CacheStore, codec core.Codec) {
	core.AddWarmPingEndpoint(mux)
	// addShutdownEndpoint(mux)
	addDownloadsHistoryEndpoint(mux, cache, codec)
}

func addDownloadsHistoryEndpoint(mux *http.ServeMux, cache core.CacheStore, codec core.Codec) {
	now := time.Now().Format(core.CacheKeyDatePattern)
	key := core.CacheKeyPrefix + now

	mux.HandleFunc("/downloadsCompleted", func(writer http.ResponseWriter, request *http.Request) {
		completed := make([]*core.TorrentData, 0)
		if err := cache.Get(key, &completed); err != nil {
			goutils.LogError(err)
		}

		jsonCompleted, err := codec.EncodeString(completed)
		if err != nil {
			goutils.LogError(err)
		}

		writer.Header().Set("Content-Type", codec.ContentType())
		writer.WriteHeader(http.StatusOK)
		fmt.Fprint(writer, jsonCompleted)
	})
}

func addShutdownEndpoint(mux *http.ServeMux) {
	if runtime.GOOS != "windows" {
		goutils.LogInfo(fmt.Sprintf("Shutdown available for windows only, OS found: %s", runtime.GOOS))
		return
	}

	mux.HandleFunc("/shutdown", func(writer http.ResponseWriter, request *http.Request) {
		values := request.URL.Query()
		secondsInt := "0"
		if seconds, exists := values["seconds"]; exists {
			secondsInt = seconds[0]
		}
		shutdown.Shutdown(secondsInt)
	})
}
