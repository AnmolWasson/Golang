package main

import "sync"

type BoundedBlockingQueue[T any] struct {
	queue         []T
	capacity      int
	lock          sync.Mutex
	isFullSignal  *sync.Cond
	isEmptySignal *sync.Cond
}

func Init[T any](capacity int) *BoundedBlockingQueue[T] {

	bbq := &BoundedBlockingQueue[T]{
		capacity: capacity,
		queue:    make([]T, 0),
		lock:     sync.Mutex{},
	}

	bbq.isEmptySignal = sync.NewCond(&bbq.lock)
	bbq.isFullSignal = sync.NewCond(&bbq.lock)

	return bbq
}

func (bbq *BoundedBlockingQueue[T]) Enqueue(item T) {
	bbq.lock.Lock()
	defer bbq.lock.Unlock()

	if len(bbq.queue) == bbq.capacity {
		bbq.isFullSignal.Wait() // waits for the signal that the queue is not full now
	}

	bbq.queue = append(bbq.queue, item)
	bbq.isEmptySignal.Signal() // sends the signal that the queue is not empty now

}

func (bbq *BoundedBlockingQueue[T]) Dequeue() T {
	bbq.lock.Lock()
	defer bbq.lock.Unlock()

	if len(bbq.queue) == 0 {
		bbq.isEmptySignal.Wait() // waits for the signal that the queue is not empty now
	}

	item := bbq.queue[0]
	bbq.queue = bbq.queue[1:]
	bbq.isFullSignal.Signal() // sends the signal that the queue is not full now

	return item

}
