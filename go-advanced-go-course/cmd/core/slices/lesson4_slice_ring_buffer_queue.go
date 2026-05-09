package main

import (
	"errors"
	"fmt"
)

/*
LESSON 4: Slice Ring Buffer Queue

Goal:
Implement a high-performance queue using slices.

A ring buffer avoids shifting elements.

Instead of moving data:
we move two indexes:

head -> where we dequeue
tail -> where we enqueue

Memory layout example:

capacity = 5

[ _ _ _ _ _ ]

enqueue 3 items

[ A B C _ _ ]

	 ↑     ↑
	head  tail

dequeue 2

[ A B C _ _ ]

	 ↑   ↑
	head tail

enqueue again (wrap around)

[ A B C D E ]

	 ↑
	head
*/

type RingQueue struct {
	data []string

	head int
	tail int
	size int
}

// create queue with fixed capacity.
func NewRingQueue(capacity int) *RingQueue {
	return &RingQueue{
		data: make([]string, capacity),
	}
}

/*
Enqueue adds elements to queue.

Time complexity: O(1)
*/
func (q *RingQueue) Enqueue(v string) error {
	if q.size == len(q.data) {
		return errors.New("queue full")
	}

	q.data[q.tail] = v

	// move tail forward
	q.tail = (q.tail + 1) % len(q.data)

	q.size++

	return nil
}

/*
Dequeue removes element from queue.

Time Complexity: O(1)
*/
func (q *RingQueue) Dequeue() (string, error) {
	if q.size == 0 {
		return "", errors.New("queue empty")
	}

	val := q.data[q.head]

	// avoid memory retention
	q.data[q.head] = ""

	// move head forward
	q.head = (q.head + 1) % len(q.data)

	q.size--

	return val, nil
}

/*
Print internal state for learning.
*/
func (q *RingQueue) Debug() {
	fmt.Println("buffer:", q.data)
	fmt.Println("head:", q.head, "tail:", q.tail, "size:", q.size)
	fmt.Println()
}

func Slice_Ring_Buffer_Queue() {
	queue := NewRingQueue(5)

	queue.Enqueue("job1")
	queue.Enqueue("job2")
	queue.Enqueue("job3")

	queue.Debug()

	fmt.Println("=== Dequeue ===")

	v, _ := queue.Dequeue()
	fmt.Println("processed:", v)

	v, _ = queue.Dequeue()
	fmt.Println("processed:", v)

	queue.Debug()

	fmt.Println("=== Wrap Around ===")

	queue.Enqueue("job4")
	queue.Enqueue("job5")
	queue.Enqueue("job6")

	queue.Debug()

	fmt.Println("=== Process Remaining ===")

	for queue.size > 0 {

		v, _ := queue.Dequeue()
		fmt.Println("processed:", v)

		queue.Debug()
	}
}
