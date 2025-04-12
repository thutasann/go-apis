package main

import (
	"fmt"

	"github.com/thutasann/playground/cmd/pkg/fundamentals"
)

// Playground App 1
func main() {
	fmt.Println("----- Playground -----")
	fmt.Println(fundamentals.Hello)

	fmt.Println("\n----- Fundamentals -----")
	fundamentals.RuneSampleOne()
	fundamentals.OuterInnerFunction()
	fundamentals.OuterInnerFunctionTwo()
	fundamentals.PublicFunction()

	fmt.Println("\n----- Pointers -----")
	fundamentals.PointerSampleOne()
	fundamentals.ModifyPointerFunctionSample()
	fundamentals.PointerStructSampleOne()
	fundamentals.DoublePointer()
	fundamentals.ArraySliceModify()
	fundamentals.StructPointerSampleOne()
	fundamentals.StringVsPointerString()
	fundamentals.StructPointerSampleTwo()
	fundamentals.WithoutPointerSample()
	fundamentals.WithPointerSample()
	fundamentals.ModifyingConfig()

	fundamentals.Person.Greet(fundamentals.Person{Name: "Thuta", Age: 20})
	// fundamentals.TickerSampleOne()
	fundamentals.MutexSamples()
	fundamentals.DeferSampleOne()
	fundamentals.DeferInsideLoop()
	fundamentals.ArraysSampleOne()
	fundamentals.SlicesSampleOne()

	fmt.Println("----- Demos -----")
	// demos.ATMMutexDemo()
	// demos.DeferDemo()
}
