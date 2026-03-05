package main

import "fmt"

// slice[low:high]
// low is inclusive, high is exclusive
// index 0 → included
// index 1 → included
// index 2 → excluded
func sample_one() {
	users := []string{"A", "B", "C", "D"}
	filtered := users[0:2]
	// fmt.Println(len(filtered))
	// fmt.Println(cap(filtered))
	fmt.Println(filtered)
}
