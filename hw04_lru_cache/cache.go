package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type cacheItem struct {
	key   Key
	value interface{}
}

type lruCache struct {
	sync.Mutex
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func (cache *lruCache) remove(item *ListItem) {
	cache.queue.Remove(item)
	delete(cache.items, item.Value.(*cacheItem).key)
}

func (cache *lruCache) createCacheItem(key Key, value interface{}) *cacheItem {
	return &cacheItem{key: key, value: value}
}

func (cache *lruCache) Set(key Key, value interface{}) bool {
	cache.Lock()
	defer cache.Unlock()
	if item, exist := cache.items[key]; exist {
		item.Value = cache.createCacheItem(key, value)
		cache.queue.MoveToFront(item)
		return true
	}
	if cache.queue.Len() >= cache.capacity {
		lastItem := cache.queue.Back()
		cache.remove(lastItem)
	}
	item := cache.queue.PushFront(cache.createCacheItem(key, value))
	cache.items[key] = item
	return false
}

func (cache *lruCache) Get(key Key) (interface{}, bool) {
	cache.Lock()
	defer cache.Unlock()
	if item, exist := cache.items[key]; exist {
		cache.queue.MoveToFront(item)
		return item.Value.(*cacheItem).value, true
	}
	return nil, false
}

func (cache *lruCache) Clear() {
	cache.Lock()
	defer cache.Unlock()
	cache.items = make(map[Key]*ListItem, cache.capacity)
	cache.queue = NewList()
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
