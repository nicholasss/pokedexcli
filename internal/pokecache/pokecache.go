package pokecache

import (
	"sync"
	"time"
)

// Cache
//
// this is the exposed struct that will be accessed
type Cache struct {
	entries map[string]cacheEntry
	mux     *sync.Mutex
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
	newCache := Cache{
		entries: make(map[string]cacheEntry),
		mux:     &sync.Mutex{},
	}

	// begins the read loop that deletes expired entires
	newCache.readLoop(interval)

	return &newCache
}

// adds a new entry to the cache
//
// take a key and a val to add to the Cache map
func (c *Cache) Add(name string, data []byte) {
	newCacheEntry := cacheEntry{
		createdAt: time.Now(),
		val:       data,
	}

	c.mux.Lock()
	defer c.mux.Unlock()
	c.entries[name] = newCacheEntry
}

// gets an entry to the cache
//
// takes a key
// returns a []byte and true if found or []byte{} and false if not
func (c *Cache) Get(name string) ([]byte, bool) {
	c.mux.Lock()
	defer c.mux.Unlock()

	foundEntry, ok := c.entries[name]
	if !ok {
		return []byte{}, false
	}

	return foundEntry.val, true
}

// removes expired cache entries
//
// called when the cache is created by NewCache
// each time an interval (time.Duration) passes/occurs
// it should remove any entries that are expired
// (older than the interval)
func (c *Cache) readLoop(interval time.Duration) {
	go func() {
		// create ticker
		ticker := time.NewTicker(interval)

		for {
			<-ticker.C

			// locks after timed loop
			c.mux.Lock()

			for name, entry := range c.entries {
				age := time.Since(entry.createdAt)
				expiAge := (time.Minute * 10) // 10min ttl

				if age > expiAge {
					delete(c.entries, name)
				}
			}

			// unlocks at end of timed loop
			c.mux.Unlock()
		}
	}()
}
