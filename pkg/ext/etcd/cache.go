package etcd

import (
	core "videosmover/pkg"
)

type cacheStore struct {
}

func NewCacheStore(connectionAddress string) core.CacheStore {
	cs := new(cacheStore)

	return cs
}

func (cs *cacheStore) Set(key string, val interface{}) error {
	return nil
}

func (cs *cacheStore) Get(key string, valHolderPointer interface{}) error {
	return nil
}

func (cs *cacheStore) Close() {
}
