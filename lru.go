package lru

import (
	"container/list"
	"sync"
)

type value struct {
	id    string
	value interface{}
}

type Cache interface {
	Set(id string, e interface{})
	Get(id string) (interface{}, bool)
	Remove(id string) bool
}

type lru struct {
	list *list.List

	maxElementsSize int
	elementsByID    map[string]*list.Element

	mu *sync.Mutex
}

func NewLRU(size int) Cache {
	return &lru{
		list: list.New(),
		mu:   &sync.Mutex{},

		maxElementsSize: size,
		elementsByID:    make(map[string]*list.Element),
	}
}

func (l lru) Set(id string, v interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	valueToAdd := value{
		id:    id,
		value: v,
	}

	if e, ok := l.elementsByID[id]; ok {
		e.Value = valueToAdd
		l.list.MoveToBack(e)
		return
	}

	elementAdded := l.list.PushBack(valueToAdd)

	l.elementsByID[id] = elementAdded

	if len(l.elementsByID) > l.maxElementsSize {
		// remove the oldest element from the list
		firstElement := l.list.Front()
		firstValue := firstElement.Value.(value)

		l.remove(firstValue.id)
	}
}

func (l lru) Get(id string) (interface{}, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()

	e, ok := l.elementsByID[id]
	if !ok {
		return nil, false
	}

	// move element to back as recently used
	l.list.MoveToBack(e)

	v := e.Value.(value)
	return v.value, true
}

func (l lru) remove(id string) bool {
	e, ok := l.elementsByID[id]

	if !ok {
		return false
	}

	delete(l.elementsByID, id)
	l.list.Remove(e)
	return true

}

func (l lru) Remove(id string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.remove(id)
}
