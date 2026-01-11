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
	// concurrencypatterns.ContextSampleOne()
	// concurrencypatterns.ThreadSafeSample()

	// concurrencypatterns.GeneratorSampleOne()
	// concurrencypatterns.Concurrent_Files_Reading_Sample()
	// concurrencypatterns.Concurrent_TCP_Server()
	// concurrencypatterns.Concurrent_TCP_Client_Pool()
	// concurrencypatterns.Network_Aggr()
	concurrencypatterns.TCP_Context_Timeout()
}
