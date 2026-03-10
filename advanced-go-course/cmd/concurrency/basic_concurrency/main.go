package main

import "fmt"

func main() {
	orders := generateOrders(20)

	go processOrders(orders)

	go updateOrderStatuses(orders)

	go reportOrderStatus(orders)

	fmt.Println("All operations completed. Existing")
}
