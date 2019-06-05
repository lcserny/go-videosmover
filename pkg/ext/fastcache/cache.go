package fastcache

import (
	"bytes"
	"encoding/gob"
	"github.com/VictoriaMetrics/fastcache"
	"videosmover/pkg"
)

type cacheStore struct {
	cacheFile    string
	maxSizeBytes int
	cache        *fastcache.Cache
}

func NewCacheStore(cacheFile string, maxSizeBytes int) core.CacheStore {
	cs := new(cacheStore)
	cs.cacheFile = cacheFile
	cs.maxSizeBytes = maxSizeBytes
	cs.cache = fastcache.LoadFromFileOrNew(cacheFile, maxSizeBytes)
	return cs
}

func (cs *cacheStore) Set(key string, val interface{}) error {
	enc, err := cs.marshal(val)
	if err != nil {
		return err
	}
	cs.cache.Set([]byte(key), enc)
	return cs.cache.SaveToFile(cs.cacheFile)
}

func (cs cacheStore) marshal(v interface{}) ([]byte, error) {
	b := new(bytes.Buffer)
	err := gob.NewEncoder(b).Encode(v)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (cs *cacheStore) Get(key string, valHolderPointer interface{}) error {
	enc := cs.cache.Get(nil, []byte(key))
	if enc == nil {
		return nil
	}
	if err := cs.unmarshal(enc, valHolderPointer); err != nil {
		return err
	}
	return nil
}

func (cs cacheStore) unmarshal(data []byte, v interface{}) error {
	b := bytes.NewBuffer(data)
	return gob.NewDecoder(b).Decode(v)
}
