package main

import "fmt"

func Label_Sample() {
OuterLoop:
	for i := range 3 {
		for j := range 3 {
			if i == 1 && j == 1 {
				fmt.Println("Breaking loop when i =", i, "and j =", j)
				break OuterLoop
			}
			fmt.Printf("i = %d, j = %d\n", i, j)
		}
	}
}

func Label_Sample_Two() {
	for i := range 5 {
		for j := range 5 {
			if i == 2 && j == 2 {
				fmt.Println("Skipping rest of inner loop")
				goto nextIteration
			}
			fmt.Printf("i = %d, j = %d\n", i, j)
		}
	nextIteration:
	}
}
