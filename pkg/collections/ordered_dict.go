package collections

import "sync"

type OrderedDict[T any] struct {
	sync.RWMutex
	lookup map[string]*Node[T]
	list   *LinkedList[T]
}

func NewOrderedDict[T any]() *OrderedDict[T] {
	return &OrderedDict[T]{
		lookup: make(map[string]*Node[T]),
		list:   NewLinkedList[T](),
	}
}

func (o *OrderedDict[T]) Set(key string, value T) {
	o.Lock()
	defer o.Unlock()

	if n, ok := o.lookup[key]; ok {
		o.list.Remove(n)
	}
	o.lookup[key] = o.list.Append(value)
}

func (o *OrderedDict[T]) Get(key string) (T, bool) {
	o.RLock()
	defer o.RUnlock()

	n, ok := o.lookup[key]
	var val T
	if n != nil {
		val = n.Value()
	}
	return val, ok
}

func (o *OrderedDict[T]) Remove(key string) bool {
	o.Lock()
	defer o.Unlock()

	if n, ok := o.lookup[key]; ok {
		if ok := o.list.Remove(n); !ok {
			return false
		}
		delete(o.lookup, key)
		return true
	}

	return false
}

func (o *OrderedDict[T]) Iterate() chan T {
	o.RLock()
	ch := make(chan T)

	go func() {
		defer close(ch)
		defer o.RUnlock()

		for v := range o.list.Iterate() {
			ch <- v.Value()
		}

	}()

	return ch
}

func (o *OrderedDict[T]) Count() int {
	return len(o.lookup)
}
