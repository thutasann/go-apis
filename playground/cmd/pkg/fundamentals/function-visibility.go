package fundamentals

import "fmt"

// PublicFunction is accessible from other packages
func PublicFunction() {
	fmt.Println("Public Function")
	privateFunction()
}

// privateFunction is only accessible within the example package
func privateFunction() {
	println("This is private")
}
