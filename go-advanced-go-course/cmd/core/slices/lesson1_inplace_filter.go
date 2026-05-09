package main

import (
	"fmt"
	"time"
)

/*
LESSON 1: In-Place Filtering (Zero Allocation Pattern)

Goal:
Remove elements from a slice WITHOUT allocating a new slice.

This is critical in:
- High-throughput systems
- Streaming pipelines
- Log processors
- Network packet filtering

Core idea:
Reuse the same backing array by re-slicing to length 0,
then append valid elements into it.

Pattern:

	dst := src[:0]
	for _, v := range src {
	    if keep(v) {
	        dst = append(dst, v)
	    }
	}

Now dst reuses the original memory.
*/

// User represents a real-world domain object
type InPlaceUser struct {
	ID        int
	Active    bool
	LastLogin time.Time
}

// filterActiveUsers keeps only active users
// without allocating new backing array.
func filterActiveUsers(users []InPlaceUser) []InPlaceUser {
	// IMPORTANT:
	// users[:0] keeps same backing array
	// but sets length = 0
	filtered := users[:0]

	for _, u := range users {
		if u.Active {
			filtered = append(filtered, u)
		}
	}
	return filtered
}

func In_Place_Filter() {
	fmt.Println("==== Before Filtering ===")

	users := []InPlaceUser{
		{1, true, time.Now()},
		{2, false, time.Now()},
		{3, true, time.Now()},
		{4, false, time.Now()},
		{5, true, time.Now()},
	}

	fmt.Println("Original length:", len(users))
	fmt.Println("Original capacity:", cap(users))

	// Store original pointer
	originalPtr := &users[0]

	users = filterActiveUsers(users)
	fmt.Println("\n=== After Filtering ===")
	fmt.Println("Filtered length:", len(users))
	fmt.Println("Filtered capacity:", cap(users))

	newPtr := &users[0]

	fmt.Println("\nBacking array reused: ", originalPtr == newPtr)

	for _, u := range users {
		fmt.Printf("ID=%d Active=%v\n", u.ID, u.Active)
	}
}
