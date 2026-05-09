package main

import (
	"context"
	"fmt"
	"time"
)

func dbCall(ctx context.Context) error {
	select {
	case <-time.After(1 * time.Second):
		fmt.Println("DB done")
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func apiCall(ctx context.Context) error {
	select {
	case <-time.After(2 * time.Second):
		fmt.Println("API done")
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// - DB succeeds
// - API exceeds deadline
// - API exceeds deadline
func Shared_Timeout() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := dbCall(ctx); err != nil {
		fmt.Println("DB error:", err)
		return
	}

	if err := apiCall(ctx); err != nil {
		fmt.Println("API error:", err)
		return
	}

	fmt.Println("Request completed successfully")
}
