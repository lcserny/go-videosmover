package core

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

type CacheStore interface {
	Set(key string, val interface{}) error
	Get(key string, valHolderPointer interface{}) error
	Close()
}

type MockCacheStore struct {
}

func (MockCacheStore) Set(key string, val interface{}) error {
	return nil
}

func (MockCacheStore) Get(key string, valHolderPointer interface{}) error {
	return nil
}

func (MockCacheStore) Close() {
}

type httpCacheStore struct {
	available bool
	address   string
	getURI    string
	setURI    string
	closeURI  string
	codec     Codec
}

func (cs httpCacheStore) Set(key string, val interface{}) error {
	if !cs.available {
		return nil
	}

	enc, err := cs.codec.EncodeString(val)
	if err != nil {
		return err
	}

	postBytes, err := cs.codec.EncodeBytes(map[string]string{
		"key": key,
		"val": enc,
	})
	if err != nil {
		return err
	}

	resp, err := http.Post(cs.address+cs.setURI, "application.json", bytes.NewBuffer(postBytes))
	defer resp.Body.Close()
	return err
}

func (cs httpCacheStore) Get(key string, valHolderPointer interface{}) error {
	if !cs.available {
		return nil
	}

	postBytes, err := cs.codec.EncodeBytes(map[string]string{"key": key})
	if err != nil {
		return err
	}

	resp, err := http.Post(cs.address+cs.getURI, "application.json", bytes.NewBuffer(postBytes))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return cs.codec.Decode([]byte(body), valHolderPointer)
}

func (cs httpCacheStore) Close() {
	if cs.available {
		http.Get(cs.address + cs.closeURI)
	}
}

func NewHttpCacheStore(address, getURI, setURI, closeURI string, codec Codec) CacheStore {
	cs := new(httpCacheStore)
	cs.address = address
	cs.getURI = getURI
	cs.setURI = setURI
	cs.closeURI = closeURI
	cs.codec = codec

	resp, err := http.Get(address)
	if err != nil || resp.StatusCode != 200 {
		return cs
	}
	cs.available = true

	return cs
}
