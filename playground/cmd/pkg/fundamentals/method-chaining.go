package fundamentals

// ChanPipeline holds a channel and supports method chaining
type ChanPipline struct {
	data chan string
}

// NewPipeline initializes the pipeline
func NwePipeline(input []string) *ChanPipline {
	ch := make(chan string, len(input))
	for _, val := range input {
		ch <- val
	}
	close(ch)
	return &ChanPipline{data: ch}
}
