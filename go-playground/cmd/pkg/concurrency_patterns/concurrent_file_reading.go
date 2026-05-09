package concurrencypatterns

import (
	"bufio"
	"fmt"
	"os"
	"sync"
)

// Read multiple files concurrently
// Counts number of lines in each file
// Collects results safely

// Result represents the output of one goroutine
type Result struct {
	FileName string
	Lines    int
	Error    error
}

// countLines reads a file and counts lines
func countLines(filepath string, wg *sync.WaitGroup, results chan<- Result) {
	defer wg.Done()

	file, err := os.Open(filepath)
	if err != nil {
		results <- Result{FileName: filepath, Error: err}
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := 0

	for scanner.Scan() {
		lines++
	}

	if err := scanner.Err(); err != nil {
		results <- Result{FileName: filepath, Error: err}
		return
	}

	results <- Result{
		FileName: filepath,
		Lines:    lines,
	}
}

func Concurrent_Files_Reading_Sample() {
	files := []string{
		"log1.txt",
		"log2.txt",
		"log3.txt",
	}

	results := make(chan Result)
	var wg sync.WaitGroup

	// Launch a goroutine per file
	for _, file := range files {
		wg.Add(1)
		go countLines(file, &wg, results)
	}

	// Close channel after all goroutines finish
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	for res := range results {
		if res.Error != nil {
			fmt.Printf("Error reading %s: %v\n", res.FileName, res.Error)
			continue
		}
		fmt.Printf("File: %s, Lines: %d\n", res.FileName, res.Lines)
	}
}
