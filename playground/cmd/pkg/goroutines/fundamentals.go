package goroutines

import (
	"context"
	"fmt"
	"net/http"
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
