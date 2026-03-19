package main

import "fmt"

type Counter struct {
	count int
}

func (c *Counter) Increment() {
	c.count += 1
}

func NewCounter() *Counter {
	return new(Counter)
}

func New_Sample() {
	counter := NewCounter()
	counter.Increment()
	fmt.Println("count: ", counter.count)
}
