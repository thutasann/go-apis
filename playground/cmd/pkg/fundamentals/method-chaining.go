package fundamentals

import (
	"fmt"
	"strings"
)

// ChanPipeline holds a channel and supports method chaining
type ChanPipeline struct {
	data chan string
}

// NewPipeline initializes the pipeline
func NewPipeline(input []string) *ChanPipeline {
	ch := make(chan string, len(input))
	for _, val := range input {
		ch <- val
	}
	close(ch)
	return &ChanPipeline{data: ch}
}

// ToUpper converts all strings to uppercase
func (p *ChanPipeline) ToUpper() *ChanPipeline {
	out := make(chan string, 10) // Buffered
	go func() {
		for val := range p.data {
			out <- strings.ToUpper(val)
		}
		close(out)
	}()
	p.data = out
	return p
}

// Filter removes strings that contain a keyword
func (p *ChanPipeline) Filter(keyword string) *ChanPipeline {
	out := make(chan string, 10) // Buffered
	go func() {
		for val := range p.data {
			if !strings.Contains(val, keyword) {
				out <- val
			}
		}
		close(out)
	}()
	p.data = out
	return p
}

// Print outputs the final result
func (p *ChanPipeline) Print() {
	for val := range p.data {
		fmt.Println(val)
	}
}

func MethodChainExample() {
	input := []string{"hello", "world", "go", "golang", "great"}
	NewPipeline(input).
		ToUpper().
		Filter("GO").
		Print()
}
