package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity    int
	valuesQueue List
	keysQueue   List // to store the consequence of keys
	valuesItems map[Key]*ListItem
	keysItems   map[Key]*ListItem // to have o(1) access
	mutex       sync.Mutex
}

func (lru *lruCache) Clear() {
	lru.mutex.Lock()
	defer lru.mutex.Unlock()
	lru.valuesQueue = NewList()
	lru.keysQueue = NewList()
	lru.valuesItems = make(map[Key]*ListItem, lru.capacity)
	lru.keysItems = make(map[Key]*ListItem, lru.capacity)
}

func (lru *lruCache) Get(key Key) (interface{}, bool) {
	lru.mutex.Lock()
	_, exist := lru.valuesItems[key]
	var result interface{}
	if exist {
		item := lru.valuesItems[key]
		result = item.Value
		lru.valuesQueue.MoveToFront(item)
		lru.keysQueue.MoveToFront(lru.keysItems[key])
	}
	defer lru.mutex.Unlock()
	return result, exist
}

func (lru *lruCache) Set(key Key, value interface{}) bool {
	lru.mutex.Lock()
	defer lru.mutex.Unlock()
	_, exist := lru.valuesItems[key]
	if exist {
		item := lru.valuesItems[key]
		item.Value = value
		lru.valuesQueue.MoveToFront(item)
		lru.keysQueue.MoveToFront(lru.keysItems[key])
	} else {
		lru.valuesQueue.PushFront(value)
		lru.keysQueue.PushFront(key)

		lru.valuesItems[key] = lru.valuesQueue.Front()
		lru.keysItems[key] = lru.keysQueue.Front()
		if lru.valuesQueue.Len() > lru.capacity {
			extraElementsCount := lru.valuesQueue.Len() - lru.capacity
			for i := 0; i < extraElementsCount; i++ {
				keyToDelete := lru.keysQueue.Back().Value.(Key)
				delete(lru.valuesItems, keyToDelete)
				delete(lru.keysItems, keyToDelete)
				lru.valuesQueue.Remove(lru.valuesQueue.Back())
				lru.keysQueue.Remove(lru.keysQueue.Back())
			}
		}
	}
	return exist
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity:    capacity,
		valuesQueue: NewList(),
		keysQueue:   NewList(),
		valuesItems: make(map[Key]*ListItem, capacity),
		keysItems:   make(map[Key]*ListItem, capacity),
	}
}
