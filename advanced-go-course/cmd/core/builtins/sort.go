package main

import (
	"fmt"
	"sort"
)

func Sort_Ints_Sample() {
	vars := []int{5, 2, 0, 3, 4, 9, 6}
	sort.Ints(vars)
	fmt.Println(vars)
}

func Sort_Strigns_Sample() {
	vars := []string{"Learning", "Golang", "on", "KodeKloud"}
	sort.Strings(vars)
	fmt.Println(vars)
}
