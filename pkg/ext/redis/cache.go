package redis

import (
	"bytes"
	"encoding/gob"
	"github.com/lcserny/goutils"
	"github.com/mediocregopher/radix/v3"
	"videosmover/pkg"
)

type cacheStore struct {
	available    bool
	redisAddress string
	connPoolSize int
	client       *radix.Pool
}

func (cs *cacheStore) Set(key string, val interface{}) error {
	if !cs.available {
		return nil
	}
	enc, err := cs.marshal(val)
	if err != nil {
		return err
	}
	return cs.client.Do(radix.Cmd(nil, "SET", key, string(enc)))
}

func (cs *cacheStore) Get(key string, valHolderPointer interface{}) error {
	if !cs.available {
		return nil
	}
	var b []byte
	if err := cs.client.Do(radix.Cmd(&b, "GET", key)); err != nil {
		return err
	}
	if len(b) == 0 {
		return nil
	}
	if err := cs.unmarshal(b, valHolderPointer); err != nil {
		return err
	}
	return nil
}

func (cs cacheStore) marshal(v interface{}) ([]byte, error) {
	b := new(bytes.Buffer)
	err := gob.NewEncoder(b).Encode(v)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (cs cacheStore) unmarshal(data []byte, v interface{}) error {
	b := bytes.NewBuffer(data)
	return gob.NewDecoder(b).Decode(v)
}

func NewCacheStore(redisAddress string, connPoolSize int) core.CacheStore {
	cs := &cacheStore{
		redisAddress: redisAddress,
		connPoolSize: connPoolSize,
	}
	client, err := radix.NewPool("tcp", redisAddress, connPoolSize)
	if err != nil {
		goutils.LogError(err)
		return cs
	}
	cs.client = client
	cs.available = true
	return cs
}