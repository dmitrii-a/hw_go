package freecache

import (
	"github.com/cespare/xxhash/v2"
	"github.com/coocood/freecache"
	"github.com/dmitrii-a/hw_go/hw12_13_14_15_calendar/internal/common"
)

// CacheDB is a freecache struct with mutex.
type CacheDB struct {
	cache *freecache.Cache
	keys  map[uint64][]byte
}

// NewCacheDB init freecache.
func NewCacheDB(size int) *CacheDB {
	return &CacheDB{cache: freecache.NewCache(size), keys: make(map[uint64][]byte)}
}

// Get is a method for getting data.
func (db *CacheDB) Get(key []byte) ([]byte, error) {
	got, err := db.cache.Get(key)
	return got, err
}

// Set is a method for saving data with expire.
func (db *CacheDB) Set(key, val []byte, expireIn int) error {
	err := db.cache.Set(key, val, expireIn)
	if common.IsErr(err) {
		return err
	}
	db.keys[xxhash.Sum64(key)] = key
	return nil
}

// Del is a method for deleting data.
func (db *CacheDB) Del(key []byte) (affected bool) {
	return db.cache.Del(key)
}

// Keys is a method for getting all keys.
func (db *CacheDB) Keys() [][]byte {
	keys := make([][]byte, 0, len(db.keys))
	for _, v := range db.keys {
		keys = append(keys, v)
	}
	return keys
}

// Clear is a method for clearing cache.
func (db *CacheDB) Clear() {
	db.cache.Clear()
	db.keys = make(map[uint64][]byte)
}
