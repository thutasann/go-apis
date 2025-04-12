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

// privateFunction is only accessible within the example package
func privateFunction() {
	println("This is private")
}
