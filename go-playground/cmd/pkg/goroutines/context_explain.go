// controlling tiemtouts
// cancelling goroutines
// passing metadata across go application
package goroutines

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// # Timeout
//
// - Creates a base context with no cancellation or timeout — like the root of a context tree.
//
// - Derives a new context from ctx that automatically cancels after 2 seconds (This is used to limit how long some work should run.)
//
// - Ensures we free resources associated with ctxWithTimeout, even if timeout doesn’t happen.
//
// - Creates a signal-only channel. Used to notify when the simulated "API call" (below) is done.
//
// - Spawns a goroutine that pretends to do work for 3 seconds, then signals completion by closing done.
//
// - API simulation finishing (<-done)
//
// - Timeout triggering (<-ctxWithTimeout.Done())
func ContextTimeoutExplain() {
	ctx := context.Background()

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	done := make(chan struct{})

	go func() {
		time.Sleep(3 * time.Second)
		close(done)
	}()

	select {
	case <-done:
		fmt.Println("Called the API")
	case <-ctxWithTimeout.Done():
		fmt.Println("oh no my timetout expired :>> ", ctxWithTimeout.Err())
		// do some logic to handle this
	}
}

// # Context with Value Explain
func ContextWithValueExplain() {
	type key int
	const UserKey key = 0

	ctx := context.Background()

	ctxWithValue := context.WithValue(ctx, UserKey, "123")

	if userID, ok := ctxWithValue.Value(UserKey).(string); ok {
		fmt.Println("this is user ID", userID)
	} else {
		fmt.Println("this is a protected route - no userID found")
	}
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	select {
	case <-time.After(3 * time.Second):
		fmt.Println("API response!")
	case <-ctx.Done():
		fmt.Println("oh no the context expired")
		http.Error(w, "Request context timeout", http.StatusRequestTimeout)
	}
}
