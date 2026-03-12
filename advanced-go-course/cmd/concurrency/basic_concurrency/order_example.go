package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Order struct {
	ID     int
	Status string
	mu     sync.Mutex
}

var (
// totalUpdates int
// updateMutex  sync.Mutex
)

func generateOrders(count int) []*Order {
	orders := make([]*Order, count)

	for i := range count {
		orders[i] = &Order{
			ID:     i + 1,
			Status: "Pending",
		}
	}

	return orders
}

func processOrders(orderCh <-chan *Order, wg *sync.WaitGroup) {
	defer wg.Done()
	for order := range orderCh {
		time.Sleep(
			time.Duration(rand.Intn(500)) * time.Millisecond,
		)
		fmt.Printf("Processing order %d\n", order.ID)
	}
}

// func updateOrderStatuses(orders []*Order) {
// 	for _, order := range orders {
// 		time.Sleep(time.Duration(rand.Intn(300)) * time.Millisecond)
// 		status := []string{"Processing", "Shipped", "Delivered"}[rand.Intn(3)]
// 		order.Status = status
// 		fmt.Printf(
// 			"Updated order %d status: %s\n", order.ID, status,
// 		)
// 	}
// }

// === Mutex Example ===
// func updateOrderStatus(order *Order) {
// 	order.mu.Lock()
// 	time.Sleep(time.Duration(rand.Intn(300)) * time.Millisecond)
// 	status := []string{"Processing", "Shipped", "Delivered"}[rand.Intn(3)]
// 	order.Status = status
// 	fmt.Printf(
// 		"Updated order %d status: %s\n", order.ID, status,
// 	)
// 	order.mu.Unlock()

// 	updateMutex.Lock()
// 	defer updateMutex.Unlock()
// 	currentUpdates := totalUpdates
// 	time.Sleep(5 * time.Millisecond)
// 	totalUpdates = currentUpdates + 1
// }

// func reportOrderStatus(orders []*Order) {
// 	fmt.Println("\n--- Order status report ---")
// 	for _, order := range orders {
// 		fmt.Printf("Order %d: %s\n", order.ID, order.Status)
// 	}
// 	fmt.Println("-----------------------------------------")
// }

func Order_Example() {
	var wg sync.WaitGroup
	wg.Add(2)
	orderCh := make(chan *Order)

	go func() {
		defer wg.Done()
		for _, order := range generateOrders(20) {
			orderCh <- order
		}
		close(orderCh)
		fmt.Println("Done with generating orders")
	}()

	// process orders goroutine
	go processOrders(orderCh, &wg)

	// update order statuses goroutine
	// for range 3 {
	// 	go func() {
	// 		defer wg.Done()
	// 		for _, order := range orders {
	// 			updateOrderStatus(order)
	// 		}
	// 	}()
	// }

	wg.Wait()

	// reportOrderStatus(orders)

	fmt.Println("All operations completed. Existing")
	// fmt.Println("totalUpdates --> ", totalUpdates)
}
