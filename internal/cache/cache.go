package cache

type store map[interface{}]interface{}

type Cache struct {
	store store
}

func NewCache() *Cache {
	return &Cache{
		store: store{},
	}
}

func (c *Cache) Set(key, val interface{}) {
	c.store[key] = val
}

func (c *Cache) Get(key interface{}) (val interface{}, ok bool) {
	val, ok = c.store[key]
	return
}

func (c *Cache) Delete(key interface{}) {
	delete(c.store, key)
}
