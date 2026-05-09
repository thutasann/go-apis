package concurrencypatterns

import (
	"fmt"
	"time"
)

func slowAPI(response chan<- string) {
	time.Sleep(3 * time.Second)
	response <- "🌐 API Response"
}

func TimeoutWithSelect() {
	response := make(chan string)
	go slowAPI(response)

	select {
	case res := <-response:
		fmt.Println("✅ Got:", res)
	case <-time.After(2 * time.Second):
		fmt.Println("⏰ Timeout: API took too long")
	}
}

func fetch_a(out chan<- string) {
	time.Sleep(1 * time.Second)
	out <- "🍎 From Service A"
}

func fetch_b(out chan<- string) {
	time.Sleep(2 * time.Second)
	out <- "🍌 From Service B"
}

func FanInPatternMergeMultipleResourcesIntoOne() {
	aChan := make(chan string)
	bChan := make(chan string)

	go fetch_a(aChan)
	go fetch_b(bChan)

	for i := 0; i < 2; i++ {
		select {
		case a := <-aChan:
			fmt.Println("🔻 Got:", a)
		case b := <-bChan:
			fmt.Println("🔻 Got:", b)
		}
	}
}

func ForSelectSampleOne() {
	charChannel := make(chan string, 3)
	chars := []string{"a", "b", "c"}

	for _, s := range chars {
		select {
		case charChannel <- s:
		default:
			fmt.Println("default")
		}
	}

	close(charChannel)

	for result := range charChannel {
		fmt.Println(result)
	}
}

func InfiniteLoopingGoRoutines() {
	go func() {
		for {
			select {
			default:
				fmt.Println("DOING WORK...")
			}
		}
	}()
	time.Sleep(time.Second * 10)
}

func worker_loop(id int, tasks <-chan int, quit <-chan struct{}) {
	for {
		select {
		case task := <-tasks:
			fmt.Printf("👷 Worker %d processing task %d\n", id, task)
		case <-quit:
			fmt.Printf("❌ Worker %d shutting down\n", id)
			return
		}
	}
}

func WorkerLoopThatListensForTasks() {
	tasks := make(chan int)
	quit := make(chan struct{})

	go worker_loop(1, tasks, quit)

	for i := 1; i <= 3; i++ {
		tasks <- i
		time.Sleep(500 * time.Millisecond)
	}

	quit <- struct{}{}
	time.Sleep(1 * time.Second)
}

func PollingWithTimeoutInBetween() {
	data := make(chan string)

	go func() {
		time.Sleep(2 * time.Second)
		data <- "✅ Data received"
	}()

	for {
		select {
		case d := <-data:
			fmt.Println("📦 Got:", d)
			return
		case <-time.After(1 * time.Second):
			fmt.Println("⏰ Still waiting...")
		}
	}
}

func do_work(done <-chan bool) {
	for {
		select {
		case <-done:
			return
		default:
			fmt.Println("DOING WORK")
		}
	}
}

func DoneChannel() {
	done := make(chan bool)
	go do_work(done)
	time.Sleep(time.Second * 3)
	close(done)
}
