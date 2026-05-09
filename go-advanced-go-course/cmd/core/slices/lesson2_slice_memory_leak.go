package main

import (
	"fmt"
	"runtime"
)

/*
LESSON 2: Hidden Memory Leaks with Slices

A slice only stores:
- pointer to backing array
- length
- capacity

When you take a subslice, it still references
the SAME backing array.

If the original slice is huge, the GC cannot free it.
*/

// simulates reading a huge file
func simulateLargeData() []byte {
	size := 100 * 1024 * 1024

	data := make([]byte, size)

	for i := range data {
		data[i] = 'A'
	}

	return data
}

/*
BAD PATTERN

Returns a small slice but keeps the entire
100MB array in memory.
*/
func leakExample() []byte {
	data := simulateLargeData()

	// Take only first 10 bytes
	header := data[:10]

	// data goes out of scope,
	// but header still references the array
	return header
}

/*
GOOD PATTERN

Copy the needed data into a new slice
so the large array can be garbage collected.
*/
func safeExample() []byte {
	data := simulateLargeData()

	header := make([]byte, 10)

	copy(header, data[:10])

	return header
}

func printMemory(label string) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	fmt.Printf(
		"%s: Alloc = %d MB\n",
		label,
		m.Alloc/1024/1024,
	)
}

func Slice_Memory_Leak() {
	printMemory("Start")

	// BAD case
	h := leakExample()

	runtime.GC()

	printMemory("After leakExample")

	fmt.Println("Header:", h)

	// GOOD case
	s := safeExample()

	runtime.GC()

	printMemory("After safeExample")

	fmt.Println("Header:", s)

	/*
		KEY LESSON:

		A slice keeps the ENTIRE backing array alive.

		If you only need a small piece,
		COPY it into a new slice.
	*/
}
