package eviction

import (
	"container/list"
	"sync"
)

type LRUNode struct {
	Key   string
	Value interface{}
}

type LeastRecentlyUsed struct {
	capacity  int
	items     map[string]*list.Element
	evictList *list.List
	mutex     sync.RWMutex
}

func NewLRU(capacity int) *LeastRecentlyUsed {
	return &LeastRecentlyUsed{
		capacity:  capacity,
		items:     make(map[string]*list.Element),
		evictList: list.New(),
	}
}

func (lru *LeastRecentlyUsed) Get(key string) (interface{}, bool) {
	lru.mutex.Lock()
	defer lru.mutex.Unlock()

	if element, exists := lru.items[key]; exists {
		lru.evictList.MoveToFront(element)
		return element.Value.(*LRUNode).Value, true
	}
	return nil, false
}

func (lru *LeastRecentlyUsed) Put(key string, value interface{}) {
	lru.mutex.Lock()
	defer lru.mutex.Unlock()

	if element, exists := lru.items[key]; exists {
		lru.evictList.MoveToFront(element)
		element.Value.(*LRUNode).Value = value
		return
	}
	if lru.evictList.Len() >= lru.capacity {
		lru.evict()
	}

	element := lru.evictList.PushFront(&LRUNode{Key: key, Value: value})
	lru.items[key] = element
}

func (lru *LeastRecentlyUsed) evict() {
	if element := lru.evictList.Back(); element != nil {
		lru.removeElement(element)
	}
}

func (lru *LeastRecentlyUsed) Remove(key string) bool {
	lru.mutex.Lock()
	defer lru.mutex.Unlock()

	if element, exists := lru.items[key]; exists {
		lru.removeElement(element)
		return true
	}
	return false
}

func (lru *LeastRecentlyUsed) removeElement(element *list.Element) {
	lru.evictList.Remove(element)
	node := element.Value.(*LRUNode)
	delete(lru.items, node.Key)
}

func (lru *LeastRecentlyUsed) Len() int {
	lru.mutex.RLock()
	defer lru.mutex.RUnlock()
	return lru.evictList.Len()
}

func (lru *LeastRecentlyUsed) Clear() {
	lru.mutex.Lock()
	defer lru.mutex.Unlock()

	lru.items = make(map[string]*list.Element)
	lru.evictList.Init()
}
