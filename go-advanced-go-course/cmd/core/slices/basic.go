package main

import "fmt"

// slice[low:high]
// low is inclusive, high is exclusive
// index 0 → included
// index 1 → included
// index 2 → excluded
func Sample_One() {
	users := []string{"A", "B", "C", "D"}
	filtered := users[0:2]
	// fmt.Println(len(filtered))
	// fmt.Println(cap(filtered))
	fmt.Println(filtered)
}

func Backing_Array() {
	arr := [4]string{"A", "B", "C", "D"}
	fmt.Println("arr", arr)

	s := arr[0:2]
	fmt.Println(s)
}

// Modifying One Slice Can Affect Another
func Modifying_Slices() {
	users := []string{"A", "B", "C", "D"}

	a := users[0:2]
	b := users[1:3]

	a[1] = "X"

	fmt.Println("users", users)
	fmt.Println("b", b)
}
