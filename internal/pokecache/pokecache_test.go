package pokecache_test

import (
	"fmt"
	"testing"
	"time"

	pokecache "github.com/nicholasss/pokedexcli/internal/pokecache"
)

func TestAddGet(t *testing.T) {
	const interval = (5 * time.Second)

	cases := []struct {
		key string
		val []byte
	}{
		{
			key: "https://example.com",
			val: []byte("thisisanexample"),
		},
		{
			key: "https://example.com/testpath",
			val: []byte("noodleisagoodcat"),
		},
		{
			key: "https://api.example.com/v1/hello",
			val: []byte("{\"hello\": \"world\"}"),
		},
		{
			key: "https://api.example.com/v2/hello",
			val: []byte("{\"hello\": \"second world\"}"),
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T) {
			cache := pokecache.NewCache(interval)
			cache.Add(c.key, c.val)

			val, ok := cache.Get(c.key)
			if !ok {
				t.Errorf("expected to find key")
				return
			}

			if string(val) != string(c.val) {
				t.Errorf("expected to find value")
				return
			}

		})
	}
}

func TestReadLoop(t *testing.T) {
	const baseTime = (5 * time.Millisecond)
	const waitTime = (baseTime + (5 * time.Millisecond))

	cache := pokecache.NewCache(baseTime)
	URL := "https://example.com"
	cache.Add(URL, []byte("testData"))

	_, ok := cache.Get(URL)
	if !ok {
		t.Errorf("expected to find key")
		return
	}

	time.Sleep(waitTime)

	_, ok = cache.Get(URL)
	if ok {
		t.Errorf("expected not to find key")
		return
	}
}
