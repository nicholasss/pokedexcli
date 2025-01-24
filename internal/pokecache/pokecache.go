// pokecache is an internal package
// It provides a caching mechanism for URLs with a cacheEntry of []byte
package pokecache

import (
	"sync"
	"time"
)

// A Cache allows for storing and retrieving cached data associated with specific URLs
type Cache struct {
	entries map[string]cacheEntry
	mux     *sync.Mutex
}

// This type is responsible for storing raw []byte,
// as well as a time that it was created at.
type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

// Initializes a new Cache.
// Creates new cache that is valid for the given duration.
func NewCache(interval time.Duration) *Cache {
	newCache := Cache{
		entries: make(map[string]cacheEntry),
		mux:     &sync.Mutex{},
	}

	// begins the read loop that deletes old entires
	newCache.readLoop(interval)

	return &newCache
}

// Adds a new entry to the cache:
// Takes a key and a val to add to the Cache entries.
func (c *Cache) Add(name string, data []byte) {
	newCacheEntry := cacheEntry{
		createdAt: time.Now(),
		val:       data,
	}

	c.mux.Lock()
	defer c.mux.Unlock()

	// Will not replace item in map if it exists already.
	// Otherwise the data gotten from cache will just be used
	// to write over the existing data in cache every time
	// it is requested.
	if _, alreadyInCache := c.entries[name]; alreadyInCache {
		return
	} else {
		c.entries[name] = newCacheEntry
	}
}

// It takes a key
// and returns []byte and bool, true if there is an entry, false if there is none.
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
				expiAge := interval // ttl equal to interval

				if age > expiAge {
					delete(c.entries, name)
				}
			}

			// unlocks at end of timed loop
			c.mux.Unlock()
		}
	}()
}
