package etcd

import (
	"context"
	"github.com/coreos/etcd/client"
	"github.com/lcserny/goutils"
	"time"
	core "videosmover/pkg"
)

type cacheStore struct {
	available bool
	api       client.KeysAPI
	codec     core.Codec
}

func NewCacheStore(connectionAddress string, codec core.Codec) core.CacheStore {
	cs := new(cacheStore)
	cfg := client.Config{Endpoints: []string{connectionAddress}, HeaderTimeoutPerRequest: time.Second}
	c, err := client.New(cfg)
	if err != nil {
		goutils.LogError(err)
		return cs
	}
	cs.api = client.NewKeysAPI(c)
	cs.codec = codec
	cs.available = true
	return cs
}

func (cs *cacheStore) Set(key string, val interface{}) error {
	if !cs.available {
		return nil
	}

	enc, err := cs.codec.EncodeString(val)
	if err != nil {
		return err
	}
	_, err = cs.api.Set(context.Background(), key, enc, nil)
	return err
}

func (cs *cacheStore) Get(key string, valHolderPointer interface{}) error {
	if !cs.available {
		return nil
	}

	resp, err := cs.api.Get(context.Background(), key, nil)
	if err != nil {
		return err
	}
	if err := cs.codec.Decode([]byte(resp.Node.Value), valHolderPointer); err != nil {
		return err
	}
	return nil
}

func (cs *cacheStore) Close() {
}
