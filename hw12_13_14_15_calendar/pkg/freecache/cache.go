package freecache

import (
	"sync"

	"github.com/cespare/xxhash/v2"
	"github.com/coocood/freecache"
	"github.com/dmitrii-a/hw_go/hw12_13_14_15_calendar/internal/common"
)

// CacheDB is a freecache struct with mutex.
type CacheDB struct {
	cache *freecache.Cache
	mx    sync.RWMutex
	keys  map[uint64][]byte
}

// NewCacheDB init freecache.
func NewCacheDB(size int) *CacheDB {
	return &CacheDB{cache: freecache.NewCache(size), keys: make(map[uint64][]byte)}
}

// Get is a method for getting data.
func (db *CacheDB) Get(key []byte) ([]byte, error) {
	db.mx.RLock()
	got, err := db.cache.Get(key)
	db.mx.RUnlock()
	return got, err
}

// Set is a method for saving data with expire.
func (db *CacheDB) Set(key, val []byte, expireIn int) error {
	err := db.cache.Set(key, val, expireIn)
	if common.IsErr(err) {
		return err
	}
	db.mx.Lock()
	db.keys[xxhash.Sum64(key)] = key
	db.mx.Unlock()
	return nil
}

// Del is a method for deleting data.
func (db *CacheDB) Del(key []byte) (affected bool) {
	db.mx.Lock()
	defer db.mx.Unlock()
	return db.cache.Del(key)
}

// Keys is a method for getting all keys.
func (db *CacheDB) Keys() [][]byte {
	keys := make([][]byte, 0, len(db.keys))
	db.mx.RLock()
	for _, v := range db.keys {
		keys = append(keys, v)
	}
	db.mx.RUnlock()
	return keys
}

// Clear is a method for clearing cache.
func (db *CacheDB) Clear() {
	db.cache.Clear()
	db.mx.Lock()
	db.keys = make(map[uint64][]byte)
	db.mx.Unlock()
}
