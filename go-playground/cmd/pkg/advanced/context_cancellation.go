package advanced

import (
	"context"
	"fmt"
	"time"
)

func ContextCancellation() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	select {
	case <-time.After(3 * time.Second):
		fmt.Println("DONE....")
	case <-ctx.Done():
		fmt.Println("Timeout: ", ctx.Err())
	}
}
