package web

import (
	"fmt"
	"github.com/lcserny/goutils"
	"net/http"
	"videosmover/pkg"
)

type cacheHandler struct {
	cacheStore core.CacheStore
	config     *core.CacheConfig
	codec      core.Codec
}

func NewCacheHandler(c core.CacheStore, cfg *core.CacheConfig, codec core.Codec) core.CacheHandler {
	return &cacheHandler{
		cacheStore: c,
		config:     cfg,
		codec:      codec,
	}
}

// TODO
func (ch *cacheHandler) Get(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		return
	}

	w.Header().Set("Content-Type", "application/json")

	goutils.LogInfo(fmt.Sprintf("%+v", r.PostFormValue("key")))
	goutils.LogInfo(fmt.Sprintf("%+v", r.FormValue("key")))

	// get post params
	// decode them?
	// save to cache
}

// TODO
func (ch *cacheHandler) Set(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		return
	}

	w.Header().Set("Content-Type", "application/json")

	goutils.LogInfo(fmt.Sprintf("%+v", r.PostFormValue("key")))
	goutils.LogInfo(fmt.Sprintf("%+v", r.PostFormValue("val")))
	goutils.LogInfo(fmt.Sprintf("%+v", r.FormValue("key")))
	goutils.LogInfo(fmt.Sprintf("%+v", r.FormValue("val")))

	// if not POST return
	// get post params
	// decode them
	// get from cache
	// encode back
	// return them
}
