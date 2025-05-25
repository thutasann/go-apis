package concurrencypatterns

import (
	"fmt"
	"math/rand"
)

// repeatFunc continuously calls the provided function `fn` and sends the result
// into a returned channel until the `done` channel is closed.
//
// Type Parameters:
//   - T: the return type of the function `fn`
//   - K: the type used for the done signal (typically empty struct{} or any placeholder)
//
// Parameters:
//   - done <-chan K: a read-only channel used to signal cancellation. When this channel is closed,
//     the goroutine stops sending values and the output channel is closed.
//   - fn func() T: the function to be repeatedly executed. Each call's return value is sent to the output channel.
//
// Returns:
//   - <-chan T: a read-only channel that emits values returned by `fn()` until the `done` channel is closed.
func repeatFunc[T any, K any](done <-chan K, fn func() T) <-chan T {
	stream := make(chan T)
	go func() {
		defer close(stream)
		for {
			select {
			case <-done:
				return
			case stream <- fn():

			}
		}
	}()
	return stream
}

// Generator Sample One
func GeneratorSampleOne() {
	done := make(chan int)
	defer close(done)

	randomNumFetcher := func() int { return rand.Intn(500000000) }
	for rand_number := range repeatFunc(done, randomNumFetcher) {
		fmt.Println(rand_number)
	}
}
