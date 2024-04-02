package lru

import (
	"container/list"
	"errors"
)

type Cache struct {
	// 字典的键为描述，值为队列的元素
	cache map[string]*list.Element
	// 一个队列，队尾为最近访问的元素
	ll *list.List

	// Cache的最大容量
	maxBytes int64
	// 当前容量
	nbytes int64

	// 回调函数
	OnEvicted func(key string, value Value)
}

type entry struct {
	key   string
	value Value
}

// Value 是一个接口，根据不同的元素内容实现
type Value interface {
	Len() int
}

// 初始化
func NewCache(maxByte int64, onEvicted func(string, Value)) *Cache {
	c := &Cache{
		cache: map[string]*list.Element{},
		ll:    list.New(),

		maxBytes: maxByte,

		OnEvicted: onEvicted,
	}
	return c
}

// 查找元素
func (c *Cache) Get(key string) (value Value, ok bool) {
	if element, ok := c.cache[key]; ok {
		c.ll.MoveToFront(element)
		kv := element.Value.(*entry)
		return kv.value, ok
	}
	return
}

// 移除元素
func (c *Cache) RemoveOldest() error {
	element := c.ll.Back()
	// 没有元素可以移除
	if element == nil {
		return errors.New("nothing can be removed")
	}
	// 移除list末尾元素
	c.ll.Remove(element)
	// 删除map中的键值对
	kv := element.Value.(*entry)
	delete(c.cache, kv.key)
	// 修改当前容量
	c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
	// 执行回调函数
	if c.OnEvicted != nil {
		c.OnEvicted(kv.key, kv.value)
	}
	return nil
}

// 增加内容，若内容已存在则删除
func (c *Cache) Add(key string, value Value) {
	// 修改
	if element, ok := c.cache[key]; ok {
		c.ll.MoveToFront(element)
		kv := element.Value.(*entry) // 旧元素
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else { // 增加
		element := c.ll.PushFront(&entry{key: key, value: value})
		c.cache[key] = element
		c.nbytes += int64(len(key)) + int64(value.Len())
	}
	for c.maxBytes > 0 && c.maxBytes < c.nbytes {
		err := c.RemoveOldest()
		if err != nil {
			break
		}
	}
}

func (c *Cache) Len() int {
	return c.ll.Len()
}
