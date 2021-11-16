package cache

import (
	"sync"
	"time"
)

type Config struct {
	Expire   time.Duration
	GcPeriod time.Duration
}

type Cache struct {
	*cache
	// If this is confusing, see the comment at the bottom of New()
}

type item struct {
	Object   interface{}
	added    time.Time
	duration time.Duration
}

func (item item) Expired() bool {
	if item.duration < 0 {
		return false
	}
	return time.Now().Sub(item.added) > item.duration
}

type cache struct {
	config   Config
	items    map[string]item
	mu       sync.RWMutex
	onExpire func(string, interface{})
	gc       *gc
}

// Add an item to the cache, replacing any existing item.
// If the duration is 0 (DefaultExpiration), the cache's default expiration time is used.
// If it is -1 (NoExpire), the item never expires.
func (c *cache) Set(k string, x interface{}, dur time.Duration) {
	if x == nil {
		return
	}
	it := item{
		Object:   x,
		added:    time.Now(),
		duration: NoExpire,
	}
	if dur == DefaultExpire {
		it.duration = c.config.Expire
	} else if dur > 0 {
		it.duration = dur
	}
	c.mu.Lock()
	c.items[k] = it
	// Calls to mu.Unlock are currently not deferred because defer
	// adds ~200 ns (as of go1.)
	c.mu.Unlock()
}

// if item exist in cache cache, replacing any existing item.
// If the duration is 0 (DefaultExpiration), the cache's default expiration time is used.
// If it is -1 (NoExpiration), the item never expires.
func (c *cache) Touch(key string, dur time.Duration) (touched bool) {
	c.mu.Lock()
	if item, ok := c.items[key]; ok {
		item.added = time.Now()
		if dur == DefaultExpire {
			item.duration = c.config.Expire
		} else if dur > 0 {
			item.duration = dur
		}
		c.items[key] = item
		touched = true
	}
	c.mu.Unlock()
	return
}

// Remove an item from the cache.
// Returns the item or nil,// and a bool indicating if the key was found and deleted.
func (c *cache) Remove(key string) (interface{}, bool) {
	c.mu.Lock()
	o, ok := c.items[key]
	if ok {
		delete(c.items, key)
	}
	c.mu.Unlock()
	return o, ok
}

// Get an item from the cache.
// Returns the item or nil,
// and a bool indicating if the key was found.
func (c *cache) Get(k string) (interface{}, bool) {
	c.mu.RLock()
	item, ok := c.items[k]
	c.mu.RUnlock()
	if !ok || item.Expired() {
		return nil, false
	}
	return item.Object, true
}

// Returns the number of items in the cache. This may include items that have
// expired, but have not yet been cleaned up.
func (c *cache) ItemCount() int {
	c.mu.RLock()
	n := len(c.items)
	c.mu.RUnlock()
	return n
}

// Delete all items from the cache.
func (c *cache) Flush() {
	c.mu.Lock()
	c.items = make(map[string]item)
	c.mu.Unlock()
}

func (c *cache) FlushExpired() {
	expiredItems := make(map[string]interface{})
	c.mu.RLock()
	for key, item := range c.items {
		if item.Expired() {
			expiredItems[key] = item.Object
		}
	}
	c.mu.RUnlock()
	go func(removed map[string]interface{}) {
		for key, obj := range removed {
			c.Remove(key)
			if c.onExpire != nil {
				go c.onExpire(key, obj)
			}
		}
	}(expiredItems)
}

// Sets an (optional) function that is called with the key and value when an
// item is evicted from the cache by expire.
func (c *cache) OnExpire(f func(string, interface{})) {
	c.mu.Lock()
	c.onExpire = f
	c.mu.Unlock()
}

type gc struct {
	Interval time.Duration
	stop     chan bool
}

func (j *gc) Run(c *cache) {
	ticker := time.NewTicker(j.Interval)
	for {
		select {
		case <-ticker.C:
			c.FlushExpired()
		case <-j.stop:
			ticker.Stop()
			return
		}
	}
}
