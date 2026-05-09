package goroutines

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"time"
)

func fetchUrl(url string, ch chan string) {
	resp, err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprintf("error: %v", err)
		return
	}
	defer resp.Body.Close()
	ch <- fmt.Sprintf("%s: %d", url, resp.StatusCode)
}

func FetchURLsSample() {
	urls := []string{"https://pokeapi.co/api/v2/pokemon", "https://jsonplaceholder.typicode.com/todos", "https://golang.org"}
	ch := make(chan string)

	for _, url := range urls {
		go fetchUrl(url, ch)
	}

	for range urls {
		fmt.Println(<-ch)
	}
}

func worker(id int, jobs <-chan int, results chan<- int) {
	for job := range jobs {
		fmt.Printf("Worker %d processing job %d\n", id, job)
		time.Sleep(time.Second)
		results <- job * 2
	}
}

func WorkerPoolSample() {
	jobs := make(chan int, 5)
	results := make(chan int, 5)

	for w := 1; w <= 3; w++ {
		go worker(w, jobs, results)
	}

	for j := 1; j <= 5; j++ {
		jobs <- j
	}

	close(jobs)

	for a := 1; a <= 5; a++ {
		fmt.Println("Result:", <-results)
	}
}

func ChannelSample() {
	fmt.Println("----> Channel Sample")
	ch := make(chan int) // create a channel of int type
	ch <- 10             // send value 10 into the channel
	val := <-ch          // Receive from channel
	fmt.Println("val -->", val)
}

func chan_worker(id int, jobs <-chan int, results chan<- int) {
	for job := range jobs {
		fmt.Printf("Worker %d started job %d\n", id, job)
		time.Sleep(time.Second)
		results <- job * 2
	}
}

func JobQueueSample() {
	jobs := make(chan int, 5)
	results := make(chan int, 5)

	for w := 1; w <= 3; w++ {
		go chan_worker(w, jobs, results)
	}

	for j := 1; j <= 5; j++ {
		jobs <- j
	}

	close(jobs)

	for a := 1; a <= 5; a++ {
		fmt.Println("Result:", <-results)
	}
}

func producer(id int, ch chan<- string) {
	for i := 0; i < 3; i++ {
		ch <- fmt.Sprintf("Producer %d: data %d", id, i)
	}
}

func MultpleProducersOneConsumerPattern() {
	ch := make(chan string)

	for i := 0; i < 3; i++ {
		go producer(i, ch)
	}

	for i := 0; i < 9; i++ {
		fmt.Println(<-ch)
	}
}

func graceful_server(ctx context.Context, done chan<- string) {
	for {
		select {
		case <-ctx.Done():
			done <- "server stopped"
			return
		default:
			fmt.Println("working...")
			time.Sleep(time.Second)
		}
	}
}

func GracefulShutdownSample() {
	ctx, channel := context.WithTimeout(context.Background(), 3*time.Second)
	defer channel()

	done := make(chan string)
	go graceful_server(ctx, done)

	fmt.Println(<-done)
}

func chat_broadcaster(messages <-chan string) {
	for msg := range messages {
		fmt.Println("Broadcast:", msg)
	}
}

func chat_client(id int, messges chan<- string) {
	for i := 0; i < 2; i++ {
		messges <- fmt.Sprintf("Client %d: Hello %d", id, i)
	}
}

func RealTimeChatServer() {
	msgs := make(chan string)
	go chat_broadcaster(msgs)

	for i := 0; i < 3; i++ {
		go chat_client(i, msgs)
	}

	time.Sleep(2 * time.Second)
}

func SelectSample() {
	ch := make(chan string)

	go func() {
		time.Sleep(2 * time.Second)
		ch <- "data"
	}()

	select {
	case res := <-ch:
		fmt.Println("Data Received:", res)
	case <-time.After(3 * time.Second):
		fmt.Println("Timeout!")
	}
}

func SelectMultiplexing() {
	ch1 := make(chan string)
	ch2 := make(chan string)

	go func() {
		time.Sleep(1 * time.Second)
		ch1 <- "from cha1"
	}()

	go func() {
		time.Sleep(2 * time.Second)
		ch2 <- "from ch2"
	}()

	select {
	case msg1 := <-ch1:
		fmt.Println(msg1)
	case msg2 := <-ch2:
		fmt.Println(msg2)
	}
}

func select_fetch_data(ch chan string) {
	time.Sleep(3 * time.Second)
	ch <- "data"
}

func TimeOutAndDeadliens() {
	ch := make(chan string)
	go select_fetch_data(ch)

	select {
	case data := <-ch:
		fmt.Println("Received:", data)
	case <-time.After(2 * time.Second):
		fmt.Println("timeout: no data")
	}
}

func fan_in_select_source(name string, delay time.Duration, out chan<- string) {
	for i := 0; ; i++ {
		time.Sleep(delay)
		out <- fmt.Sprintf("%s: %d", name, i)
	}
}

// Fan-In Pattern
func FanInPattern() {
	ch1 := make(chan string)
	ch2 := make(chan string)

	go fan_in_select_source("API-A", 1*time.Second, ch1)
	go fan_in_select_source("API-B", 2*time.Second, ch2)
}

func doWork(ctx context.Context, ch chan<- string) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("work cancelledd: ", ctx.Err())
			return
		default:
			time.Sleep(500 * time.Millisecond)
			ch <- "working..."
		}
	}
}

// select + context.Context (for cancellation & timeout)
func CancellationAndTimeout() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	ch := make(chan string)
	go doWork(ctx, ch)

	for {
		select {
		case msg := <-ch:
			fmt.Println(msg)
		case <-ctx.Done():
			fmt.Println("Main context done: ", ctx.Err())
			return
		}
	}
}

func heartbeat(ping <-chan struct{}, done <-chan struct{}) {
	for {
		select {
		case <-ping:
			fmt.Println("Received ping")
		case <-time.After(1 * time.Second):
			fmt.Println("NO ping: system might be down")
		case <-done:
			fmt.Println("Shutting down heartbeat...")
			return
		}
	}
}

// HeatBeat Sample
func HeartBeatSample() {
	ping := make(chan struct{})
	done := make(chan struct{})

	go heartbeat(ping, done)

	time.Sleep(500 * time.Millisecond)
	ping <- struct{}{}

	time.Sleep(2 * time.Second)
	close(done)
}

// Context Help you stop goroutines that would otherwise keep running
func LeakPrevention() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		select {
		case <-ctx.Done():
			fmt.Println("Goroutine cleaned up:", ctx.Err())
		default:
			fmt.Println("default")
		}
	}()
}

func fan_out_distribute(ch chan<- string) {
	names := []string{"Alice", "Bob", "Charlie"}
	for _, name := range names {
		ch <- name
	}
	close(ch)
}

// fan_out_distribute pushes data out -- it only needs chan<- string
func FanOutDistrubuteSample() {
	nameChannel := make(chan string)
	go fan_out_distribute(nameChannel)

	for name := range nameChannel {
		fmt.Println("Received:", name)
	}
}

// Sort Sample
func SortSample() {
	strs := []string{"c", "a", "b"}
	slices.Sort(strs)
	fmt.Println("Strings:  ", strs)

	ints := []int{7, 2, 4}
	slices.Sort(ints)
	fmt.Println("Ints:  ", ints)

	s := slices.IsSorted(ints)
	fmt.Println("Sorted: ", s)
}
