package collections

import (
	"container/list"
	"sync"
)

type SafeOrderedDict[T any] struct {
	sync.RWMutex
	lookup    map[string]*list.Element
	list      *list.List
	listCache []T
}

func NewSafeOrderedDict[T any]() *SafeOrderedDict[T] {
	return &SafeOrderedDict[T]{
		lookup:    make(map[string]*list.Element),
		list:      list.New(),
		listCache: make([]T, 0, 0),
	}
}

func (o *SafeOrderedDict[T]) Set(key string, value T) {
	o.Lock()
	if n, ok := o.lookup[key]; ok {
		o.list.Remove(n)
	}

	o.lookup[key] = o.list.PushBack(value)
	o.Unlock()
	go o.updateListCache()
}

func (o *SafeOrderedDict[T]) Get(key string) (T, bool) {
	o.RLock()
	defer o.RUnlock()

	n, ok := o.lookup[key]
	var val T
	if n != nil {
		val = n.Value()
	}
	return val, ok
}

func (o *SafeOrderedDict[T]) Remove(key string) bool {
	o.Lock()
	defer o.Unlock()

	if n, ok := o.lookup[key]; ok {
		delete(o.lookup, key)
		if removed := o.list.Remove(n); removed == nil {
			return false
		}
		go o.updateListCache()

		return true
	}

	return false
}

func (o *SafeOrderedDict[T]) Iterate() chan T {
	o.RLock()
	ch := make(chan T)

	go func() {
		defer close(ch)
		defer o.RUnlock()

		for e := o.list.Front(); e != nil; e = e.Next() {
			ch <- e.Value
		}

	}()

	return ch
}

func (o *SafeOrderedDict[T]) List() []T {
	o.RLock()
	defer o.RUnlock()

	return o.listCache
}

func (o *SafeOrderedDict[T]) Count() int {
	return len(o.lookup)
}

func (o *SafeOrderedDict[T]) updateListCache() {
	o.RLock()
	defer o.RUnlock()

	index := 0
	o.listCache = make([]T, o.Count(), o.Count())
	for element := range o.Iterate() {
		o.listCache[index] = element
		index += 1
	}
}
