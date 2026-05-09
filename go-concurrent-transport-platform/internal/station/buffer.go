package station

import "context"

// Buffer is a bounded channel-based queue
type Buffer[T any] struct {
	ch chan T
}

func NewBuffer[T any](size int) *Buffer[T] {
	return &Buffer[T]{
		ch: make(chan T, size),
	}
}

func (b *Buffer[T]) Enqueue(ctx context.Context, v T) bool {
	select {
	case b.ch <- v:
		return true
	case <-ctx.Done():
		return false
	}
}

func (b *Buffer[T]) Dequeue() <-chan T {
	return b.ch
}
