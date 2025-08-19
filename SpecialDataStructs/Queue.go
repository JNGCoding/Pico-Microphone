package DataStructures

type Queue[T any] struct {
	items    []T
	head     int
	tail     int
	size     int
	capacity int
}

func CreateQueue[T any](capacity int) *Queue[T] {
	return &Queue[T]{
		items:    make([]T, capacity),
		capacity: capacity,
	}
}

func (q *Queue[T]) Enqueue(item T) bool {
	if q.size == q.capacity {
		return false // Queue is full
	}
	q.items[q.tail] = item
	q.tail = (q.tail + 1) % q.capacity
	q.size++
	return true
}

func (q *Queue[T]) Dequeue() (T, bool) {
	var zero T
	if q.size == 0 {
		return zero, false
	}
	item := q.items[q.head]
	q.head = (q.head + 1) % q.capacity
	q.size--
	return item, true
}

func (q *Queue[T]) Peek() (T, bool) {
	var zero T
	if q.size == 0 {
		return zero, false
	}
	return q.items[q.head], true
}

func (q *Queue[T]) IsEmpty() bool {
	return q.size == 0
}

func (q *Queue[T]) IsFull() bool {
	return q.size == q.capacity
}

func (q *Queue[T]) Size() int {
	return q.size
}

func (q *Queue[T]) Clear() {
	q.head = 0
	q.tail = 0
	q.size = 0

	var zero T
	for i := range q.items {
		q.items[i] = zero
	}
}
