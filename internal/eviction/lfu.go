package eviction

import (
	"container/list"
	"sync"
)

type LFUNode struct {
	Key       string
	Value     interface{}
	Frequency int
}

type FrequencyNode struct {
	Frequency int
	Items     *list.List
}

type LeastFrequentlyUsed struct {
	capacity    int
	items       map[string]*list.Element
	frequencies *list.List
	freqMap     map[int]*list.Element
	mutex       sync.Mutex
}

func NewLFU(capacity int) *LeastFrequentlyUsed {
	lfu := &LeastFrequentlyUsed{
		capacity:    capacity,
		items:       make(map[string]*list.Element),
		frequencies: list.New(),
		freqMap:     make(map[int]*list.Element),
	}

	freqNode := &FrequencyNode{
		Frequency: 0,
		Items:     list.New(),
	}
	freqElem := lfu.frequencies.PushFront(freqNode)
	lfu.freqMap[0] = freqElem

	return lfu
}

func (lfu *LeastFrequentlyUsed) Get(key string) (interface{}, bool) {
	lfu.mutex.Lock()
	defer lfu.mutex.Unlock()

	if element, exists := lfu.items[key]; exists {
		node := element.Value.(*LFUNode)
		freqElement := lfu.freqMap[node.Frequency]
		freqNode := freqElement.Value.(*FrequencyNode)

		freqNode.Items.Remove(element)

		if freqNode.Items.Len() == 0 && node.Frequency != 0 {
			lfu.frequencies.Remove(freqElement)
			delete(lfu.freqMap, node.Frequency)
		}

		node.Frequency++

		nextFreqElement, exists := lfu.freqMap[node.Frequency]
		if !exists {
			nextFreqElement = lfu.frequencies.InsertAfter(&FrequencyNode{
				Frequency: node.Frequency,
				Items:     list.New(),
			}, freqElement)
			lfu.freqMap[node.Frequency] = nextFreqElement
		}

		nextFreqNode := nextFreqElement.Value.(*FrequencyNode)
		lfu.items[key] = nextFreqNode.Items.PushFront(node)

		return node.Value, true
	}

	return nil, false
}

func (lfu *LeastFrequentlyUsed) Put(key string, value interface{}) {
	lfu.mutex.Lock()
	defer lfu.mutex.Unlock()

	if element, exists := lfu.items[key]; exists {
		node := element.Value.(*LFUNode)
		node.Value = value
		lfu.Get(key)
		return
	}

	if len(lfu.items) >= lfu.capacity {
		lfu.evict()
	}

	node := &LFUNode{
		Key:       key,
		Value:     value,
		Frequency: 0,
	}

	freqElement := lfu.freqMap[0]
	freqNode := freqElement.Value.(*FrequencyNode)
	lfu.items[key] = freqNode.Items.PushFront(node)
}

func (lfu *LeastFrequentlyUsed) Remove(key string) bool {
	lfu.mutex.Lock()
	defer lfu.mutex.Unlock()

	element, exists := lfu.items[key]
	if !exists {
		return false
	}

	node := element.Value.(*LFUNode)
	freqElement := lfu.freqMap[node.Frequency]
	freqNode := freqElement.Value.(*FrequencyNode)

	freqNode.Items.Remove(element)

	if freqNode.Items.Len() == 0 && node.Frequency != 0 {
		lfu.frequencies.Remove(freqElement)
		delete(lfu.freqMap, node.Frequency)
	}

	delete(lfu.items, key)

	return true
}

func (lfu *LeastFrequentlyUsed) evict() {
	freqElement := lfu.frequencies.Front()
	freqNode := freqElement.Value.(*FrequencyNode)

	element := freqNode.Items.Back()
	if element == nil {
		return
	}

	node := element.Value.(*LFUNode)

	delete(lfu.items, node.Key)

	freqNode.Items.Remove(element)

	if freqNode.Items.Len() == 0 && freqNode.Frequency != 0 {
		lfu.frequencies.Remove(freqElement)
		delete(lfu.freqMap, freqNode.Frequency)
	}
}

func (lfu *LeastFrequentlyUsed) Len() int {
	lfu.mutex.Lock()
	defer lfu.mutex.Unlock()
	return len(lfu.items)
}

func (lfu *LeastFrequentlyUsed) Clear() {
	lfu.mutex.Lock()
	defer lfu.mutex.Unlock()

	lfu.items = make(map[string]*list.Element)
	lfu.frequencies = list.New()
	lfu.freqMap = make(map[int]*list.Element)

	freqNode := &FrequencyNode{
		Frequency: 0,
		Items:     list.New(),
	}
	freqElem := lfu.frequencies.PushFront(freqNode)
	lfu.freqMap[0] = freqElem
}
