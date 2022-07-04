package collections

import (
	"sync"
)

type SafeOrderedDict[T any] struct {
	sync.RWMutex
	lookup     map[string]*Node[T]
	linkedList *LinkedList[T]
	listCache  *SafeSlice[T]
}

func NewSafeOrderedDict[T any]() *SafeOrderedDict[T] {
	return &SafeOrderedDict[T]{
		lookup:     make(map[string]*Node[T]),
		linkedList: NewLinkedList[T](),
		listCache:  NewSafeSlice[T](0, 0),
	}
}

func (o *SafeOrderedDict[T]) Set(key string, value T) {
	o.Lock()
	if n, ok := o.lookup[key]; ok {
		o.linkedList.Remove(n)
	}

	o.lookup[key] = o.linkedList.Append(value)
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
		if ok := o.linkedList.Remove(n); !ok {
			return false
		}
		delete(o.lookup, key)
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

		for v := range o.linkedList.Iterate() {
			ch <- v.Value()
		}

	}()

	return ch
}

func (o *SafeOrderedDict[T]) List() []T {
	o.RLock()
	defer o.RUnlock()

	return o.listCache.Slice()
}

func (o *SafeOrderedDict[T]) Count() int {
	return len(o.lookup)
}

func (o *SafeOrderedDict[T]) updateListCache() {
	o.RLock()
	defer o.RUnlock()

	index := 0
	o.listCache = NewSafeSlice[T](o.Count(), o.Count())
	for element := range o.Iterate() {
		o.listCache.Set(index, element)
		index += 1
	}
}
