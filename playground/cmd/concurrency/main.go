package main

import concurrencypatterns "github.com/thutasann/playground/cmd/pkg/concurrency_patterns"

// Concurrency Patterns
func main() {
	// concurrencypatterns.GoRoutineSampleOne()
	// concurrencypatterns.ReceiveOnlyChannelSample()
	// concurrencypatterns.ChannelSampleOne()

	// concurrencypatterns.TimeoutWithSelect()
	// concurrencypatterns.FanInPatternMergeMultipleResourcesIntoOne()
	// concurrencypatterns.ForSelectSampleOne()
	// concurrencypatterns.InfiniteLoopingGoRoutines()
	// concurrencypatterns.WorkerLoopThatListensForTasks()
	// concurrencypatterns.PollingWithTimeoutInBetween()
	// concurrencypatterns.DoneChannel()

	// concurrencypatterns.PipelineSampleOne()
	concurrencypatterns.ContextSampleOne()
}
