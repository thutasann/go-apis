package main

import (
	"fmt"
	"time"
)

func paymentAPI(result chan string) {
	time.Sleep(3 * time.Second)
	result <- "payment success"
}

func Select_Timeout() {
	result := make(chan string)

	go paymentAPI(result)

	select {
	case res := <-result:
		fmt.Println("Received: ", res)
	case <-time.After(2 * time.Second):
		fmt.Println("Timeout: payment service not responding")
	}
}
