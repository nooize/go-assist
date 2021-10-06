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
	Object interface{}
	Expire int64
}

func (item item) Expired() bool {
	if item.Expire == 0 {
		return false
	}
	return time.Now().UnixNano() > item.Expire
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
// If it is -1 (NoExpiration), the item never expires.
func (c *cache) Set(k string, x interface{}, dur time.Duration) {
	exp := time.Now()
	if dur == DefaultExpire {
		exp = exp.Add(c.config.Expire)
	} else if dur > 0 {
		exp = exp.Add(dur)
	}
	c.mu.Lock()
	c.items[k] = item{
		Object: x,
		Expire: exp.UnixNano(),
	}
	// Calls to mu.Unlock are currently not deferred because defer
	// adds ~200 ns (as of go1.)
	c.mu.Unlock()
}

// Remove an item from the cache.
// Returns the item or nil,
// and a bool indicating if the key was found and deleted.
func (c *cache) Remove(k string) (interface{}, bool) {
	c.mu.RLock()
	o, ok := c.items[k]
	if ok {
		c.delete(k)
	}
	return o, ok
}

// Get an item from the cache.
// Returns the item or nil,
// and a bool indicating if the key was found.
func (c *cache) Get(k string) (interface{}, bool) {
	c.mu.RLock()
	// "Inlining" of get and Expired
	item, ok := c.items[k]
	if !ok {
		c.mu.RUnlock()
		return nil, false
	}
	if item.Expired() {
		c.mu.RUnlock()
		return nil, false
	}
	c.mu.RUnlock()
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
	now := time.Now().UnixNano()
	c.mu.Lock()
	for k, v := range c.items {
		if v.Expire > 0 && now > v.Expire {
			ov, evicted := c.delete(k)
			if evicted {
				expiredItems[k] = ov
			}
		}
	}
	c.mu.Unlock()
	for k, v := range expiredItems {
		c.onExpire(k, v)
	}
}

// Sets an (optional) function that is called with the key and value when an
// item is evicted from the cache by expire.
func (c *cache) OnExpire(f func(string, interface{})) {
	c.mu.Lock()
	c.onExpire = f
	c.mu.Unlock()
}

func (c *cache) delete(k string) (interface{}, bool) {
	delete(c.items, k)
	return nil, false
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
