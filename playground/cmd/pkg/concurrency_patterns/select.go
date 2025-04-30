package concurrencypatterns

import (
	"fmt"
	"time"
)

func slowAPI(response chan<- string) {
	time.Sleep(3 * time.Second)
	response <- "🌐 API Response"
}

// 🧵 1. Timeout with select (e.g., calling an API)
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

func fetchA(out chan<- string) {
	time.Sleep(1 * time.Second)
	out <- "🍎 From Service A"
}

func fetchB(out chan<- string) {
	time.Sleep(2 * time.Second)
	out <- "🍌 From Service B"
}

func FanInPatternMergeMultipleResourcesIntoOne() {
	aChan := make(chan string)
	bChan := make(chan string)

	go fetchA(aChan)
	go fetchB(bChan)

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

}
