package concurrencypatterns

import (
	"bufio"
	"fmt"
	"net"
)

// handleConnection runs in its own goroutine
func handleConnection(conn net.Conn) {
	defer conn.Close()

	clientAddr := conn.RemoteAddr().String()
	fmt.Println("New client connected: ", clientAddr)

	reader := bufio.NewReader(conn)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Client disconnected: ", clientAddr)
			return
		}

		fmt.Printf("Received from %s: %s", clientAddr, message)

		response := fmt.Sprintf("Echo: %s", message)
		conn.Write([]byte(response))
	}
}

func Concurrent_TCP_Server() {
	listener, err := net.Listen("tcp", ":8000")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	fmt.Println("TCP server is listening on : 8000")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Accept error: ", err)
			continue
		}

		go handleConnection(conn)
	}
}
