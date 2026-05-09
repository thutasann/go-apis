package concurrencypatterns

import (
	"fmt"
	"time"
)

// Send only Channel
func sampleChan(ch chan<- string, data string) {
	ch <- data
}

// Channel and Select Sample One
func ChannelSampleOne() {
	myChannel := make(chan string)
	anotherChannel := make(chan string)

	go sampleChan(myChannel, "data")
	go sampleChan(anotherChannel, "cow")

	select {
	case msgFromMyChannel := <-myChannel:
		fmt.Println("msgFromMyChannel:", msgFromMyChannel)
	case msgFromAnotherChannel := <-anotherChannel:
		fmt.Println("msgFromAnotherChannel:", msgFromAnotherChannel)
	}
}

// Worker only RECEIVES from the jobs channel (receive-only)
func channel_worker(id int, jobs <-chan int, done chan<- bool) {
	for job := range jobs {
		fmt.Printf("Worker %d started job %d\n", id, job)
		time.Sleep(time.Second)
		fmt.Printf("Worker %d finished job %d\n", id, job)
	}
	done <- true
}

// Receive only channel sample
func ReceiveOnlyChannelSample() {
	jobs := make(chan int, 5)
	done := make(chan bool)

	// start 2 workers
	go channel_worker(1, jobs, done)
	go channel_worker(2, jobs, done)

	// send 5 jobs
	for j := 1; j <= 5; j++ {
		jobs <- j
	}

	close(jobs) // close channel so workers know there's no more data

	// wait for both workers to finish
	<-done
	<-done

	fmt.Println("All jobs completed")
}
