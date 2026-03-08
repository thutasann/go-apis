package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Order struct {
	ID     int
	Status string
}

func generateOrders(count int) []*Order {
	orders := make([]*Order, count)

	for i := 0; i < count; i++ {
		orders[i] = &Order{
			ID:     i + 1,
			Status: "Pending",
		}
	}

	return orders
}

func processOrders(orders []*Order) {
	for _, order := range orders {
		time.Sleep(
			time.Duration(rand.Intn(500)) * time.Millisecond,
		)
		fmt.Printf("Processing order %d\n", order.ID)
	}
}

func updateOrderStatuses(orders []*Order) {
	for _, order := range orders {
		time.Sleep(time.Duration(rand.Intn(300)) * time.Millisecond)
		status := []string{"Processing", "Shipped", "Delivered"}[rand.Intn(3)]
		order.Status = status
		fmt.Printf(
			"Updated order %d status: %s\n", order.ID, status,
		)
	}
}

func reportOrderStatus(orders []*Order) {
	for range 5 {
		time.Sleep(1 * time.Second)
		fmt.Println("\n--- Order status report ---")
		for _, order := range orders {
			fmt.Printf("Order %d: %s\n", order.ID, order.Status)
		}
		fmt.Println("-----------------------------------------")
	}
}

func Order_Example() {
	orders := generateOrders(20)

	processOrders(orders)

	updateOrderStatuses(orders)

	reportOrderStatus(orders)

	fmt.Println("All operations completed. Existing")
}
