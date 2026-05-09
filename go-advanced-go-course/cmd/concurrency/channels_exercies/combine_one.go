package main

import (
	"context"
	"fmt"
	"time"
)

type Result struct {
	Name string
	Data string
	Err  error
}

func fetchProfile(ctx context.Context, out chan<- Result) {
	select {
	case <-time.After(1 * time.Second):
		out <- Result{Name: "profile", Data: "Alice"}
	case <-ctx.Done():
		return
	}
}

func fetchBalance(ctx context.Context, out chan<- Result) {
	select {
	case <-time.After(1500 * time.Millisecond):
		out <- Result{Name: "balance", Data: "$420"}
	case <-ctx.Done():
		return
	}
}

func fetchRecommendations(ctx context.Context, out chan<- Result) {
	select {
	case <-time.After(3 * time.Second): // slow service
		out <- Result{Name: "recommendations", Data: "Book, Laptop"}
	case <-ctx.Done():
		return
	}
}

func Combine_One() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	results := make(chan Result)

	// fan-out: parallel work
	go fetchProfile(ctx, results)
	go fetchBalance(ctx, results)
	go fetchRecommendations(ctx, results)

	collected := make(map[string]string)

	for range 3 {
		select {
		case res := <-results:
			if res.Err != nil {
				fmt.Println("Error:", res.Err)
				cancel()
				return
			}
			fmt.Println("Received:", res.Name)
			collected[res.Name] = res.Data

		case <-ctx.Done():
			fmt.Println("Request cancelled:", ctx.Err())
			fmt.Println("Partial results:", collected)
			return
		}
	}

	fmt.Println("All data collected:", collected)
}
