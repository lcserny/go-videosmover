package etcd

import (
	"bytes"
	"context"
	"encoding/gob"
	"github.com/coreos/etcd/client"
	"github.com/lcserny/goutils"
	"time"
	core "videosmover/pkg"
)

type cacheStore struct {
	api client.KeysAPI
}

func NewCacheStore(connectionAddress string) core.CacheStore {
	cs := new(cacheStore)
	cfg := client.Config{Endpoints: []string{connectionAddress}, Transport: client.DefaultTransport, HeaderTimeoutPerRequest: time.Second}
	c, err := client.New(cfg)
	if err != nil {
		goutils.LogFatal(err)
	}
	cs.api = client.NewKeysAPI(c)
	return cs
}

func (cs *cacheStore) Set(key string, val interface{}) error {
	enc, err := cs.marshal(val)
	if err != nil {
		return err
	}
	_, err = cs.api.Set(context.Background(), key, enc, nil)
	return err
}

func (cs *cacheStore) Get(key string, valHolderPointer interface{}) error {
	resp, err := cs.api.Get(context.Background(), key, nil)
	if err != nil {
		return err
	}
	if err := cs.unmarshal(resp.Node.Value, valHolderPointer); err != nil {
		return err
	}
	return nil
}

func (cs *cacheStore) Close() {
}

func (cs cacheStore) marshal(v interface{}) (string, error) {
	b := new(bytes.Buffer)
	err := gob.NewEncoder(b).Encode(v)
	if err != nil {
		return "", err
	}
	return b.String(), nil
}

func (cs cacheStore) unmarshal(data string, v interface{}) error {
	b := bytes.NewBufferString(data)
	return gob.NewDecoder(b).Decode(v)
}
