package simpleCache

import (
	"runtime"
	"sync"
	"time"
)

type SimpleCache interface {
	Set(key string, value interface{})
	Get(key string) *interface{}
	Delete(key string) *interface{}
}

type simpleMapCache struct {
	storage map[string]interface{}
}

func NewMapCache() SimpleCache {
	return &simpleMapCache{storage: make(map[string]interface{})}
}

func (c *simpleMapCache) Set(key string, value interface{}) {
	c.storage[key] = value
}

func (c *simpleMapCache) Get(key string) *interface{} {
	value, ok := c.storage[key]
	if !ok {
		return nil
	}
	return &value
}

func (c *simpleMapCache) Delete(key string) *interface{} {
	currentValue := c.Get(key)
	if currentValue == nil {
		return nil
	}
	delete(c.storage, key)
	return currentValue
}

type itemWithExpiration struct {
	item      interface{}
	expiredAt int64
}

type purifier struct {
	done chan bool
}

type concurrencyMapCache struct {
	storage  map[string]itemWithExpiration
	ttl      time.Duration
	mutex    sync.RWMutex
	purifier *purifier
}

func addPurifier(c *concurrencyMapCache) {
	c.purifier = &purifier{done: make(chan bool)}
	if c.ttl > 0 {
		go c.purifier.runPurifier(c)
		runtime.SetFinalizer(c.purifier, stopPurifier)
	}
}

func (p *purifier) runPurifier(c *concurrencyMapCache) {
	ticker := time.NewTicker(c.ttl)
	for {
		select {
		case <-ticker.C:
			c.purify()
		case <-p.done:
			ticker.Stop()
			return
		}
	}
}

func (c *concurrencyMapCache) purify() {
	c.mutex.Lock()
	now := time.Now().UnixNano()
	for key, item := range c.storage {
		if item.expiredAt < now {
			delete(c.storage, key)
		}
	}
	c.mutex.Unlock()
}

func stopPurifier(p *purifier) {
	p.done <- true
}

// NewConcurrentCache returns concurrency cache with provided TTL
// if ttl is zero cache won't be purified automatically
func NewConcurrentCache(ttl time.Duration) SimpleCache {
	c := &concurrencyMapCache{
		storage: make(map[string]itemWithExpiration),
		ttl:     ttl,
		mutex:   sync.RWMutex{},
	}
	addPurifier(c)
	return c
}

func (c *concurrencyMapCache) Set(key string, value interface{}) {
	c.mutex.Lock()
	if c.ttl > 0 {
		c.storage[key] = itemWithExpiration{
			item:      value,
			expiredAt: time.Now().Add(c.ttl).UnixNano(),
		}
	} else {
		c.storage[key] = itemWithExpiration{item: value}
	}
	c.mutex.Unlock()
}

func (c *concurrencyMapCache) Get(key string) *interface{} {
	c.mutex.RLock()
	value, ok := c.storage[key]
	c.mutex.RUnlock()
	if !ok || value.expiredAt < time.Now().UnixNano() {
		return nil
	}
	return &(value.item)
}

func (c *concurrencyMapCache) Delete(key string) *interface{} {
	currentValue := c.Get(key)
	if currentValue == nil {
		return nil
	}
	c.mutex.Lock()
	delete(c.storage, key)
	c.mutex.Unlock()
	return currentValue
}
