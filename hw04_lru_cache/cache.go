package hw04lrucache

import (
	"sync"
)

type Key string

type val struct {
	A Key
	B interface{}
}

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	mu       sync.Mutex
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func NewCache(capacity int) Cache {
	if capacity < 0 {
		capacity = 0
	}

	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.capacity == 0 {
		return false
	}

	if item, ok := c.items[key]; ok {
		c.queue.MoveToFront(item)
		item.Value = val{key, value}
		return true
	}

	if c.capacity == c.queue.Len() {
		temp := c.queue.Back()
		delete(c.items, temp.Value.(val).A)
		c.queue.Remove(temp)
	}

	c.items[key] = c.queue.PushFront(val{key, value})

	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, ok := c.items[key]
	if !ok {
		return nil, false
	}

	c.queue.MoveToFront(item)

	return item.Value.(val).B, true
}

func (c *lruCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}
