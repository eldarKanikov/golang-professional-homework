package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
	mutex    sync.Mutex
}

func (lru *lruCache) Clear() {
	lru.mutex.Lock()
	defer lru.mutex.Unlock()

	lru.queue = NewList()
	lru.items = make(map[Key]*ListItem, lru.capacity)
}

func (lru *lruCache) Get(key Key) (interface{}, bool) {
	lru.mutex.Lock()
	defer lru.mutex.Unlock()
	item, exist := lru.items[key]
	if !exist {
		return nil, false
	}
	lru.queue.MoveToFront(item)
	value := item.Value
	return value.(*Pair).el2, true
}

func (lru *lruCache) Set(key Key, value interface{}) bool {
	lru.mutex.Lock()
	defer lru.mutex.Unlock()
	item, exist := lru.items[key]
	if exist {
		pair := item.Value.(*Pair)
		pair.el2 = value
		lru.queue.MoveToFront(item)
	} else {
		lru.queue.PushFront(newPair(key, value))
		lru.items[key] = lru.queue.Front()

		if lru.queue.Len() > lru.capacity {
			back := lru.queue.Back()
			delete(lru.items, back.Value.(*Pair).el1)
			lru.queue.Remove(back)
		}
	}
	return exist
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

type Pair struct {
	el1 Key
	el2 interface{}
}

func newPair(el1 Key, el2 interface{}) *Pair {
	return &Pair{
		el1: el1,
		el2: el2,
	}
}
