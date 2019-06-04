package core

import "net/http"

type CacheStore interface {
	Set(key string, val interface{}) error
	Get(key string) (val interface{}, found bool)
	Persist() error
}

type CacheHandler interface {
	Get(w http.ResponseWriter, r *http.Request)
	Set(w http.ResponseWriter, r *http.Request)
}
