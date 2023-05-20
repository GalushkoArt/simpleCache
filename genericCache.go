package simpleCache

import (
	"sync"
	"time"
)

type GenericCache[T interface{}] interface {
	Set(key string, value T)
	Get(key string) *T
	Delete(key string) *T
}

type ExpiryGenericCache[T interface{}] interface {
	GenericCache[T]
	SetWithExpiry(key string, nano int64, value T)
}

type genericMapCache[T interface{}] struct {
	storage map[string]interface{}
}

func NewGenericMapCache[T interface{}]() GenericCache[T] {
	return &genericMapCache[T]{storage: make(map[string]interface{})}
}

func (c *genericMapCache[T]) Set(key string, value T) {
	c.storage[key] = value
}

func (c *genericMapCache[T]) Get(key string) *T {
	value, ok := c.storage[key].(T)
	if !ok {
		return nil
	}
	return &value
}

func (c *genericMapCache[T]) Delete(key string) *T {
	currentValue := c.Get(key)
	if currentValue == nil {
		return nil
	}
	delete(c.storage, key)
	return currentValue
}

type concurrencyGenericMapCache[T interface{}] struct {
	concurrencyMapCache
}

// NewGenericConcurrentCache returns parametrized concurrency cache with provided TTL
// if ttl is zero cache won't be purified automatically
func NewGenericConcurrentCache[T interface{}](ttl time.Duration) ExpiryGenericCache[T] {
	c := &concurrencyGenericMapCache[T]{
		concurrencyMapCache{
			storage: make(map[string]itemWithExpiration),
			ttl:     ttl,
			mutex:   sync.RWMutex{},
		},
	}
	addPurifier(&c.concurrencyMapCache)
	return c
}

func (c *concurrencyGenericMapCache[T]) Set(key string, value T) {
	c.SetWithExpiry(key, time.Now().Add(c.ttl).UnixNano(), value)
}

func (c *concurrencyGenericMapCache[T]) SetWithExpiry(key string, nano int64, value T) {
	if nano < time.Now().UnixNano() {
		return
	}
	c.mutex.Lock()
	if c.ttl > 0 {
		c.storage[key] = itemWithExpiration{
			item:      value,
			expiredAt: nano,
		}
	} else {
		c.storage[key] = itemWithExpiration{item: value}
	}
	c.mutex.Unlock()
}

func (c *concurrencyGenericMapCache[T]) Get(key string) *T {
	c.mutex.RLock()
	value, ok := c.storage[key]
	c.mutex.RUnlock()
	if !ok || value.expiredAt < time.Now().UnixNano() {
		return nil
	}
	result := value.item.(T)
	return &result
}

func (c *concurrencyGenericMapCache[T]) Delete(key string) *T {
	currentValue := c.Get(key)
	if currentValue == nil {
		return nil
	}
	c.mutex.Lock()
	delete(c.storage, key)
	c.mutex.Unlock()
	return currentValue
}
