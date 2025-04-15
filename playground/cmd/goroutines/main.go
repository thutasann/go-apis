package main

import (
	"fmt"
	"net/http"
)

func fetchURL(url string, ch chan string) {
	resp, err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprintf("error: %v", err)
		return
	}
	defer resp.Body.Close()
	ch <- fmt.Sprintf("%s: %d", url, resp.StatusCode)
}

func fetchURLsSample() {
	urls := []string{"https://pokeapi.co/api/v2/pokemon", "https://jsonplaceholder.typicode.com/todos", "https://golang.org"}
	ch := make(chan string)

	for _, url := range urls {
		go fetchURL(url, ch)
	}

	for range urls {
		fmt.Println(<-ch)
	}
}

func main() {
	fetchURLsSample()
}
