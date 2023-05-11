package simpleCache

type SimpleCache interface {
	Set(key string, value interface{})
	Get(key string) *interface{}
	Delete(key string) *interface{}
}

type SimpleMapCache struct {
	storage map[string]interface{}
}

func NewMapCache() SimpleCache {
	return &SimpleMapCache{storage: make(map[string]interface{})}
}

func (c *SimpleMapCache) Set(key string, value interface{}) {
	c.storage[key] = value
}

func (c *SimpleMapCache) Get(key string) *interface{} {
	value, ok := c.storage[key]
	if !ok {
		return nil
	}
	return &value
}

func (c *SimpleMapCache) Delete(key string) *interface{} {
	currentValue := c.Get(key)
	if currentValue == nil {
		return nil
	}
	delete(c.storage, key)
	return currentValue
}
