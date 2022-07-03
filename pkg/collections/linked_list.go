package collections

type Node[T any] struct {
	next *Node[T]
	prev *Node[T]

	value T
}

func (n *Node[T]) Value() T {
	return n.value
}

type LinkedList[T any] struct {
	head *Node[T]
	tail *Node[T]
}

func NewLinkedList[T any]() *LinkedList[T] {
	var (
		head T
		tail T
	)
	list := &LinkedList[T]{
		head: &Node[T]{value: head},
		tail: &Node[T]{value: tail},
	}
	list.head.next = list.tail
	list.tail.prev = list.head
	return list
}

func (l *LinkedList[T]) Append(value T) *Node[T] {
	n := &Node[T]{
		prev:  l.tail.prev,
		next:  l.tail,
		value: value,
	}

	l.tail.prev.next = n
	l.tail.prev = n

	return n
}

func (l *LinkedList[T]) Remove(n *Node[T]) bool {
	if n == l.head || n == l.tail {
		return false
	}

	n.prev.next = n.next
	n.next.prev = n.prev

	return true
}

func (l *LinkedList[T]) Iterate() chan *Node[T] {
	ch := make(chan *Node[T])

	go func() {
		defer close(ch)
		n := l.head

		for n.next != l.tail {
			ch <- n.next
			n = n.next
		}
	}()

	return ch
}
