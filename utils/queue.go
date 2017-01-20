package utils

import "sync"

type queuenode struct {
	data interface{}
	next *queuenode
}

// Queue go-routine safe FIFO (first in first out) data stucture.
type Queue struct {
	head  *queuenode
	tail  *queuenode
	count int
	lock  *sync.Mutex
}

// NewQueue creates a new pointer to a new queue.
func NewQueue() *Queue {
	q := &Queue{}
	q.lock = &sync.Mutex{}
	return q
}

// Len returns the number of elements in the queue (i.e. size/length)
func (q *Queue) Len() int {
	q.lock.Lock()
	defer q.lock.Unlock()
	return q.count
}

// Push pushes/inserts a value at the end/tail of the queue.
// This function mutates the queue.
func (q *Queue) Push(item interface{}) {
	q.lock.Lock()
	defer q.lock.Unlock()

	node := &queuenode{data: item}

	if q.tail == nil {
		q.tail = node
		q.head = node
	} else {
		q.tail.next = node
		q.tail = node
	}
	q.count++
}

// Poll returns the value at the front of the queue.
// This function mutates the queue.
func (q *Queue) Poll() interface{} {
	q.lock.Lock()
	defer q.lock.Unlock()

	if q.head == nil {
		return nil
	}

	node := q.head
	q.head = node.next

	if q.head == nil {
		q.tail = nil
	}
	q.count--

	return node.data
}

// Peek returns a read value at the front of the queue.
// This function does NOT mutate the queue.
func (q *Queue) Peek() interface{} {
	q.lock.Lock()
	defer q.lock.Unlock()

	node := q.head
	if node == nil {
		return nil
	}

	return node.data
}
