package main

import (
	"fmt"

	"github.com/thutasann/playground/cmd/pkg/fundamentals"
)

// Playground Main
func main() {
	fmt.Println("----- Playground -----")
	fmt.Println(fundamentals.Hello)

	fmt.Println("----- Fundamentals -----")
	fundamentals.PublicFunction()
	fundamentals.PointerSampleOne()
	fundamentals.ModifyPointerFunctionSample()
	fundamentals.PointerStructSampleOne()
	fundamentals.DoublePointer()
	fundamentals.ArraySliceModify()
	fundamentals.StructPointerSampleOne()
	fundamentals.Person.Greet(fundamentals.Person{Name: "Thuta", Age: 20})
	fundamentals.StructPointerSampleTwo()
}
