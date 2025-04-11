package main

import (
	"fmt"

	"github.com/thutasann/playground/cmd/pkg/bytes"
)

// Bytes
func main() {
	fmt.Println("===> Playground App 2")
	bytes.BufferSampleOne()
	bytes.BytesSamples()
	// bytes.BufferedNetworkIO()
	bytes.StreamingJSONToBuffer()
	bytes.BufferedFileWriter()
}
