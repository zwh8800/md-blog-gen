package util

import "sync"

type AutoClearKVCache struct {
	timesBeforeClear int
	n                int
	data             map[string]string
	mut              sync.RWMutex
}

func NewAutoClearKVCache(timesBeforeClear int) *AutoClearKVCache {
	return &AutoClearKVCache{
		timesBeforeClear: timesBeforeClear,
		n:                0,
		data:             make(map[string]string),
	}
}

func (c *AutoClearKVCache) add1s() {
	c.mut.Lock()
	defer c.mut.Unlock()
	c.n++
	if c.n >= c.timesBeforeClear {
		c.data = make(map[string]string)
		c.n = 0
	}
}

func (c *AutoClearKVCache) lockAndGet(key string) (string, bool) {
	c.mut.RLock()
	defer c.mut.RUnlock()
	ret, ok := c.data[key]
	return ret, ok

}

func (c *AutoClearKVCache) Get(key string) (string, bool) {
	ret, ok := c.lockAndGet(key)
	c.add1s()
	return ret, ok
}

func (c *AutoClearKVCache) Put(key, value string) {
	c.mut.Lock()
	defer c.mut.Unlock()
	c.data[key] = value
}
