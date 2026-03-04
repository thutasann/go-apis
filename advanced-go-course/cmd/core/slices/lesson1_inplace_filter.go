package main

import "time"

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
type User struct {
	ID        int
	Active    bool
	LastLogin time.Time
}

// filterActiveUsers keeps only active users
// without allocating new backing array.
// func filterActiveUsers(users []User) []User {
// 	// IMPORTANT:
// 	// users[:0] keeps same backing array
// 	// but sets length = 0
// }
