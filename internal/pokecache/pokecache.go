package pokecache

import (
	"sync"
	"time"
)

// Cache
//
// this is the exposed struct that will be accessed
type Cache struct {
	entires map[string]cacheEntry
	mux     sync.Mutex
}

// cacheEntry
//
// responsible for storing raw []byte
// as well as a time that it was created at
type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

// initializes a new Cache
//
// creates new cache that is valid for the given duration
func NewCache(interval time.Duration) *Cache {

	return &Cache{}
}

// adds a new entry to the cache
//
// take a key and a val to add to the Cache map
func (c *Cache) Add(name string, data []byte) {

}

// gets an entry to the cache
//
// takes a key
// returns a []byte and true if found or []byte{} and false if not
func (c *Cache) Get(name string) ([]byte, bool) {

	return []byte{}, false
}

// removes expired cache entries
//
// called when the cache is created by NewCache
// each time an interval (time.Duration) passes/occurs
// it should remove any entries that are expired
// (older than the interval)
func (c *Cache) readLoop(interval time.Duration) {

}
