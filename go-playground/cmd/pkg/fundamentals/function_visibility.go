package fundamentals

import "fmt"

// PublicFunction is accessible from other packages
func PublicFunction() {
	fmt.Println("\n--> Public Function")
	privateFunction()
}

// Outer Inner Function
func OuterInnerFunction() {
	fmt.Println("\n--> Outer Inner Function")

	inner := func(name string) {
		fmt.Println("Hello", name)
	}

	inner("Thuta Sanne")
}

// Outer Inner Function Two
func OuterInnerFunctionTwo() {
	fmt.Println("\n---> Outer Inner Function Two")
	result := func(a, b int) int {
		return a + b
	}(2, 4)
	fmt.Println(result)
}

// privateFunction is only accessible within the example package
func privateFunction() {
	println("This is private")
}
