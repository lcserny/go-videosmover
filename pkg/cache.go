package core

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
