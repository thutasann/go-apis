package main

import "fmt"

// nums ...int === nums []int
func v_sum(nums ...int) int {
	total := 0

	for _, n := range nums {
		total += n
	}

	return total
}

func Variadic_Function() {
	fmt.Println(v_sum(1, 2))
	fmt.Println(v_sum(1, 2, 3))
}
