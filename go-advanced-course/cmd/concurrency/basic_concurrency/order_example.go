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

func processOrders(inCh <-chan *Order, outCh chan<- *Order, wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
		close(outCh)
	}()
	for order := range inCh {
		time.Sleep(
			time.Duration(rand.Intn(500)) * time.Millisecond,
		)
		order.Status = "Processed"
		outCh <- order
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
	orderCh := make(chan *Order, 20)
	processedCh := make(chan *Order, 20)

	go func() {
		defer wg.Done()
		for _, order := range generateOrders(20) {
			orderCh <- order
		}
		close(orderCh)
		fmt.Println("Done with generating orders")
	}()

	// process orders goroutine
	go processOrders(orderCh, processedCh, &wg)

	go func() {
		defer wg.Done()

		for {
			select {
			case processedOrder, ok := <-processedCh:
				if !ok {
					fmt.Println(
						"Processing channel closed",
					)
					return
				}
				fmt.Printf("Processed order %d with status: %s\n", processedOrder.ID, processedOrder.Status)
			case <-time.After(10 * time.Second):
				fmt.Println("Timeout waiting for operations...")
				return
			}
		}
	}()

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
