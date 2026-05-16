package main

import (
	"fmt"
	"time"
)

func Multiple_Buffered_Signals() {
	rainAlarm := make(chan string, 1)
	doorAlarm := make(chan string, 1)

	// Trigger fire alarm after 1 second
	go func() {
		time.Sleep(1 * time.Second)
		rainAlarm <- "Rain detected!"
	}()

	// Trigger door alaram after 2 seconds
	go func() {
		time.Sleep(2 * time.Second)
		doorAlarm <- "Door opened!"
	}()

	select {
	case msg := <-rainAlarm:
		fmt.Println("Handle: ", msg)
	case msg := <-doorAlarm:
		fmt.Println("Handle: ", msg)
	}
}
