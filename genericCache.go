package simpleCache

type GenericCache[T interface{}] interface {
	Set(key string, value T)
	Get(key string) *T
	Delete(key string) *T
}

type GenericMapCache[T interface{}] struct {
	storage map[string]interface{}
}

func NewGenericMapCache[T interface{}]() GenericCache[T] {
	return &GenericMapCache[T]{storage: make(map[string]interface{})}
}

func (c *GenericMapCache[T]) Set(key string, value T) {
	c.storage[key] = value
}

func (c *GenericMapCache[T]) Get(key string) *T {
	value, ok := c.storage[key].(T)
	if !ok {
		return nil
	}
	return &value
}

func (c *GenericMapCache[T]) Delete(key string) *T {
	currentValue := c.Get(key)
	if currentValue == nil {
		return nil
	}
	delete(c.storage, key)
	return currentValue
}
