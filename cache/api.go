package cache

import (
	"runtime"
	"time"
)

const (
	// For use with functions that take an expiration time.
	NoExpire time.Duration = -1
	// For use with functions that take an expiration time. Equivalent to
	// passing in the same expiration duration as was given to New() or
	// when the cache was created (e.g. 5 minutes.)
	DefaultExpire time.Duration = 0

	DefaultCgPeriod = 5 * time.Minute
)

func New(conf *Config) *Cache {
	return newWithGC(conf, make(map[string]item))
}

func newWithGC(conf *Config, items map[string]item) *Cache {
	cfg := Config{
		Expire:   NoExpire,
		GcPeriod: DefaultCgPeriod,
	}
	if conf != nil && conf.Expire > 0 {
		cfg.Expire = conf.Expire
	}
	if conf != nil && conf.GcPeriod > 0 {
		cfg.GcPeriod = conf.GcPeriod
	}
	c := &cache{
		config: cfg,
		items:  items,
	}
	// This trick ensures that the janitor goroutine (which--granted it
	// was enabled--is running DeleteExpired on c forever) does not keep
	// the returned C object from being garbage collected. When it is
	// garbage collected, the finalizer stops the janitor goroutine, after
	// which c can be collected.
	C := &Cache{c}
	if cfg.GcPeriod > 0 {
		runGc(c)
		runtime.SetFinalizer(C, stopGc)
	}
	return C
}
