package use

type Collection[Key comparable, Value any] struct {
	values     map[Key]Value
	order      []Key
	currentKey int
}

func NewCollection[Key comparable, Value any]() *Collection[Key, Value] {
	return &Collection[Key, Value]{values: make(map[Key]Value), order: make([]Key, 0), currentKey: -1}
}

func (c *Collection[Key, Value]) Add(key Key, value Value) {
	c.values[key] = value
	c.order = append(c.order, key)
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

func (c *Collection[Key, Value]) Slice() []Value {
	v := make([]Value, 0, c.Count())

	for _, value := range c.values {
		v = append(v, value)
	}

	return v
}

func (c *Collection[Key, Value]) Map() map[Key]Value {
	return c.values
}

func (c *Collection[Key, Value]) Next() Value {
	c.currentKey++

	return c.values[c.Key()]
}

func (c *Collection[Key, Value]) First() Value {
	c.currentKey = -1

	return c.Next()
}

func (c *Collection[Key, Value]) Current() Value {
	if c.currentKey < 0 {
		return c.First()
	}

	return c.values[c.Key()]
}

func (c *Collection[Key, Value]) Key() Key {
	mykey := c.currentKey
	if mykey < 0 {
		mykey = 0
	}

	return c.order[mykey]
}

func (c *Collection[Key, Value]) Last() Value {
	c.currentKey = len(c.order) - 1

	return c.Current()
}

func (c *Collection[Key, Value]) HasNext() bool {
	return c.currentKey < len(c.order)-1
}

func (c *Collection[Key, Value]) HasPrevious() bool {
	return c.currentKey > 0
}

func (c *Collection[Key, Value]) Previous() interface{} {
	c.currentKey--

	return c.values[c.Key()]
}

func (c *Collection[Key, Value]) Merge(other *Collection[Key, Value]) {
	other.Iterate(func(key Key, value Value) {
		c.Add(key, value)
	})
}
