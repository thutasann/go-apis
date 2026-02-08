package main

import (
	"fmt"
	"time"
)

func payment_service(confirm chan bool) {
	fmt.Println("Payment service: processing payment...")
	time.Sleep(2 * time.Second)
	fmt.Println("Payment service: payment received")
	confirm <- true // blocks until main receives
}

func Unbuffered_Payment() {
	confirm := make(chan bool)

	go payment_service(confirm)

	fmt.Println("Main: waiting for payment confirmation...")
	<-confirm // blocks until payment_service sends

	fmt.Println("Main: payment confirmed, continue workflow")
}
