package use

type Collection[Key comparable, Value any] struct {
	values map[Key]Value
}

func NewCollection[Key comparable, Value any]() *Collection[Key, Value] {
	return &Collection[Key, Value]{values: make(map[Key]Value)}
}

func (c *Collection[Key, Value]) Add(key Key, value Value) {
	c.values[key] = value
}

func (c *Collection[Key, Value]) Get(key Key) Value {
	return c.values[key]
}

func (c *Collection[Key, Value]) Set(data map[Key]Value) {
	c.values = data
}

func (c *Collection[Key, Value]) Has(key Key) bool {
	_, isPresent := c.values[key]

	return isPresent
}

func (c *Collection[Key, Value]) Iterate(fn func(key Key, value Value)) {
	for key, value := range c.values {
		fn(key, value)
	}
}

func (c *Collection[Key, Value]) Remove(key Key) bool {
	if c.Has(key) {
		delete(c.values, key)
		return true
	}

	return false
}

func (c *Collection[Key, Value]) RemoveAll() {
	c.values = make(map[Key]Value)
}

func (c *Collection[Key, Value]) Count() int {
	return len(c.values)
}
