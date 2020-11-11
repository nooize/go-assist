package cache

func stopGc(c *Cache) {
	c.gc.stop <- true
}

func runGc(c *cache) {
	j := &gc{
		Interval: c.config.GcPeriod,
		stop:     make(chan bool),
	}
	c.gc = j
	go j.Run(c)
}
