package core

type CacheStore interface {
	Set(key string, val interface{}) error
	Get(key string) (val interface{}, found bool)
}
