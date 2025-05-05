package concurrencypatterns

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
func RepeatFunc[T any, K any](done <-chan K, fn func() T) <-chan T {
	return nil
}
