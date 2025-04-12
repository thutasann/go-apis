package main

import (
	"fmt"

	"github.com/thutasann/playground/cmd/pkg/fundamentals"
)

// Playground App 1
func main() {
	fmt.Println("----- Playground -----")
	fmt.Println(fundamentals.Hello)

	fmt.Println("----- Fundamentals -----")
	fundamentals.RuneSampleOne()
	fundamentals.OuterInnerFunction()
	fundamentals.PublicFunction()
	fundamentals.PointerSampleOne()
	fundamentals.ModifyPointerFunctionSample()
	fundamentals.PointerStructSampleOne()
	fundamentals.DoublePointer()
	fundamentals.ArraySliceModify()
	fundamentals.StructPointerSampleOne()
	fundamentals.Person.Greet(fundamentals.Person{Name: "Thuta", Age: 20})
	fundamentals.StructPointerSampleTwo()
	fundamentals.StringVsPointerString()
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
