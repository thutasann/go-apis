package main

import "fmt"

func sample_one() {
	users := []string{"A", "B", "C", "D"}
	filtered := users[:0]
	fmt.Println(len(filtered))
	fmt.Println(cap(filtered))
}
