package database

import (
	"github.com/patrickmn/go-cache"
)

type GoCache struct {
	Cache *cache.Cache
}

func (g GoCache) Get(key []byte) [][]byte {
	value, _ := g.Cache.Get(string(key))
	bytes, ok := value.([][]byte)
	if !ok {
		return [][]byte{value.([]byte)}
	} else {
		return bytes
	}
}

func (g GoCache) Put(key []byte, value []byte) {
	g.Cache.Set(string(key), value, cache.NoExpiration)
}

func (g GoCache) Append(key, value []byte) {
	slice := g.Get(key)
	slice = append(slice, value)
	g.Cache.Set(string(key), slice, cache.NoExpiration)
}

func (g GoCache) Del(key []byte) {
	g.Cache.Delete(string(key))
}