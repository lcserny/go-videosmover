package boltdb

import (
	"bytes"
	"encoding/gob"
	"github.com/lcserny/goutils"
	"github.com/schollz/boltdb-server/connect"
	core "videosmover/pkg"
)

type cacheStore struct {
	available  bool
	bucketName string
	conn       *connect.Connection
}

func NewCacheStore(connectionAddress, dbName, bucketName string) core.CacheStore {
	cs := new(cacheStore)
	cs.bucketName = bucketName
	connection, err := connect.Open(connectionAddress, dbName)
	if err != nil {
		goutils.LogError(err)
		return cs
	}
	if err = connection.CreateBuckets([]string{bucketName}); err != nil {
		goutils.LogError(err)
		return cs
	}
	cs.available = true
	cs.conn = connection
	return cs
}

func (cs *cacheStore) Set(key string, val interface{}) error {
	if !cs.available {
		return nil
	}
	enc, err := cs.marshal(val)
	if err != nil {
		return err
	}
	return cs.conn.Post(cs.bucketName, map[string]string{key: string(enc)})
}

func (cs *cacheStore) Get(key string, valHolderPointer interface{}) error {
	if !cs.available {
		return nil
	}
	resMap, err := cs.conn.Get(cs.bucketName, []string{key})
	if err != nil {
		return err
	}
	enc, ok := resMap[key]
	if !ok || len(enc) == 0 {
		return nil
	}
	if err := cs.unmarshal([]byte(enc), valHolderPointer); err != nil {
		return err
	}
	return nil
}

func (cs *cacheStore) Close() {
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
